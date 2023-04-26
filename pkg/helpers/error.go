package helpers

import (
	log "github.com/sirupsen/logrus"
)

// Wrap log a custom message for each error that is wrapped into another
func WrapError(custom, original error) error {
	log.Errorf("%v:%v", custom, original)
	return original
}
