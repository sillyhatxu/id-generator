package idgenerator

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestTimeInMillis(t *testing.T) {
	for {
		fmt.Println(fmt.Sprintf("%s", strconv.FormatInt(time.Now().Unix()/60, 10)))
		time.Sleep(time.Second)
	}
}

func TestHex(t *testing.T) {
	test1 := "Hello Gopher!"
	src := []byte(test1)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	fmt.Printf("%s ==> %s \n", test1, dst)

	str := "Hello from ADMFactory.com"
	hx := hex.EncodeToString([]byte(str))
	fmt.Printf("%s ==> %s \n", str, hx)
}

func TestUUIDCheckVersion1(t *testing.T) {
	for i := 0; i < 100; i++ {
		generatorUUID, err := uuid.NewUUID()
		if err != nil {
			t.Fatalf("could not create UUID: %v", err)
		}
		fmt.Println(generatorUUID.String())
		time.Sleep(200 * time.Millisecond)
	}
}

func TestUUIDCheckVersion4(t *testing.T) {
	for i := 0; i < 100; i++ {
		generatorUUID := uuid.New()
		fmt.Println(generatorUUID.String())
		time.Sleep(200 * time.Millisecond)
	}
}

func TestUUID(t *testing.T) {
	generatorUUID, err := uuid.NewUUID()
	if err != nil {
		t.Fatalf("could not create UUID: %v", err)
	}
	fmt.Println(generatorUUID.String())
	id := strings.ToUpper(strings.ReplaceAll(generatorUUID.String(), "-", ""))
	fmt.Println(id)
	fmt.Println(hex.EncodeToString([]byte(id)))
}

func TestGeneratorClient_GeneratorGroupIdDefault(t *testing.T) {
	client, err := NewGeneratorClient("test", Prefix("GT"))
	assert.Nil(t, err)
	for i := 0; i < 100; i++ {
		id, err := client.GeneratorGroupId("asdgadsfhgdhj")
		assert.Nil(t, err)
		fmt.Println(id)
	}
}

func TestGeneratorClient_GeneratorGroupId(t *testing.T) {
	client, err := NewGeneratorClient("test", Prefix("GT"), GroupLength(3), SequenceFormat("%03d"))
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
		SequenceFormat("%03d"),
		Instance("8"),
		LifeCycle(Minute),
	)
	assert.Nil(t, err)
	for i := 0; i < 500; i++ {
		id, err := client.GeneratorGroupId("group")
		assert.Nil(t, err)
		fmt.Println(id)
		time.Sleep(500 * time.Millisecond)
		//GT826176437022516
	}
}
