package builder

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"os"
	"path"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	dkarch "github.com/docker/docker/pkg/archive"
	"github.com/hashicorp/go-getter"
	"github.com/morty-faas/registry/pkg/archive"
	"github.com/morty-faas/registry/pkg/helpers"
	"github.com/morty-faas/registry/pkg/sys"
	cp "github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
)

// Builder is used to encapsulate logic used to build functions.
type (
	Builder struct {
		client            *client.Client
		runtimesDirectory string
	}

	BuildOptions struct {
		// The identifier for the current build
		Id string
		// The runtime to use to build the function
		Runtime string
		Archive io.Reader
	}
)

const (
	destination     = "/tmp/morty-runtimes"
	runtimeEndpoint = "github.com/morty-faas/runtimes.git"
	branch          = "main"
	alphaInitScript = `#!/bin/sh
source /app/env.sh
/usr/bin/alpha
	`
)

var (
	ErrFailedToInitializeDockerClient        = errors.New("failed to initialize docker client")
	ErrFailedToDownloadRuntimes              = errors.New("unable to download runtimes repository")
	ErrInvalidRuntime                        = errors.New("runtime is invalid")
	ErrFailedToSaveArchiveToWorkingDirectory = errors.New("failed to save archive file to the current working directory")
	ErrInjectingCodeIntoRuntime              = errors.New("the function code can't be injected into the given runtime")
	ErrCreatingTemporaryBuildDirectory       = errors.New("failed to create a temporary build directory")
	ErrCreateBuildContainer                  = errors.New("failed to create the build container")
	ErrStartBuildContainer                   = errors.New("failed to start the build container")
	ErrInjectAlphaScript                     = errors.New("failed to inject the custom /sbin/init script")
)

func NewBuilder() (*Builder, error) {
	log.Info("bootstrapping new function builder")

	dir, err := downloadRuntimes()
	if err != nil {
		return nil, err
	}

	// Initialize the docker client for future requests
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, helpers.WrapError(ErrFailedToInitializeDockerClient, err)
	}

	return &Builder{client, dir}, nil
}

// ImageBuild package an entire function image, and return the path to the image on disk.
func (b *Builder) ImageBuild(ctx context.Context, opts *BuildOptions) (string, error) {
	buildId := opts.Id

	log.Infof("build/%s: starting new build", buildId)

	pathToRuntime := path.Join(b.runtimesDirectory, opts.Runtime)

	// Check that the runtime exists before going further in the process
	if _, err := os.Stat(pathToRuntime); os.IsNotExist(err) {
		return "", ErrInvalidRuntime
	}

	// Create a temporary directory based with a copy of the runtime
	workingDirectory := path.Join("/tmp", buildId)
	if err := cp.Copy(pathToRuntime, workingDirectory); err != nil {
		return "", helpers.WrapError(ErrCreatingTemporaryBuildDirectory, err)
	}

	// Write the received archive in the working directory
	zipPath := path.Join(workingDirectory, "function.zip")

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", helpers.WrapError(ErrFailedToSaveArchiveToWorkingDirectory, err)
	}
	defer zipFile.Close()

	if _, err := io.Copy(zipFile, opts.Archive); err != nil {
		return "", helpers.WrapError(ErrFailedToSaveArchiveToWorkingDirectory, err)
	}

	// Decompress the given function code archive into the runtime template
	// We can't continue the process if an errors occurs here because we
	// need to inject the user provided code to be able to build the function.
	if err := archive.Unzip(zipPath, path.Join(workingDirectory, "function")); err != nil {
		return "", helpers.WrapError(ErrInjectingCodeIntoRuntime, err)
	}

	log.Tracef("build context set to %s", workingDirectory)

	buildCtx, err := dkarch.TarWithOptions(workingDirectory, &dkarch.TarOptions{})
	if err != nil {
		return "", err
	}

	buildOpts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{buildId},
		Remove:     true,
	}

	buildResponse, err := b.client.ImageBuild(ctx, buildCtx, buildOpts)
	if err != nil {
		return "", err
	}
	defer buildResponse.Body.Close()

	// Stream the logs into the console
	streamDockerLogs(buildId, buildResponse.Body)

	log.Infof("build/%s: docker build successful", buildId)

	// Now that we have our docker image, we need to create a ext4
	// mount point for our rootfs

	// Size of the image
	size := b.getImageSize(ctx, buildId)

	rootfs, mountpoint, err := makeExt4FS(workingDirectory, size)
	if err != nil {
		return "", err
	}

	// Options for the temporary container

	containerOpts := &container.Config{
		Image: buildId,
		Cmd: []string{"sh", "-c", `
for d in app bin etc lib root sbin usr; do tar c "/${d}" | tar x -C /my-rootfs; done;
for dir in dev proc run sys var; do mkdir /my-rootfs/${dir}; done;
exit;
		`},
		Tty: false,
	}

	containerHostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: mountpoint,
				Target: "/my-rootfs",
			},
			{
				Type:   mount.TypeBind,
				Source: "/dev/urandom",
				Target: "/dev/random",
			},
		},
	}

	log.Tracef("build/%s: creating build container", buildId)
	c, err := b.client.ContainerCreate(ctx, containerOpts, containerHostConfig, nil, nil, "")
	if err != nil {
		return "", helpers.WrapError(ErrCreateBuildContainer, err)
	}

	log.Tracef("build/%s: starting container with id=%s", buildId, c.ID)
	if err := b.client.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
		return "", helpers.WrapError(ErrStartBuildContainer, err)
	}

	statusCh, errCh := b.client.ContainerWait(ctx, c.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}

	if err := injectAlphaLauncher(mountpoint); err != nil {
		return "", helpers.WrapError(ErrInjectAlphaScript, err)
	}

	if err := sys.Umount(mountpoint); err != nil {
		return "", err
	}

	log.Infof("build/%s: rootfs successfully exported from docker", buildId)

	return rootfs, nil
}

