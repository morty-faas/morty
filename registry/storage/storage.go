package storage

import "io"

type Storage interface {
	// GetDownloadLink will request a download link for the given function key in the storage.
	GetDownloadLink(key string) (string, string, error)

	// PutFile insert a new file at the given key in the storage
	PutFile(key string, body io.Reader) error

	Healthcheck() error
}
