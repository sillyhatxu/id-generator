package idgenerator

import (
	"fmt"
	"github.com/sillyhatxu/gocache-client"
	"hash/fnv"
	"strconv"
	"sync"
	"time"
)

const (
	defaultGroupLength    int    = 2
	defaultSequenceFormat string = "%02d"
	defaultLifeCycle             = Minute
)

type GeneratorClient struct {
	key         string
	config      *Config
	cacheClient *cacheclient.CacheClient
	mu          sync.Mutex
}

func NewGeneratorClient(key string, opts ...Option) (*GeneratorClient, error) {
	//default
	config := &Config{
		Prefix:         "",
		GroupLength:    defaultGroupLength,
		SequenceFormat: defaultSequenceFormat,
		LifeCycle:      defaultLifeCycle,
	}
	for _, opt := range opts {
		opt(config)
	}
	l, c := getLifeWindowAndCleanWindow(config.LifeCycle)
	client, err := cacheclient.NewCacheClient(cacheclient.LifeWindow(l), cacheclient.CleanWindow(c))
	if err != nil {
		return nil, err
	}
	return &GeneratorClient{
		key:         key,
		config:      config,
		cacheClient: client,
	}, nil
}

func getLifeWindowAndCleanWindow(lifeCycle LifeCycleType) (time.Duration, time.Duration) {
	if lifeCycle == Minute {
		return 1 * time.Minute, 2 * time.Minute
	} else if lifeCycle == Hour {
		return 1 * time.Hour, 2 * time.Hour
	} else {
		return 1 * time.Second, 2 * time.Second
	}
}

func (gc GeneratorClient) validate() error {
	if gc.key == "" {
		return fmt.Errorf("key cannot empty")
	}
	if gc.config == nil {
		return fmt.Errorf("config is nil")
	}
	if gc.cacheClient == nil {
		return fmt.Errorf("cache client is nil")
	}
	return nil
}

func (gc GeneratorClient) GeneratorId() (string, error) {
	return gc.GeneratorGroupId("")
}

func (gc GeneratorClient) GeneratorGroupId(src string) (string, error) {
	err := gc.validate()
	if err != nil {
		return "", err
	}
	gc.mu.Lock()
	defer gc.mu.Unlock()
	group, err := gc.formatGroup(src)
	if err != nil {
		return "", err
	}
	sequence, err := gc.getSequence(group)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s%s%s%s", gc.config.Prefix, gc.config.Instance, gc.getTimeInMillis(), sequence, group), nil
}

func (gc GeneratorClient) formatGroup(src string) (string, error) {
	if src == "" {
		return "", nil
	}
	hashSrc, err := hash(src)
	if err != nil {
		return "", err
	}
	formatUintSrc := strconv.FormatUint(hashSrc, 10)
	if len(formatUintSrc) > gc.config.GroupLength {
		return string(formatUintSrc[len(formatUintSrc)-gc.config.GroupLength:]), nil
	} else {
		return formatUintSrc, nil
	}
}

func (gc GeneratorClient) getSequence(group string) (string, error) {
	key := fmt.Sprintf("%s_%s_%s", gc.key, group, gc.getKeySuffix())
	sequence, err := gc.cacheClient.IncrementInt64(key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(gc.config.SequenceFormat, sequence), nil
}

func (gc GeneratorClient) getKeySuffix() string {
	hr, min, sec := time.Now().Clock()
	if gc.config.LifeCycle == Minute {
		return fmt.Sprintf("%d_%d", hr, min)
	} else if gc.config.LifeCycle == Hour {
		return fmt.Sprintf("%d", hr)
	} else {
		return fmt.Sprintf("%d_%d_%d", hr, min, sec)
	}
}

func (gc GeneratorClient) getTimeInMillis() string {
	return strconv.FormatInt(time.Now().Unix()/getLifeCycleNumber(gc.config.LifeCycle), 10)
}

func getLifeCycleNumber(lifeCycle LifeCycleType) int64 {
	if lifeCycle == Minute {
		return 60
	} else if lifeCycle == Hour {
		return 60 * 60
	} else {
		return 1
	}
}

func hash(s string) (uint64, error) {
	h := fnv.New64a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}
