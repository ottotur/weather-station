package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"weather-station/internal/cache"
	"weather-station/internal/fmi"
	"weather-station/internal/waterlevel"
)

func TestBuiltUIFromGoServer(t *testing.T) {
	dist := os.Getenv("WEATHER_TEST_WEB_DIST")
	if dist == "" {
		t.Skip("set WEATHER_TEST_WEB_DIST to the built web/dist directory")
	}
	dist, err := filepath.Abs(dist)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now
	service := waterlevel.NewService(fmi.NewClient("http://127.0.0.1:1"), cache.Store{Dir: t.TempDir()}, now)
	server := httptest.NewServer(handler(dist, service, now))
	defer server.Close()
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil || resp.StatusCode != http.StatusOK || !strings.Contains(string(body), "Haminan vedenkorkeus") {
		t.Fatalf("status %d body %s error %v", resp.StatusCode, body, err)
	}
	asset := regexp.MustCompile(`src="(/assets/[^"]+\.js)"`).FindSubmatch(body)
	if len(asset) != 2 {
		t.Fatal("built script not found in HTML")
	}
	resp, err = http.Get(server.URL + string(asset[1]))
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("script status %d", resp.StatusCode)
	}
	resp, err = http.Get(server.URL + "/api/water-level/hamina")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || resp.Header.Get("Cache-Control") != "no-store" {
		t.Fatalf("API status %d headers %v", resp.StatusCode, resp.Header)
	}
}
