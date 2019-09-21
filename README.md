# id-generator

## sample

[example](https://github.com/sillyhatxu/id-generator/blob/master/generator_test.go)

```
func TestGeneratorClient_GeneratorId(t *testing.T) {
	client :=.NewGeneratorClient("test")
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
```