package storage

type FileStore interface {
	Upload([]byte) (string, error)
	Download(string) ([]byte, error)
}
