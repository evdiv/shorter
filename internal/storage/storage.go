package storage

type Storer interface {
	Set(url string) string
	Get(key string) string
	IsAvailable() bool
}
