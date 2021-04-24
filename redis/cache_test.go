package redis

import (
	"fmt"
	"testing"
)

func TestRedisCache(t *testing.T) {
	c := NewCache("localhost", 6379)
	err := c.Set("hello", []byte("world123"), 1500000)
	fmt.Println(err)
	if err != nil {
		t.FailNow()
	}
}

func TestGetMissing(t *testing.T) {
	c := NewCache("localhost", 6379)
	val, err := c.Get("hello")
	t.Log("err", err)
	t.Log("val", val)

	if err != nil {
		t.FailNow()
	}
}

func TestGet(t *testing.T) {
	initialVal := "world123"
	c := NewCache("localhost", 6379)
	c.Set("hello", []byte(initialVal), 1500000)
	val, err := c.Get("hello")
	t.Log("err", err)
	t.Log("val", val)

	if err != nil {
		t.FailNow()
	} else if string(val) != initialVal {
		t.Logf("expected value of %q, received value %q", initialVal, string(val))
		t.FailNow()
	}
}

func TestDelete(t *testing.T) {
	c := NewCache("localhost", 6379)
	err := c.Delete("hello")
	t.Log("err", err)
	if err != nil {
		t.FailNow()
	}
}

func TestDeleteFail(t *testing.T) {
	c := NewCache("localhostddddd", 6379)
	err := c.Delete("hello")
	t.Log("err", err)
	if err != nil {
		t.FailNow()
	}
}
