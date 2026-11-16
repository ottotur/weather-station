package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {
	dir := t.TempDir()
	base := map[string]string{"WEATHER_STATIC_DIR": dir, "WEATHER_CACHE_DIR": filepath.Join(dir, "cache")}
	get := func(k string) string { return base[k] }
	c, err := loadConfig(get)
	if err != nil || c.listen != "127.0.0.1:8080" || c.now == nil {
		t.Fatalf("%+v %v", c, err)
	}
	for _, tc := range []struct{ key, value string }{
		{"WEATHER_LISTEN", ":8080"}, {"WEATHER_LISTEN", "localhost:99999"},
		{"WEATHER_SOURCE_URL", "file:///tmp/x"}, {"WEATHER_SOURCE_URL", "https://user:pass@example.com"},
		{"WEATHER_STATIC_DIR", filepath.Join(dir, "missing")}, {"WEATHER_CACHE_DIR", " "},
	} {
		old := base[tc.key]
		base[tc.key] = tc.value
		_, err := loadConfig(get)
		if err == nil || !strings.Contains(err.Error(), "WEATHER_") && tc.key != "WEATHER_CACHE_DIR" {
			t.Errorf("%s=%q: %v", tc.key, tc.value, err)
		}
		base[tc.key] = old
	}
}
