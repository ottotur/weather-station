package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type config struct {
	listen, cacheDir, staticDir, sourceURL string
	now                                    func() time.Time
}

func loadConfig(getenv func(string) string) (config, error) {
	value := func(key, fallback string) string {
		if v := getenv(key); v != "" {
			return v
		}
		return fallback
	}
	c := config{
		listen:    value("WEATHER_LISTEN", "127.0.0.1:8080"),
		cacheDir:  value("WEATHER_CACHE_DIR", "cache"),
		staticDir: value("WEATHER_STATIC_DIR", "web/dist"),
		sourceURL: value("WEATHER_SOURCE_URL", "https://opendata.fmi.fi/wfs"),
		now:       time.Now,
	}
	if fixed := getenv("WEATHER_TEST_NOW"); fixed != "" {
		t, err := time.Parse(time.RFC3339, fixed)
		if err != nil {
			return c, fmt.Errorf("WEATHER_TEST_NOW: expected RFC3339 time: %w", err)
		}
		c.now = func() time.Time { return t }
	}
	host, port, err := net.SplitHostPort(c.listen)
	if err != nil || host == "" {
		return c, fmt.Errorf("WEATHER_LISTEN: expected host:port")
	}
	if n, err := strconv.Atoi(port); err != nil || n < 1 || n > 65535 {
		return c, fmt.Errorf("WEATHER_LISTEN: invalid port")
	}
	u, err := url.Parse(c.sourceURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.User != nil || u.Fragment != "" {
		return c, fmt.Errorf("WEATHER_SOURCE_URL: expected HTTP(S) URL without credentials or fragment")
	}
	if strings.TrimSpace(c.cacheDir) == "" || strings.TrimSpace(c.staticDir) == "" {
		return c, fmt.Errorf("cache and static directories must not be empty")
	}
	info, err := os.Stat(c.staticDir)
	if err != nil || !info.IsDir() {
		return c, fmt.Errorf("WEATHER_STATIC_DIR: directory is not accessible: %s", c.staticDir)
	}
	if err := os.MkdirAll(c.cacheDir, 0700); err != nil {
		return c, fmt.Errorf("WEATHER_CACHE_DIR: %w", err)
	}
	info, err = os.Stat(c.cacheDir)
	if err != nil || !info.IsDir() {
		return c, fmt.Errorf("WEATHER_CACHE_DIR: not a directory")
	}
	return c, nil
}
