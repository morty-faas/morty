package helpers

import (
	"bytes"

	"log"
	"os"

	"github.com/pierrec/lz4"
	"github.com/sirupsen/logrus"
)

func CompressLZ4(file, dest string) error {
	logrus.Tracef("helpers/lz4: compress '%s' into '%s'", file, dest)
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	var buff bytes.Buffer
	w := lz4.NewWriter(&buff)
	defer w.Close()

	if _, err := w.Write(data); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return os.WriteFile(dest, buff.Bytes(), 0644)
}
