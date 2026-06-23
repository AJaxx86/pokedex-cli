package pokeapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	cl := NewClient()
	if cl.cache == nil {
		t.Errorf("expected cache to be initialized")
	}
}

func TestGetAreasEmptyURL(t *testing.T) {
	cl := NewClient()
	areas, next, prev, err := cl.GetAreas("")
	if err == nil {
		t.Errorf("expected error for empty URL, got nil")
	}
	if areas != nil {
		t.Errorf("expected nil areas, got %v", areas)
	}
	if next != "" || prev != "" {
		t.Errorf("expected empty next/prev, got %q / %q", next, prev)
	}
}

func TestGetAreasSuccess(t *testing.T) {
	body := `{
		"count": 2,
		"next": "https://example.com/next",
		"previous": null,
		"results": [
			{"name": "area1", "url": "https://example.com/area1"},
			{"name": "area2", "url": "https://example.com/area2"}
		]
	}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, body)
	}))
	defer server.Close()

	cl := NewClient()
	areas, next, prev, err := cl.GetAreas(server.URL)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(areas) != 2 {
		t.Fatalf("expected 2 areas, got %d (%v)", len(areas), areas)
	}
	if areas[0] != "area1" || areas[1] != "area2" {
		t.Errorf("unexpected area names: %v", areas)
	}
	if next != "https://example.com/next" {
		t.Errorf("unexpected next: %q", next)
	}
	if prev != "" {
		t.Errorf("expected empty previous, got %q", prev)
	}
}

func TestGetAreasCacheHit(t *testing.T) {
	body := `{
		"count": 1,
		"next": null,
		"previous": null,
		"results": [
			{"name": "cached-area", "url": "https://example.com/cached"}
		]
	}`
	hits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, body)
	}))
	defer server.Close()

	cl := NewClient()

	// First call: cache miss, hits the server.
	areas, _, _, err := cl.GetAreas(server.URL)
	if err != nil {
		t.Errorf("unexpected error on first call: %v", err)
	}
	if len(areas) != 1 || areas[0] != "cached-area" {
		t.Fatalf("unexpected areas on first call: %v", areas)
	}

	// Second call: should be served from cache, no extra HTTP hit.
	areas, _, _, err = cl.GetAreas(server.URL)
	if err != nil {
		t.Errorf("unexpected error on second call: %v", err)
	}
	if len(areas) != 1 || areas[0] != "cached-area" {
		t.Fatalf("unexpected areas on second call: %v", areas)
	}
	if hits != 1 {
		t.Errorf("expected server to be hit once, got %d hits", hits)
	}
}

func TestGetAreasPreviousURL(t *testing.T) {
	body := `{
		"count": 2,
		"next": null,
		"previous": "https://example.com/prev",
		"results": [
			{"name": "area1", "url": "https://example.com/area1"}
		]
	}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, body)
	}))
	defer server.Close()

	cl := NewClient()
	areas, next, prev, err := cl.GetAreas(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(areas) != 1 {
		t.Fatalf("expected 1 area, got %d", len(areas))
	}
	if next != "" {
		t.Errorf("expected empty next, got %q", next)
	}
	if prev != "https://example.com/prev" {
		t.Errorf("unexpected previous: %q", prev)
	}
}

func TestGetAreasTransportError(t *testing.T) {
	cl := NewClient()
	areas, _, _, err := cl.GetAreas("http://[::1]:0")
	if err == nil {
		t.Errorf("expected transport error, got nil")
	}
	if areas != nil {
		t.Errorf("expected nil areas on error, got %v", areas)
	}
}

func TestGetAreasHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cl := NewClient()
	areas, _, _, err := cl.GetAreas(server.URL)
	if err == nil {
		t.Errorf("expected error for 404 response, got nil")
	}
	if areas != nil {
		t.Errorf("expected nil areas on error, got %v", areas)
	}
}

func TestGetAreasInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "not valid json {{{")
	}))
	defer server.Close()

	cl := NewClient()
	areas, _, _, err := cl.GetAreas(server.URL)
	if err == nil {
		t.Errorf("expected JSON unmarshal error, got nil")
	}
	if areas != nil {
		t.Errorf("expected nil areas on error, got %v", areas)
	}
}