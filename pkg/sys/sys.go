package sys

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

const (
	mkfsExt4Cmd = "mkfs.ext4"
	ddCmd       = "dd"
)

// MkfsExt4 transform the directory given in parameters into
// and ext4 file system.
func MkfsExt4(directory string) error {
	if err := Exec(mkfsExt4Cmd, directory); err != nil {
		log.Errorf("failed to create ext4 FS into '%s': %v", directory, err)
		return err
	}
	return nil
}

func DD(path string, size int64) error {
	if err := Exec(ddCmd, "if=/dev/zero", fmt.Sprintf("of=%s", path), "bs=1M", fmt.Sprintf("count=%d", size)); err != nil {
		log.Errorf("failed to create device at path '%s' with size %dMB: %v", path, size, err)
		return err
	}
	return nil
}
