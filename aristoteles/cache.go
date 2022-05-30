package aristoteles

import (
	"github.com/odysseia/plato/cache"
)

func (c *Config) getBadgerClient() (cache.Client, error) {
	badger, err := cache.NewBadgerClient("")
	if err != nil {
		return nil, err
	}

	return badger, nil
}
