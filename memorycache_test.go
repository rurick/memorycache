package memorycache

import (
	"testing"
	"time"
)

const (
	testKey      string = "cache:test"
	testKeyEmpty string = "cache:empty"
	testValue    string = "Hello world"
)

func TestKeyGen(t *testing.T) {
	tests := []struct {
		name string
		args interface{}
		want string
	}{
		{
			"strings in key",
			[]string{"k1", "k2", "k3"},
			"eb98f46aca624b1e402947677c54d025dc463a67",
		}, {
			"int64 in key",
			[]int64{1, 2, 3},
			"6d780b01458b623aa5f77db71ac9a02ff1d5ecda",
		}, {
			"mixed key",
			[]interface{}{1, "k2", 3},
			"ebdf5f5fd2817d5534d668548298bd135ae3dd0e",
		}, {
			"mixed key",
			[]interface{}{1, "k2", 3, map[int]string{1: "a"}},
			"60b428a15bf1e9eca224ef24e37c38ec8d8f86f9",
		}, {
			"mixed key with nil",
			[]interface{}{1, "k2", nil, map[int]string{1: "a"}},
			"3b381ec1e7defe06ddf6a4eeb062211c579888bf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KeyGen(tt.args)
			if k != tt.want {
				t.Errorf("KeyGen(%v), want %v, got %v ", tt.args, tt.want, k)
			}
		})
	}
}

// Test_Get get cache by key
func Test_Get(t *testing.T) {
	cache := New(10*time.Minute, 1*time.Hour)
	cache.Set(testKey, testValue, 1*time.Minute)
	value, found := cache.Get(testKey)

	if value != testValue {
		t.Error("Error: ", "Set and Get not simple:", value, testValue)
	}

	if found != true {
		t.Error("Error: ", "Could not get cache")
	}

	value, found = cache.Get(testKeyEmpty)
	if value != nil || found != false {
		t.Error("Error: ", "Value does not exist and must be empty", value)
	}
}

// Test_Get get cache by key
func Test_GetWithoutOpt(t *testing.T) {
	cache := New(10*time.Minute, 1*time.Hour)
	cache.Set(testKey, testValue)
	value, found := cache.Get(testKey)

	if value != testValue {
		t.Error("Error: ", "Set and Get not simple:", value, testValue)
	}

	if found != true {
		t.Error("Error: ", "Could not get cache")
	}

	value, found = cache.Get(testKeyEmpty)
	if value != nil || found != false {
		t.Error("Error: ", "Value does not exist and must be empty", value)
	}
}

// Test_Delete delete cache by key
func Test_Delete(t *testing.T) {
	cache := New(10*time.Minute, 1*time.Hour)
	cache.Set(testKey, testValue, 1*time.Minute)
	err := cache.Delete(testKey)

	if err != nil {
		t.Error("Error: ", "Cache delete failed")
	}

	value, found := cache.Get(testKey)
	if found {
		t.Error("Error: ", "Should not be found because it was deleted")
	}

	if value != nil {
		t.Error("Error: ", "Value is not nil:", value)
	}

	err = cache.Delete(testKeyEmpty)
	if err == nil {
		t.Error("Error: ", "An empty cache should return an error")
	}

}
