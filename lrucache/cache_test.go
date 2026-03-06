package lrucache

import "testing"

func TestCacheEvictsLeastRecentlyUsed(t *testing.T) {
	c := New[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)

	if _, ok := c.Get("a"); !ok {
		t.Fatalf("expected key a")
	}

	c.Put("c", 3)

	if _, ok := c.Get("b"); ok {
		t.Fatalf("expected key b to be evicted")
	}
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("expected key a to exist with value 1")
	}
	if v, ok := c.Get("c"); !ok || v != 3 {
		t.Fatalf("expected key c to exist with value 3")
	}
}

func TestCacheDelete(t *testing.T) {
	c := New[string, int](1)
	c.Put("x", 42)

	if ok := c.Delete("x"); !ok {
		t.Fatalf("expected delete to return true")
	}
	if ok := c.Delete("x"); ok {
		t.Fatalf("expected second delete to return false")
	}
}
