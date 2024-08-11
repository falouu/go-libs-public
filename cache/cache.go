package cache

import (
	"encoding/json"
	"github.com/falouu/go-libs-public/b"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

type FetcherFn func() (result interface{}, err error, shouldCache bool)

type Cache interface {
	Get(key string, fetcher FetcherFn, resultTypeHint interface{}) (result interface{}, err error)
}

type cacheImpl struct {
	useCache bool
	dir      string
	log      *logrus.Entry
}

func NewCache(useCache bool, dir string) Cache {

	return &cacheImpl{useCache: useCache, dir: dir, log: logrus.WithField("src", "cache")}
}

func (c *cacheImpl) Get(key string, fetcher FetcherFn, ptrToCached interface{}) (result interface{}, err error) {
	cachedExists := false
	if c.useCache {
		cachedExists, err = c.tryGetCached(key, ptrToCached)
		if err != nil {
			return nil, err
		}
	}

	if cachedExists {
		c.log.Debug("Using cached value for key ", key)
		result = reflect.Indirect(reflect.ValueOf(ptrToCached)).Interface()
	} else {
		c.log.Debug("Calculating value for key ", key)
		var fetcherError error
		var shouldCache bool
		result, fetcherError, shouldCache = fetcher()
		if fetcherError != nil {
			err = fetcherError
			return result, err
		}
		if shouldCache {
			warn := c.cache(key, result)
			if warn != nil {
				c.log.Error(b.Wrap(warn, "failed to cache entry '%v'", key))
			}
		}
	}

	return
}

func (c *cacheImpl) tryGetCached(key string, ptrToCached interface{}) (bool, error) {

	bytes, err := ioutil.ReadFile(c.path(key))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	err = json.Unmarshal(bytes, ptrToCached)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *cacheImpl) path(key string) string {

	return filepath.Join(c.dir, key+".json")
}

func (c *cacheImpl) cache(key string, val interface{}) error {
	err := os.MkdirAll(c.dir, 0700)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(val, "", " ")
	if err != nil {
		return err
	}

	path := c.path(key)
	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bytes, 0600)
}
