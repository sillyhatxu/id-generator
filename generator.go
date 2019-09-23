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
	defaultSequenceFormat string = "%04d"
	defaultLifeCycle             = Minute
)

type GeneratorClient struct {
	key         string
	config      *Config
	cacheClient *client.CacheClient
	mu          sync.Mutex
}

func NewGeneratorClient(key string, opts ...Option) *GeneratorClient {
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
	return &GeneratorClient{
		key:         key,
		config:      config,
		cacheClient: client.NewCacheClient(),
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
	key := fmt.Sprintf("%s%s", gc.key, group)
	sequence, err := gc.cacheClient.IncrementInt64WithExpiration(key, time.Duration(gc.getLifeCycleNumber())*time.Second+500*time.Millisecond)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(gc.config.SequenceFormat, sequence), nil
}

func (gc GeneratorClient) getTimeInMillis() string {
	return strconv.FormatInt(time.Now().Unix()/gc.getLifeCycleNumber(), 10)
}

func (gc GeneratorClient) getLifeCycleNumber() int64 {
	if gc.config.LifeCycle == Minute {
		return 60
	} else if gc.config.LifeCycle == Hour {
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
