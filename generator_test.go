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
	client, err := NewGeneratorClient("test", LifeCycle(Second))
	assert.Nil(t, err)
	for i := 0; i < 100; i++ {
		id, err := client.GeneratorId()
		assert.Nil(t, err)
		fmt.Println(time.Now().Format("2006-01-02T15:04:05"), id)
		time.Sleep(500 * time.Millisecond)
	}
	iterator, err := client.cacheClient.Iterator()
	assert.Nil(t, err)
	for iterator.SetNext() {
		info, err := iterator.Value()
		assert.Nil(t, err)
		fmt.Println(fmt.Sprintf("key : %s ; value : %s", info.Key(), string(info.Value())))
	}
}

func TestGeneratorClient_GeneratorGroupId(t *testing.T) {
	client, err := NewGeneratorClient("test", Prefix("GT"), GroupLength(3), SequenceFormat("%02d"))
	assert.Nil(t, err)
	for i := 0; i < 100; i++ {
		id, err := client.GeneratorGroupId("group")
		assert.Nil(t, err)
		fmt.Println(id)
	}
}

func TestGeneratorClient_GeneratorGroupIdInstance(t *testing.T) {
	client, err := NewGeneratorClient(
		"test",
		Prefix("GT"),
		GroupLength(3),
		SequenceFormat("%02d"),
		Instance("8"),
		LifeCycle(Second),
	)
	assert.Nil(t, err)
	for i := 0; i < 500; i++ {
		id, err := client.GeneratorGroupId("group")
		assert.Nil(t, err)
		fmt.Println(id)
		time.Sleep(500 * time.Millisecond)
	}
}
