package cache

import (
	"encoding/json"
	"fmt"
	"github.com/allegro/bigcache"
	"strconv"
	"sync"
	"time"
)

const (
	defaultShards             = 2048
	defaultLifeWindow         = 24 * time.Hour
	defaultCleanWindow        = 48 * time.Hour
	defaultMaxEntriesInWindow = 1000 * 10 * 60
	defaultMaxEntrySize       = 500
	defaultVerbose            = true
	defaultHardMaxCacheSize   = 8192
)

type CacheClient struct {
	client *bigcache.BigCache
	config *bigcache.Config
	mu     sync.Mutex
}

func NewCacheClient(opts ...Option) (*CacheClient, error) {
	//default
	config := &Config{
		shards:             defaultShards,
		lifeWindow:         defaultLifeWindow,
		cleanWindow:        defaultCleanWindow,
		maxEntriesInWindow: defaultMaxEntriesInWindow,
		maxEntrySize:       defaultMaxEntrySize,
		verbose:            defaultVerbose,
		hardMaxCacheSize:   defaultHardMaxCacheSize,
		logger:             bigcache.DefaultLogger(),
	}
	for _, opt := range opts {
		opt(config)
	}

	bigcacheConfig := bigcache.Config{
		Shards:             config.shards,
		LifeWindow:         config.lifeWindow,
		CleanWindow:        config.cleanWindow,
		MaxEntriesInWindow: config.maxEntriesInWindow,
		MaxEntrySize:       config.maxEntrySize,
		Verbose:            config.verbose,
		HardMaxCacheSize:   config.hardMaxCacheSize,
		Logger:             config.logger,
	}
	cache, err := bigcache.NewBigCache(bigcacheConfig)
	if err != nil {
		return nil, err
	}
	return &CacheClient{
		client: cache,
		config: &bigcacheConfig,
	}, nil
}

func (cc CacheClient) getClient() (*bigcache.BigCache, error) {
	if cc.client == nil {
		c, err := bigcache.NewBigCache(*cc.config)
		if err != nil {
			return nil, err
		}
		cc.client = c
		return cc.client, nil
	}
	return cc.client, nil
}

func (cc CacheClient) SetSrc(key string, value string) error {
	return cc.Set(key, []byte(value))
}

func (cc CacheClient) GetSrc(key string) (string, error) {
	body, err := cc.Get(key)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (cc CacheClient) SetObj(key string, value interface{}) error {
	bodyJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("value to byte error. %v", err)
	}
	return cc.Set(key, bodyJSON)
}

func (cc CacheClient) GetObj(key string, value interface{}) error {
	body, err := cc.Get(key)
	if err != nil {
		return err
	}
	if body == nil {
		return nil
	}
	err = json.Unmarshal(body, &value)
	if err != nil {
		return err
	}
	return nil
}

func (cc CacheClient) Get(key string) ([]byte, error) {
	client, err := cc.getClient()
	if err != nil {
		return nil, err
	}
	body, err := client.Get(key)
	if err != nil && err == bigcache.ErrEntryNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return body, nil
}

func (cc CacheClient) Set(key string, value []byte) error {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	client, err := cc.getClient()
	if err != nil {
		return err
	}
	return client.Set(key, value)
}

func (cc CacheClient) Delete(key string) error {
	client, err := cc.getClient()
	if err != nil {
		return err
	}
	return client.Delete(key)
}

func (cc CacheClient) Exist(key string) (bool, error) {
	body, err := cc.Get(key)
	if err != nil {
		return false, err
	}
	return body != nil, nil
}

func (cc CacheClient) Iterator() (*bigcache.EntryInfoIterator, error) {
	client, err := cc.getClient()
	if err != nil {
		return nil, err
	}
	return client.Iterator(), nil
}

func (cc CacheClient) IncrementInt(key string) (int, error) {
	index, err := cc.IncrementInt64(key)
	if err != nil {
		return 0, nil
	}
	return int(index), nil
}

func (cc CacheClient) IncrementInt64(key string) (int64, error) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	client, err := cc.getClient()
	if err != nil {
		return 0, err
	}
	exist := true
	body, err := client.Get(key)
	if err != nil && err == bigcache.ErrEntryNotFound {
		exist = false
	} else if err != nil {
		return 0, err
	}
	var index int64 = 0
	if exist {
		n, err := strconv.ParseInt(string(body), 10, 64)
		if err != nil {
			return 0, err
		}
		index = n
	}
	index++
	indexSrc := strconv.FormatInt(index, 10)
	err = client.Set(key, []byte(indexSrc))
	if err != nil {
		return 0, err
	}
	return index, nil
}
