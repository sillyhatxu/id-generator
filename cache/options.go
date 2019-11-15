package cache

import (
	"github.com/allegro/bigcache"
	"time"
)

type Config struct {
	shards             int
	lifeWindow         time.Duration
	cleanWindow        time.Duration
	maxEntriesInWindow int
	maxEntrySize       int
	verbose            bool
	hardMaxCacheSize   int
	logger             bigcache.Logger
}

type Option func(*Config)

func Shards(shards int) Option {
	return func(c *Config) {
		c.shards = shards
	}
}

func LifeWindow(lifeWindow time.Duration) Option {
	return func(c *Config) {
		c.lifeWindow = lifeWindow
	}
}

func CleanWindow(cleanWindow time.Duration) Option {
	return func(c *Config) {
		c.cleanWindow = cleanWindow
	}
}

func MaxEntriesInWindow(maxEntriesInWindow int) Option {
	return func(c *Config) {
		c.maxEntriesInWindow = maxEntriesInWindow
	}
}

func MaxEntrySize(maxEntrySize int) Option {
	return func(c *Config) {
		c.maxEntrySize = maxEntrySize
	}
}

func Verbose(verbose bool) Option {
	return func(c *Config) {
		c.verbose = verbose
	}
}

func HardMaxCacheSize(hardMaxCacheSize int) Option {
	return func(c *Config) {
		c.hardMaxCacheSize = hardMaxCacheSize
	}
}

func Logger(logger bigcache.Logger) Option {
	return func(c *Config) {
		c.logger = logger
	}
}
