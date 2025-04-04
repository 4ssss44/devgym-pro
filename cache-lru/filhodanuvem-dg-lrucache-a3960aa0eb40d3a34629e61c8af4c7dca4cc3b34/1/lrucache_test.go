package lru1

import (
	"testing"
)

func TestReplaceCorrectKey(t *testing.T) {
	c := New(3)
	c.Set("a", 3)
	c.Set("b", 2)
	c.Set("c", 1)
	c.Get("a")
	c.Get("a")
	c.Get("c")
	c.Set("d", 4)

	removedKey := "b"

	for i := range c.items {
		if c.items[i].key == removedKey {
			t.Errorf("expected key %s to be removed and it was not", removedKey)
		}
	}

	if len(c.items) != 3 {
		t.Errorf("expected to have 3 items found %d", len(c.items))
	}
}

func TestGetNotfoundValue(t *testing.T) {
	c := New(3)
	c.Set("a", 3)
	c.Set("b", 2)
	c.Set("c", 1)
	if value := c.Get("a"); value != 3 {
		t.Errorf("expected to get 3, got %v", value)
	}
	if value := c.Get("a"); value != 3 {
		t.Errorf("expected to get 3, got %v", value)
	}
	if value := c.Get("c"); value != 1 {
		t.Errorf("expected to get 1, got %v", value)
	}
	c.Set("d", 4)
	if value := c.Get("d"); value != 4 {
		t.Errorf("expected to get d, got %v", value)
	}
}

func TestGetCorrectValues(t *testing.T) {
	c := New(3)
	c.Set("a", 3)
	c.Set("b", 2)
	c.Set("c", 1)
	if value := c.Get("x"); value != -1 {
		t.Errorf("expected to get -1, got %v", value)
	}
}