// getImageSize returns the size of an image, in MB
func (b *Builder) getImageSize(ctx context.Context, tag string) int64 {
	image, _, err := b.client.ImageInspectWithRaw(ctx, tag)
	if err != nil {
		// We don't want to crash, so we return a default image size
		// in case of an error as occured
		return 500
	}

	// We add 5% of the image size to be sure everything will work
	return int64(math.Round(float64(image.Size/1000000) * 1.05))
}

// injectAlphaLauncher will add the required script into the rootfs to start Alpha when the VM will start
func injectAlphaLauncher(workingDirectory string) error {
	initScript := path.Join(workingDirectory, "sbin", "init")

	// I don't know why we need to keep it, but
	// if we remove it, the init is never launched.
	if err := os.Rename(initScript, path.Join(workingDirectory, "sbin", "init.old")); err != nil {
		return err
	}

	// Create and write the file /sbin/init that contains our
	// init script that will be launched on VM boot
	f, err := os.OpenFile(initScript, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write([]byte(alphaInitScript)); err != nil {
		return err
	}

	return nil
}

// makeExt4FS create rootfs.ext4 with a mountpoint on disk
func makeExt4FS(workingDirectory string, size int64) (string, string, error) {
	ext4 := path.Join(workingDirectory, "rootfs.ext4")

	// create the mount point on disk
	mountpoint := path.Join(workingDirectory, "mount")
	if err := os.MkdirAll(mountpoint, 0755); err != nil {
		return "", "", err
	}

	if err := sys.DD(ext4, size); err != nil {
		return "", "", err
	}

	if err := sys.MkfsExt4(ext4); err != nil {
		return "", "", err
	}

	if err := sys.Mount(ext4, mountpoint); err != nil {
		return "", "", err
	}

	return ext4, mountpoint, nil
}

// downloadRuntimes will clone the runtime repository to the local disk.
// It will be executed only once during the initialization of the builder.
func downloadRuntimes() (string, error) {
	log.Infof("downloading runtimes from %s into %s", runtimeEndpoint, destination)

	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		log.Debugf("trying to remove %s as it is not empty", destination)

		if err := os.RemoveAll(destination); err != nil {
			return "", err
		}
	}

	if err := getter.Get(destination, runtimeEndpoint); err != nil {
		return "", helpers.WrapError(ErrFailedToDownloadRuntimes, err)
	}

	return path.Join(destination, "template"), nil
}

// streamDockerLogs is a little helper function that help to stream the docker engine
// logs during an image build or a container run. The logs will be displayed to the console
// only if TRACE level is enabled.
func streamDockerLogs(buildId string, r io.Reader) {
	var line string

	// This type should not be used outside this function
	type buildLog struct {
		Message string `json:"stream"`
	}

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line = sc.Text()

		l := &buildLog{}
		json.Unmarshal([]byte(line), l)

		if l.Message == "" || l.Message == "\n" {
			continue
		}

		log.Tracef("build/%s: %s", buildId, strings.TrimSpace(l.Message))
	}
}
