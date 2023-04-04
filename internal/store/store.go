package store

type Store interface {
	Get(col, key string) (string, error)
	Set(col, key, value string) error
}
