package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func pageJSON(next string, previous string, names ...string) string {
	res := make([]byte, 0, 128)
	for i, n := range names {
		if i > 0 {
			res = append(res, ',')
		}
		res = append(res, []byte(fmt.Sprintf(`{"name":%q,"url":""}`, n))...)
	}
	return fmt.Sprintf(`{"next":%q,"previous":%q,"results":[%s]}`, next, previous, res)
}

func TestCommandMapNavigatesPages(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/p1":
			fmt.Fprint(w, pageJSON("http://"+r.Host+"/p2", "", "twinleaf-town", "sandgem-town"))
		case "/p2":
			fmt.Fprint(w, pageJSON("", "http://"+r.Host+"/p1", "oreburgh-gate"))
		default:
			http.Error(w, "unexpected path", http.StatusNotFound)
		}
	}))
	defer srv.Close()

	cfg := newTestConfig()
	cfg.next = srv.URL + "/p1"

	out := captureStdout(t, func() { _ = command_map(cfg) })
	if !strings.Contains(out, "twinleaf-town") || !strings.Contains(out, "sandgem-town") {
		t.Errorf("expected first page town names in output, got %q", out)
	}
	if cfg.next != srv.URL+"/p2" {
		t.Errorf("expected next = %q, got %q", srv.URL+"/p2", cfg.next)
	}
	if cfg.prev != "" {
		t.Errorf("expected prev = %q, got %q", "", cfg.prev)
	}

	out = captureStdout(t, func() { _ = command_map(cfg) })
	if !strings.Contains(out, "oreburgh-gate") {
		t.Errorf("expected second page town name in output, got %q", out)
	}
	if cfg.next != "" {
		t.Errorf("expected next = %q after last page, got %q", "", cfg.next)
	}

	out = captureStdout(t, func() { _ = command_map(cfg) })
	if out != "last page\n" {
		t.Errorf("expected exhausted-page message, got %q", out)
	}
}

func TestCommandMapBackPage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/p1" {
			http.Error(w, "unexpected path", http.StatusNotFound)
			return
		}
		fmt.Fprint(w, pageJSON("", "http://"+r.Host+"/p0", "pallet-town"))
	}))
	defer srv.Close()

	cfg := newTestConfig()
	cfg.prev = srv.URL + "/p1"

	out := captureStdout(t, func() { _ = command_bmap(cfg) })
	if !strings.Contains(out, "pallet-town") {
		t.Errorf("expected town name in output, got %q", out)
	}
	if cfg.next != "" {
		t.Errorf("expected next = %q, got %q", "", cfg.next)
	}
	if cfg.prev != srv.URL+"/p0" {
		t.Errorf("expected prev = %q, got %q", srv.URL+"/p0", cfg.prev)
	}
}

func TestCommandMapBoundariesShortCircuit(t *testing.T) {
	cfg := newTestConfig()

	out := captureStdout(t, func() { _ = command_map(cfg) })
	if out != "last page\n" {
		t.Errorf("expected %q when next is empty, got %q", "last page\n", out)
	}

	out = captureStdout(t, func() { _ = command_bmap(cfg) })
	if out != "you're on the first page\n" {
		t.Errorf("expected %q when prev is empty, got %q", "you're on the first page\n", out)
	}
}

func TestGetPageStatusError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "area not found")
	}))
	defer srv.Close()

	cfg := newTestConfig()
	if _, err := getPage(cfg, srv.URL); err == nil {
		t.Fatal("expected error for status >299, got nil")
	} else if !strings.Contains(err.Error(), "404") || !strings.Contains(err.Error(), "area not found") {
		t.Errorf("expected error to include status and body, got %q", err.Error())
	}

	if _, ok := cfg.cache.Get(srv.URL); ok {
		t.Errorf("expected failed response to not be cached")
	}
}

func TestGetPageInvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "not-json")
	}))
	defer srv.Close()

	cfg := newTestConfig()
	if _, err := getPage(cfg, srv.URL); err == nil {
		t.Fatal("expected unmarshal error for non-JSON body, got nil")
	}
}

func TestGetPageUsesCache(t *testing.T) {
	var hits int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&hits, 1)
		fmt.Fprint(w, pageJSON("", "", "server-town"))
	}))
	defer srv.Close()

	cfg := newTestConfig()
	if _, ok := cfg.cache.Add(srv.URL, []byte(pageJSON("", "", "cached-town"))); ok {
		t.Fatalf("expected fresh cache key")
	}

	page, err := getPage(cfg, srv.URL)
	if err != nil {
		t.Fatalf("getPage: %v", err)
	}
	if hits != 0 {
		t.Errorf("expected zero server hits on cache hit, got %d", hits)
	}
	if len(page.Results) != 1 || page.Results[0].Name != "cached-town" {
		t.Errorf("expected cached page content, got %+v", page.Results)
	}
}
