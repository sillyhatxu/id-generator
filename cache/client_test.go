package cache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type UserInfo struct {
	Id                  string    `json:"id" mapstructure:"id"`
	MobileNumber        string    `json:"mobile_number" mapstructure:"mobile_number"`
	Name                string    `json:"Name" mapstructure:"Name"`
	Paid                bool      `json:"Paid" mapstructure:"Paid"`
	FirstActionDeviceId string    `json:"first_action_device_id" mapstructure:"first_action_device_id"`
	TestNumber          int       `json:"test_number" mapstructure:"test_number"`
	TestNumber64        int64     `json:"test_number_64" mapstructure:"test_number_64"`
	TestDate            time.Time `json:"test_date" mapstructure:"test_date"`
	Member              *UserInfo `json:"member" mapstructure:"member"`
}

func TestCacheGetNil(t *testing.T) {
	client, err := NewCacheClient()
	if err != nil {
		panic(err)
	}
	test, err := client.Get("test1")
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, test, "")

	testsrc, err := client.GetSrc("test1")
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, testsrc, "")

	var user *UserInfo
	err = client.GetObj("test1", &user)
	if err != nil {
		panic(err)
	}
	assert.Nil(t, user)

	var resultArray []UserInfo
	err = client.GetObj("array", &resultArray)
	if err != nil {
		panic(err)
	}
	assert.Nil(t, resultArray)
}

func TestCacheSetGet(t *testing.T) {
	client, err := NewCacheClient()
	if err != nil {
		panic(err)
	}
	err = client.Set("test1", []byte("testhaha"))
	if err != nil {
		panic(err)
	}
	testsrc, err := client.Get("test1")
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, testsrc, "testhaha")
}

func TestCacheSetGetSrc(t *testing.T) {
	client, err := NewCacheClient()
	if err != nil {
		panic(err)
	}
	err = client.SetSrc("test1", "testhaha")
	if err != nil {
		panic(err)
	}
	testsrc, err := client.GetSrc("test1")
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, testsrc, "testhaha")
}

func TestCacheSetGetObj(t *testing.T) {
	client, err := NewCacheClient()
	if err != nil {
		panic(err)
	}
	member := &UserInfo{Id: "ID_2222", MobileNumber: "m_555555", Name: "m_test", Paid: false, FirstActionDeviceId: "m_deviceid", TestNumber: 11, TestNumber64: 22}
	userinfo := &UserInfo{Id: "ID_1001", MobileNumber: "555555", Name: "test", Paid: true, FirstActionDeviceId: "deviceid", TestNumber: 10, TestNumber64: 64, Member: member}
	err = client.SetObj("userinfo", userinfo)
	if err != nil {
		panic(err)
	}
	var user *UserInfo
	err = client.GetObj("userinfo", &user)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, user.Id, "ID_1001")
	assert.EqualValues(t, user.Name, "test")
	assert.EqualValues(t, user.MobileNumber, "555555")
	assert.EqualValues(t, user.Member.Id, "ID_2222")
	assert.EqualValues(t, user.Member.Name, "m_test")
	assert.EqualValues(t, user.Member.MobileNumber, "m_555555")
	assert.EqualValues(t, user, userinfo)
}

func TestCacheSetGetArray(t *testing.T) {
	client, err := NewCacheClient()
	if err != nil {
		panic(err)
	}
	var array []UserInfo
	array = append(array, UserInfo{Id: "ID_1001", MobileNumber: "MOBILE_1001", Name: "NAME_1001"})
	array = append(array, UserInfo{Id: "ID_1002", MobileNumber: "MOBILE_1002", Name: "NAME_1002"})
	array = append(array, UserInfo{Id: "ID_1003", MobileNumber: "MOBILE_1003", Name: "NAME_1003"})
	array = append(array, UserInfo{Id: "ID_1004", MobileNumber: "MOBILE_1004", Name: "NAME_1004"})
	array = append(array, UserInfo{Id: "ID_1005", MobileNumber: "MOBILE_1005", Name: "NAME_1005"})
	err = client.SetObj("array", array)
	if err != nil {
		panic(err)
	}

	var resultArray []UserInfo
	err = client.GetObj("array", &resultArray)
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, len(resultArray), len(array))
	assert.EqualValues(t, resultArray, array)
}

func TestExpired(t *testing.T) {
	client, err := NewCacheClient(LifeWindow(1*time.Second), CleanWindow(2*time.Second))
	if err != nil {
		panic(err)
	}
	err = client.SetSrc("test", "this is value")
	if err != nil {
		panic(err)
	}
	testsrc, err := client.GetSrc("test")
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, testsrc, "this is value")
	time.Sleep(3 * time.Second)
	testsrc, err = client.GetSrc("test")
	if err != nil {
		panic(err)
	}
	assert.EqualValues(t, testsrc, "")
}

func TestNoExpired(t *testing.T) {
	client, err := NewCacheClient(LifeWindow(1*time.Second), CleanWindow(-1*time.Second))
	if err != nil {
		panic(err)
	}
	err = client.SetSrc("test", "this is value")
	if err != nil {
		panic(err)
	}
	for {
		testsrc, err := client.GetSrc("test")
		if err != nil {
			panic(err)
		}
		fmt.Println(testsrc == "this is value")
		time.Sleep(5 * time.Second)
	}
}

func TestIncrementInt64(t *testing.T) {
	client, err := NewCacheClient()
	if err != nil {
		panic(err)
	}
	for i := 1; i <= 10000; i++ {
		index, err := client.IncrementInt("test")
		if err != nil {
			panic(err)
		}
		assert.EqualValues(t, index, i)
	}
}

func TestIterator(t *testing.T) {
	client, err := NewCacheClient()
	if err != nil {
		panic(err)
	}
	err = client.Set("test1", []byte("test1 value"))
	assert.Nil(t, err)
	err = client.SetSrc("test2", "test2 value")
	assert.Nil(t, err)
	var array []UserInfo
	array = append(array, UserInfo{Id: "ID_1001", MobileNumber: "MOBILE_1001", Name: "NAME_1001"})
	array = append(array, UserInfo{Id: "ID_1002", MobileNumber: "MOBILE_1002", Name: "NAME_1002"})
	err = client.SetObj("array", array)
	assert.Nil(t, err)

	_, err = client.IncrementInt("test3")
	assert.Nil(t, err)
	_, err = client.IncrementInt("test3")
	assert.Nil(t, err)
	_, err = client.IncrementInt("test3")
	assert.Nil(t, err)
	_, err = client.IncrementInt("test3")
	assert.Nil(t, err)

	iterator, err := client.Iterator()
	assert.Nil(t, err)
	for iterator.SetNext() {
		info, err := iterator.Value()
		assert.Nil(t, err)
		fmt.Println(fmt.Sprintf("key : %s ; value : %s", info.Key(), string(info.Value())))
	}
}
