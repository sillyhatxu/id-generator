package idgenerator

type LifeCycleType int

const (
	Second LifeCycleType = iota
	Minute
	Hour
)

type Config struct {
	Instance       string
	Prefix         string
	GroupLength    int
	SequenceFormat string
	LifeCycle      LifeCycleType
}

type Option func(*Config)

func Instance(instance string) Option {
	return func(c *Config) {
		c.Instance = instance
	}
}

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

func LifeCycle(lifeCycle LifeCycleType) Option {
	return func(c *Config) {
		c.LifeCycle = lifeCycle
	}
}
