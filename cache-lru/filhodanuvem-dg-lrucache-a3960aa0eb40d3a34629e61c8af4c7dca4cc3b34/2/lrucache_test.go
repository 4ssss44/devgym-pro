package lru2

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

	expectedValues := map[string]any{"a": 3, "c": 1, "d": 4}
	removedKey := "b"
	for i := range expectedValues {
		if _, ok := c.items[i]; !ok {
			t.Errorf("expected key %s, found nil. items: %+v", i, c.items)
			return
		}
		if c.items[i].value != expectedValues[i] {
			t.Errorf("Expected %v, got %v", expectedValues[i], c.items[i].value)
		}
	}

	if _, ok := c.items[removedKey]; ok {
		t.Errorf("expected key a to be removed and it was not")
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

func TestCapacityOneAlwaysReplace(t *testing.T) {
	c := New(1)

	criteria := []struct {
		items      map[string]any
		removedKey string
		addedKey   string
		addedValue any
	}{
		{map[string]any{"a": 3}, "", "a", 3},
		{map[string]any{"b": 2}, "a", "b", 2},
		{map[string]any{"c": 1}, "b", "c", 1},
		{map[string]any{"d": 4}, "c", "d", 4},
	}

	for _, criterion := range criteria {
		c.Set(criterion.addedKey, criterion.addedValue)

		for i := range criterion.items {
			if _, ok := c.items[i]; !ok {
				t.Errorf("expected key %s, found nil. items: %+v", i, c.items)
				return
			}
			if c.items[i].value != criterion.items[i] {
				t.Errorf("Expected %v, got %v", criterion.items[i], c.items[i].value)
			}
		}

		if _, ok := criterion.items[criterion.removedKey]; ok {
			t.Errorf("expected key a to be removed and it was not")
		}

		if len(c.items) != 1 {
			t.Errorf("expected to have 1 items, found %d", len(c.items))
		}
	}
}
