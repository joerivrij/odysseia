package cache

import (
	"github.com/dgraph-io/badger/v3"
)

func (b *Badger) Close() error {
	return b.db.Close()
}

func (b *Badger) Set(key, value string) error {
	txn := b.db.NewTransaction(true)
	defer txn.Discard()

	err := txn.Set([]byte(key), []byte(value))
	if err != nil {
		return err
	}

	if err := txn.Commit(); err != nil {
		return err
	}

	return nil
}

func (b *Badger) Read(key string) ([]byte, error) {
	var copiedValue []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			copiedValue = append([]byte{}, val...)
			return nil
		})
		return err
	})

	return copiedValue, err
}
