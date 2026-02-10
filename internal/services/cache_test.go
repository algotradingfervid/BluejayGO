package services_test

import (
	"testing"
	"time"

	"github.com/narendhupati/bluejay-cms/internal/services"
)

func TestCache_GetSet(t *testing.T) {
	c := services.NewCache()

	c.Set("key1", "value1", 60)

	val, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected key1 to be found")
	}
	if val.(string) != "value1" {
		t.Errorf("expected 'value1', got %q", val)
	}
}

func TestCache_GetMissing(t *testing.T) {
	c := services.NewCache()

	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("expected nonexistent key to not be found")
	}
}

func TestCache_GetExpired(t *testing.T) {
	c := services.NewCache()

	c.Set("expire-me", "data", 1)
	time.Sleep(1100 * time.Millisecond)

	_, ok := c.Get("expire-me")
	if ok {
		t.Error("expected expired item to not be found")
	}
}

func TestCache_Delete(t *testing.T) {
	c := services.NewCache()

	c.Set("del-key", "data", 60)
	c.Delete("del-key")

	_, ok := c.Get("del-key")
	if ok {
		t.Error("expected deleted key to not be found")
	}
}

func TestCache_DeleteByPrefix(t *testing.T) {
	c := services.NewCache()

	c.Set("page:products", "list", 60)
	c.Set("page:products:detectors", "cat", 60)
	c.Set("page:products:detectors:alpha", "detail", 60)
	c.Set("page:blog", "blog", 60)

	c.DeleteByPrefix("page:products")

	_, ok := c.Get("page:products")
	if ok {
		t.Error("expected 'page:products' to be deleted")
	}
	_, ok = c.Get("page:products:detectors")
	if ok {
		t.Error("expected 'page:products:detectors' to be deleted")
	}
	_, ok = c.Get("page:products:detectors:alpha")
	if ok {
		t.Error("expected 'page:products:detectors:alpha' to be deleted")
	}

	// blog should remain
	val, ok := c.Get("page:blog")
	if !ok {
		t.Error("expected 'page:blog' to still exist")
	}
	if val.(string) != "blog" {
		t.Errorf("expected 'blog', got %q", val)
	}
}

func TestCache_SetOverwrite(t *testing.T) {
	c := services.NewCache()

	c.Set("key", "old", 60)
	c.Set("key", "new", 60)

	val, ok := c.Get("key")
	if !ok {
		t.Fatal("expected key to be found")
	}
	if val.(string) != "new" {
		t.Errorf("expected 'new', got %q", val)
	}
}

func TestCache_DeleteNonexistent(t *testing.T) {
	c := services.NewCache()
	// Should not panic
	c.Delete("nonexistent")
	c.DeleteByPrefix("nonexistent")
}
