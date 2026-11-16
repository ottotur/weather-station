package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"weather-station/internal/waterlevel"
)

type empty struct{}

func (empty) Observation(context.Context, time.Time) (*waterlevel.Observation, error) {
	return nil, errors.New("missing")
}
func (empty) Forecast(context.Context, time.Time) (*waterlevel.Forecast, error) {
	return nil, errors.New("missing")
}
func (empty) LoadObservation() (*waterlevel.Observation, error) { return nil, errors.New("missing") }
func (empty) LoadForecast() (*waterlevel.Forecast, error)       { return nil, errors.New("missing") }
func (empty) SaveObservation(*waterlevel.Observation) error     { return nil }
func (empty) SaveForecast(*waterlevel.Forecast) error           { return nil }

func TestEmptyResponse(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 0, 0, 0, time.UTC)
	clock := func() time.Time { return now }
	s := waterlevel.NewService(empty{}, empty{}, clock)
	api := API{Service: s, Now: clock}
	res := httptest.NewRecorder()
	api.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/api/water-level/hamina", nil))
	if res.Code != 200 || res.Header().Get("Cache-Control") != "no-store" {
		t.Fatal(res.Code, res.Header())
	}
	var data map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &data); err != nil {
		t.Fatal(err)
	}
	obs := data["observation"].(map[string]any)
	if obs["status"] != "missing" || obs["data"] != nil || data["summary"] != nil {
		t.Fatal(data)
	}
	s.RefreshObservation(context.Background())
	res = httptest.NewRecorder()
	api.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	if err := json.Unmarshal(res.Body.Bytes(), &data); err != nil {
		t.Fatal(err)
	}
	attempt := data["observation"].(map[string]any)["attempt"].(map[string]any)
	if attempt["result"] != "error" || attempt["errorCode"] != "source_error" {
		t.Fatal(attempt)
	}
}
