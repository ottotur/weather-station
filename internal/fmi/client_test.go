package fmi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientQueries(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 8, 0, 0, time.UTC)
	obs, forecast := sample(t, "observation"), sample(t, "forecast")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		for key, value := range map[string]string{"service": "WFS", "version": "2.0.0", "request": "GetFeature", "fmisid": "134254", "timezone": "UTC"} {
			if q.Get(key) != value {
				t.Errorf("%s=%q", key, q.Get(key))
			}
		}
		switch q.Get("storedquery_id") {
		case "fmi::observations::mareograph::instant::timevaluepair":
			if q.Get("timestep") != "1" || q.Get("parameters") != "WLEVN2K_PT1S_INSTANT" || q.Get("starttime") != "2026-11-16T17:08:00Z" || q.Get("endtime") != "2026-11-16T18:08:00Z" {
				t.Errorf("observation query: %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(obs))
		case "fmi::forecast::sealevel::point::timevaluepair":
			if q.Get("timestep") != "60" || q.Has("parameters") || q.Get("starttime") != "2026-11-16T18:08:00Z" || q.Get("endtime") != "2026-11-18T18:08:00Z" {
				t.Errorf("forecast query: %s", r.URL.RawQuery)
			}
			_, _ = w.Write([]byte(forecast))
		default:
			t.Error("unexpected query")
		}
	}))
	defer server.Close()
	c := NewClient(server.URL)
	if _, err := c.Observation(context.Background(), now); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Forecast(context.Background(), now); err != nil {
		t.Fatal(err)
	}
}

func TestClientFailures(t *testing.T) {
	for _, tc := range []struct {
		name    string
		handler http.HandlerFunc
	}{
		{"HTTP", func(w http.ResponseWriter, r *http.Request) { http.Error(w, "busy", 429) }},
		{"WFS", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(`<ExceptionReport/>`)) }},
		{"size", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(strings.Repeat("x", maxResponseBytes+1)))
		}},
		{"slow", func(w http.ResponseWriter, r *http.Request) { <-r.Context().Done() }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)
			defer server.Close()
			c := NewClient(server.URL)
			c.HTTP.Timeout = 50 * time.Millisecond
			if _, err := c.Observation(context.Background(), time.Now()); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
