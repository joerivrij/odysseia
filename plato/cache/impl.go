package cache

import (
	"github.com/dgraph-io/badger/v3"
)

type Client interface {
	Close() error
	Set(key, value string) error
	Read(key string) ([]byte, error)
}

type Badger struct {
	db *badger.DB
}

func NewInMemoryBadgerClient() (Client, error) {
	opt := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	badgerDb := Badger{db: db}

	return &badgerDb, nil
}

func NewBadgerClient(badgerPath string) (Client, error) {
	if badgerPath == "" {
		badgerPath = "/tmp/badger"
	}
	db, err := badger.Open(badger.DefaultOptions(badgerPath))
	if err != nil {
		return nil, err
	}

	badgerDb := Badger{db: db}

	return &badgerDb, nil
}
