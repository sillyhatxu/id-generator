package idgenerator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestTimeInMillis(t *testing.T) {
	for {
		fmt.Println(fmt.Sprintf("%s", strconv.FormatInt(time.Now().Unix()/60, 10)))
		time.Sleep(time.Second)
	}
}
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
		LifeCycle(Minute),
	)
	for i := 0; i < 10000; i++ {
		id, err := client.GeneratorGroupId("group")
		assert.Nil(t, err)
		fmt.Println(id)
		time.Sleep(1 * time.Second)
	}
}
