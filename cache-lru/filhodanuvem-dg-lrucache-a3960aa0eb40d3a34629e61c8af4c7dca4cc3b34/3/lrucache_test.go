package lru3

import (
	"fmt"
	"testing"
)

func TestReplaceCorrectKeyWhenKeyIsAccessed(t *testing.T) {
	c := New(3)
	c.Set("a", 3)
	c.Set("b", 2)
	c.Set("c", 1)
	c.Get("a")
	c.Get("c")
	c.Get("a")
	c.Get("b")
	c.Set("d", 4)

	expectedValues := map[string]any{"a": 3, "b": 2, "d": 4}
	removedKey := "c"
	for i := range expectedValues {
		if _, ok := c.dict[i]; !ok {
			t.Errorf("expected key %s, found nil.", i)
			for i := range c.dict {
				fmt.Println(i)
			}

			return
		}
		if c.dict[i].value != expectedValues[i] {
			t.Errorf("Expected %v, got %v", expectedValues[i], c.dict[i].value)
		}
	}

	if _, ok := c.dict[removedKey]; ok {
		t.Errorf("expected key a to be removed and it was not")
	}

	if len(c.dict) != 3 {
		t.Errorf("expected to have 3 dict found %d", len(c.dict))
	}
}

func TestReplaceCorrectKeyWhenKeyIsNeverAccessed(t *testing.T) {
	c := New(3)
	c.Set("a", 3)
	c.Set("b", 2)
	c.Set("c", 1)
	c.Get("c")
	c.Set("d", 4)

	expectedValues := map[string]any{"a": 3, "c": 1, "d": 4}
	removedKey := "b"
	for i := range expectedValues {
		if _, ok := c.dict[i]; !ok {
			for i := c.items; i != nil; i = i.next {
				fmt.Println(i.key)
			}
			t.Errorf("expected key %s, found nil.", i)

			return
		}
		if c.dict[i].value != expectedValues[i] {
			t.Errorf("Expected %v, got %v", expectedValues[i], c.dict[i].value)
		}
	}

	if _, ok := c.dict[removedKey]; ok {
		t.Errorf("expected key a to be removed and it was not")
	}

	if len(c.dict) != 3 {
		t.Errorf("expected to have 3 dict found %d", len(c.dict))
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
		dict       map[string]any
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

		for i := range criterion.dict {
			if _, ok := c.dict[i]; !ok {
				t.Errorf("expected key %s, found nil", i)
				return
			}
			if c.dict[i].value != criterion.dict[i] {
				t.Errorf("Expected %v, got %v", criterion.dict[i], c.dict[i].value)
			}
		}

		if _, ok := criterion.dict[criterion.removedKey]; ok {
			t.Errorf("expected key a to be removed and it was not")
		}

		if len(c.dict) != 1 {
			t.Errorf("expected to have 1 dict, found %d", len(c.dict))
		}
	}
}

func TestWithSameKeyDifferentValue(t *testing.T) {
	c := New(3)
	c.Set("a", 3)
	c.Set("a", 2)
	c.Set("a", 1)

	if value := c.Get("a"); value != 1 {
		t.Errorf("expected to get 1, got %v", value)
	}

	if len(c.dict) != 1 {
		t.Errorf("expected to have 1 dict, found %d", len(c.dict))
	}

	if lenItems(c) != 1 {
		t.Errorf("expected to have 1 item, found %d", lenItems(c))
	}
}

func TestWithSameKeyDifferentValue2(t *testing.T) {
	c := New(3)
	c.Set("a", 3)
	c.Set("a", 4)
	c.Set("b", 1)
	c.Set("a", 6)
	c.Set("c", 3)
	c.Set("a", 5)
	c.Set("b", 7)
	c.Set("b", 9)

	if value := c.Get("a"); value != 5 {
		t.Errorf("expected to get 5, got %v", value)
	}

	if value := c.Get("b"); value != 9 {
		t.Errorf("expected to get 9, got %v", value)
	}

	if value := c.Get("c"); value != 3 {
		t.Errorf("expected to get 3, got %v", value)
	}

	if len(c.dict) != 3 {
		t.Errorf("expected to have 3 items on dict, found %d", len(c.dict))
	}

	if lenItems(c) != 3 {
		t.Errorf("expected to have 3 item, found %d", lenItems(c))
	}
}

func lenItems(c *Cache) int {
	count := 0
	for i := c.items; i != nil; i = i.next {
		// fmt.Println(i.key)
		count++
	}
	return count
}

func printItems(c *Cache) {
	fmt.Println("")
	for i := c.items; i != nil; i = i.next {
		p := ""
		if i.prev != nil {
			p = "<-"
		}
		n := ""
		if i.next != nil {
			n = "->"
		}
		fmt.Printf("%s=%v%s%s\n", i.key, i.value, p, n)
	}
}
