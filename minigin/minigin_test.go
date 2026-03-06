package minigin

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareOrderAndNext(t *testing.T) {
	app := New()
	steps := make([]string, 0, 4)

	app.Use(func(c *Context) {
		steps = append(steps, "mw-before")
		c.Next()
		steps = append(steps, "mw-after")
	})

	app.GET("/ping", func(c *Context) {
		steps = append(steps, "handler")
		c.String(http.StatusOK, "pong")
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	expected := []string{"mw-before", "handler", "mw-after"}
	if len(steps) != len(expected) {
		t.Fatalf("expected steps %v, got %v", expected, steps)
	}
	for i := range expected {
		if steps[i] != expected[i] {
			t.Fatalf("expected steps %v, got %v", expected, steps)
		}
	}
}

func TestPathParams(t *testing.T) {
	app := New()
	app.GET("/users/:id", func(c *Context) {
		c.String(http.StatusOK, c.Param("id"))
	})

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != "42" {
		t.Fatalf("expected body 42, got %q", rr.Body.String())
	}
}

func TestMethodNotAllowed(t *testing.T) {
	app := New()
	app.GET("/resource", func(c *Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodPost, "/resource", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}
