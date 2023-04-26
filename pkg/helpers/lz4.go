package helpers

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/pierrec/lz4"
	"github.com/sirupsen/logrus"
)

func CompressLZ4(file, dest string) error {
	logrus.Tracef("helpers/lz4: compress '%s' into '%s'", file, dest)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	var buff bytes.Buffer
	w := lz4.NewWriter(&buff)
	defer w.Close()

	w.Write(data)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(dest, buff.Bytes(), 0644)
}
