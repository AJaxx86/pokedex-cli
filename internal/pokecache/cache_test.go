package pokecache

import (
	"testing"
	"fmt"
	"time"
	"sync"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			defer cache.Stop()
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}


func TestAddGetErrors(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)
	defer cache.Stop()

	cases := []struct {
		name    string
		key     string
		val     []byte
		wantErr bool
	}{
		{
			name:    "empty key",
			key:     "",
			val:     []byte("testdata"),
			wantErr: true,
		},
		{
			name:    "empty value",
			key:     "https://example.com",
			val:     []byte(""),
			wantErr: true,
		},
		{
			name:    "valid key and value",
			key:     "https://example.com",
			val:     []byte("testdata"),
			wantErr: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := cache.Add(c.key, c.val)
			if (err != nil) != c.wantErr {
				t.Errorf("Add(%q, %q) error = %v, wantErr = %v", c.key, c.val, err, c.wantErr)
			}
		})
	}
}

func TestGetMissingKey(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)
	defer cache.Stop()

	val, ok := cache.Get("https://does-not-exist.com")
	if ok {
		t.Errorf("expected to not find key, got value %q", val)
		return
	}
	if val != nil {
		t.Errorf("expected nil value for missing key, got %q", val)
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5 * time.Millisecond
	cache := NewCache(baseTime)
	defer cache.Stop()
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Fatalf("expected to find key")
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestAddOverwrite(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)
	defer cache.Stop()

	cache.Add("https://example.com", []byte("first"))
	cache.Add("https://example.com", []byte("second"))

	val, ok := cache.Get("https://example.com")
	if !ok {
		t.Fatalf("expected to find key")
	}
	if string(val) != "second" {
		t.Errorf("expected overwritten value 'second', got %q", val)
	}
}

func TestConcurrentAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)
	defer cache.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", n)
			cache.Add(key, []byte("val"))
			cache.Get(key)
		}(i)
	}
	wg.Wait()

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		if _, ok := cache.Get(key); !ok {
			t.Errorf("expected to find key %q", key)
		}
	}
}
