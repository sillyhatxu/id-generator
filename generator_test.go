package idgenerator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGeneratorClient_GeneratorId(t *testing.T) {
	client := NewGeneratorClient("test")
	for i := 0; i < 100; i++ {
		id, err := client.GeneratorId()
		assert.Nil(t, err)
		fmt.Println(id)
	}
}

func TestGeneratorClient_GeneratorGroupId(t *testing.T) {
	client := NewGeneratorClient("test", Prefix("GT"), GroupLength(3), SequenceFormat("%02d"))
	for i := 0; i < 100; i++ {
		id, err := client.GeneratorGroupId("group")
		assert.Nil(t, err)
		fmt.Println(id)
	}
}

func TestGeneratorClient_GeneratorGroupIdInstance(t *testing.T) {
	client := NewGeneratorClient(
		"test",
		Prefix("GT"),
		GroupLength(3),
		SequenceFormat("%02d"),
		Instance("8"),
		LifeCycle(5*time.Second),
	)
	for i := 0; i < 100; i++ {
		id, err := client.GeneratorGroupId("group")
		assert.Nil(t, err)
		fmt.Println(id)
	}
}
