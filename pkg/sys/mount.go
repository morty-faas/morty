package sys

import (
	log "github.com/sirupsen/logrus"
)

const (
	mountCmd  = "mount"
	umountCmd = "umount"
)

// Mount execute a syscall to the mount command to
// create a mount point on the host.
// See: https://man7.org/linux/man-pages/man8/mount.8.html
func Mount(source, target string) error {
	if err := Exec(mountCmd, source, target); err != nil {
		log.Errorf("failed to mount '%s' into '%s': %v", source, target, err)
		return err
	}
	return nil
}

// Umount execute a syscall to umount the target
// See: https://man7.org/linux/man-pages/man2/umount.2.html
func Umount(target string) error {
	if err := Exec(umountCmd, target); err != nil {
		log.Errorf("failed to umount '%s': %v", target, err)
		return err
	}
	return nil
}
