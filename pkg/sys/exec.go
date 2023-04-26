package sys

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func Exec(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	log.Tracef("sys/exec: %s", cmd.String())
	return cmd.Run()
}
