package cache

import (
	"bytes"
	"encoding/json"
	"github.com/falouu/go-libs-public/b"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
)

type fetchFn[T any] func() (result *T, err error, shouldCacheSuccess bool)

func Fetch[T any](fn fetchFn[T]) *Fetcher[T] {
	return &Fetcher[T]{
		fn: fn,
	}
}

type Fetcher[T any] struct {
	fn     fetchFn[T]
	Result T
}

func (f *Fetcher[T]) fetcher() fetchFnAny {
	return func() (result any, err error, shouldcache bool) {
		return f.fn()
	}
}

func (f *Fetcher[T]) setResultPtr(resPtr any) {
	resPtrCasted := resPtr.(*T)
	f.Result = *resPtrCasted
}

func (f *Fetcher[T]) resultPtr() any {
	return &f.Result
}

type fetchFnAny = func() (resultPtr any, err error, shouldcache bool)

type fetcher interface {
	fetcher() fetchFnAny
	setResultPtr(res any)
	resultPtr() any
}

type Cache interface {
	Get(key string, fetcher fetcher) error
}

type cacheImpl struct {
	useCache bool
	dir      string
	log      *logrus.Entry
}

func NewCache(useCache bool, dir string) Cache {

	return &cacheImpl{useCache: useCache, dir: dir, log: logrus.WithField("src", "cache")}
}

func (c *cacheImpl) Get(key string, fetcher fetcher) (err error) {
	cachedExists := false
	if c.useCache {
		cachedExists, err = c.tryGetCached(key, fetcher)
		if err != nil {
			return err
		}
	}

	if cachedExists {
		c.log.Debug("Using cached value for key ", key)
	} else {
		c.log.Debug("Calculating value for key ", key)
		resultPtr, fetcherError, shouldCache := fetcher.fetcher()()
		if fetcherError != nil {
			err = fetcherError
			return err
		}
		fetcher.setResultPtr(resultPtr)

		if shouldCache {
			warn := c.cache(key, resultPtr)
			if warn != nil {
				c.log.Error(b.Wrap(warn, "failed to cache entry '%v'", key))
			}
		}
	}

	return
}

func (c *cacheImpl) tryGetCached(key string, fetcher fetcher) (bool, error) {

	bytes, err := ioutil.ReadFile(c.path(key))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	err = unmarshal(bytes, fetcher.resultPtr())
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *cacheImpl) path(key string) string {

	return filepath.Join(c.dir, key+".json")
}

func (c *cacheImpl) cache(key string, valPtr any) error {
	err := os.MkdirAll(c.dir, 0700)
	if err != nil {
		return err
	}

	bytes, err := marshal(valPtr)
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

func marshal(valPtr any) ([]byte, error) {
	if valBytes, ok := valPtr.(*[]byte); ok {
		dst := bytes.NewBuffer(make([]byte, 0, int(float64(len(*valBytes))*1.3)))
		err := json.Indent(dst, *valBytes, "", " ")
		return dst.Bytes(), err
	} else {
		return json.MarshalIndent(valPtr, "", " ")
	}
}

func unmarshal(data []byte, resultPtr any) error {
	if resultBytes, ok := resultPtr.(*[]byte); ok {
		*resultBytes = data
		return nil
	} else {
		return json.Unmarshal(data, resultPtr)
	}
}
