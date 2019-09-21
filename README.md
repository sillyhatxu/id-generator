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
```

> output

```
15690628250001
15690628250002
15690628250003
15690628250004
15690628250005
15690628250006
15690628250007
15690628250008
15690628250009
15690628250010
```


```
func TestGeneratorClient_GeneratorGroupId(t *testing.T) {
	client := NewGeneratorClient("test", Prefix("GT"), GroupLength(3), SequenceFormat("%02d"))
	for i := 0; i < 100; i++ {
		id, err := client.GeneratorGroupId("group")
		assert.Nil(t, err)
		fmt.Println(id)
	}
}
```

> output

```
GT156906285601516
GT156906285602516
GT156906285603516
GT156906285604516
GT156906285605516
GT156906285606516
GT156906285607516
GT156906285608516
```