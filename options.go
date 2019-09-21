package idgenerator

import (
	"time"
)

type Config struct {
	Prefix         string
	GroupLength    int
	SequenceFormat string
	LifeCycle      time.Duration
}

type Option func(*Config)

func Prefix(prefix string) Option {
	return func(c *Config) {
		c.Prefix = prefix
	}
}

func GroupLength(groupLength int) Option {
	return func(c *Config) {
		c.GroupLength = groupLength
	}
}

func SequenceFormat(sequenceFormat string) Option {
	return func(c *Config) {
		c.SequenceFormat = sequenceFormat
	}
}

func LifeCycle(lifeCycle time.Duration) Option {
	return func(c *Config) {
		c.LifeCycle = lifeCycle
	}
}
