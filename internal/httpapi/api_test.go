package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"weather-station/internal/cache"
	"weather-station/internal/fmi"
	"weather-station/internal/waterlevel"
)

func TestObservationResultTimeSurvivesRestart(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 8, 0, 0, time.UTC)
	fixture, err := os.ReadFile("../../docs/feasibility/evidence/hamina-observation.xml")
	if err != nil {
		t.Fatal(err)
	}
	xml := strings.Replace(string(fixture), "2026-11-16T18:00:00Z</gml:timePosition>", "2026-11-16T17:59:30Z</gml:timePosition>", 1)
	if xml == string(fixture) {
		t.Fatal("fixture result time not found")
	}
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(xml)) }))
	defer source.Close()
	store := cache.Store{Dir: t.TempDir()}
	clock := func() time.Time { return now }
	service := waterlevel.NewService(fmi.NewClient(source.URL), store, clock)
	service.RefreshObservation(context.Background())
	check := func(s *waterlevel.Service) {
		t.Helper()
		res := httptest.NewRecorder()
		(API{Service: s, Now: clock}).ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
		var doc map[string]any
		if err := json.Unmarshal(res.Body.Bytes(), &doc); err != nil {
			t.Fatal(err)
		}
		obs, ok := doc["observation"].(map[string]any)["data"].(map[string]any)
		if !ok {
			t.Fatalf("missing observation: %s", res.Body.String())
		}
		if obs["resultTime"] != "2026-11-16T17:59:30Z" || obs["observedAt"] != "2026-11-16T18:00:00Z" || obs["fetchedAt"] != "2026-11-16T18:08:00Z" {
			t.Fatalf("source times not preserved separately: %v", obs)
		}
	}
	check(service)
	check(waterlevel.NewService(fmi.NewClient(source.URL), store, clock))
}

func TestUnitMetadataErrorPreservesPublishedDatasets(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 8, 0, 0, time.UTC)
	obs, err := os.ReadFile("../../docs/feasibility/evidence/hamina-observation.xml")
	if err != nil {
		t.Fatal(err)
	}
	forecast, err := os.ReadFile("../../docs/feasibility/evidence/hamina-forecast.xml")
	if err != nil {
		t.Fatal(err)
	}
	var invalid atomic.Bool
	var externalCalls atomic.Int32
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/unit" {
			externalCalls.Add(1)
			http.Error(w, "unexpected request", 500)
			return
		}
		text := string(forecast)
		marker := `<wml2:MeasurementTimeseries gml:id="mts-1-1-SeaLevelN2000">`
		if strings.Contains(r.URL.Query().Get("storedquery_id"), "observations") {
			text = string(obs)
			marker = `<wml2:MeasurementTimeseries gml:id="obs-obs-1-1-WLEVN2K_PT1S_INSTANT">`
		}
		if invalid.Load() {
			text = strings.Replace(text, marker, marker+`<wml2:uom xlink:href="`+server.URL+`/unit"/>`, 1)
		}
		_, _ = w.Write([]byte(text))
	}))
	defer server.Close()
	store := cache.Store{Dir: t.TempDir()}
	clock := func() time.Time { return now }
	service := waterlevel.NewService(fmi.NewClient(server.URL), store, clock)
	service.RefreshObservation(context.Background())
	service.RefreshForecast(context.Background())
	before := service.Snapshot()
	if before.Observation == nil || before.Forecast == nil {
		t.Fatal("initial fixtures were not published")
	}
	invalid.Store(true)
	service.RefreshObservation(context.Background())
	service.RefreshForecast(context.Background())
	after := service.Snapshot()
	for _, attempt := range []waterlevel.Attempt{after.ObservationAttempt, after.ForecastAttempt} {
		if attempt.Result != "error" || attempt.ErrorCode == nil || *attempt.ErrorCode != "unit_metadata_error" {
			t.Fatalf("attempt: %+v", attempt)
		}
	}
	if after.Observation.ValueCM != before.Observation.ValueCM || !after.Forecast.AnalysisTime.Equal(before.Forecast.AnalysisTime) {
		t.Fatal("published data changed on unit failure")
	}
	restarted := waterlevel.NewService(fmi.NewClient(server.URL), store, clock).Snapshot()
	if restarted.Observation == nil || restarted.Forecast == nil || restarted.Observation.ValueCM != before.Observation.ValueCM || !restarted.Forecast.AnalysisTime.Equal(before.Forecast.AnalysisTime) {
		t.Fatal("last valid data did not survive restart")
	}
	if externalCalls.Load() != 0 {
		t.Fatal("external unit reference was requested")
	}
	res := httptest.NewRecorder()
	(API{Service: service, Now: clock}).ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	if !strings.Contains(res.Body.String(), `"errorCode":"unit_metadata_error"`) {
		t.Fatalf("missing API error code: %s", res.Body.String())
	}
}

type empty struct{}

type seeded struct {
	empty
	observation *waterlevel.Observation
	forecast    *waterlevel.Forecast
}

type countingSource struct{ calls atomic.Int32 }

func (s *countingSource) Observation(context.Context, time.Time) (*waterlevel.Observation, error) {
	s.calls.Add(1)
	return nil, errors.New("unexpected fetch")
}
func (s *countingSource) Forecast(context.Context, time.Time) (*waterlevel.Forecast, error) {
	s.calls.Add(1)
	return nil, errors.New("unexpected fetch")
}

func (s seeded) LoadObservation() (*waterlevel.Observation, error) { return s.observation, nil }
func (s seeded) LoadForecast() (*waterlevel.Forecast, error)       { return s.forecast, nil }

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

func TestForecastCoverageAndRange(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 0, 0, 0, time.UTC)
	v := 40.0
	forecast := &waterlevel.Forecast{AnalysisTime: now, ResultTime: now, FetchedAt: now, UnitStatus: "assumed-cm", Points: []waterlevel.Point{
		{Time: now.Add(-time.Hour), ValueCM: &v}, {Time: now.Add(time.Hour), ValueCM: &v}, {Time: now.Add(2 * time.Hour)}, {Time: now.Add(49 * time.Hour), ValueCM: &v},
	}}
	seed := seeded{observation: &waterlevel.Observation{ObservedAt: now, FetchedAt: now, ValueCM: 32}, forecast: forecast}
	api := API{Service: waterlevel.NewService(seed, seed, func() time.Time { return now }), Now: func() time.Time { return now }}
	res := httptest.NewRecorder()
	api.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	var data map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &data); err != nil {
		t.Fatal(err)
	}
	f := data["forecast"].(map[string]any)
	meta := f["data"].(map[string]any)
	if _, ok := meta["points"]; ok {
		t.Fatal("unfiltered source points exposed")
	}
	if f["lastShownAt"] != now.Add(time.Hour).Format(time.RFC3339) || meta["lastValidAt"] != now.Add(49*time.Hour).Format(time.RFC3339) {
		t.Fatalf("%v", f)
	}
	points := f["points"].([]any)
	if len(points) != 2 || points[1].(map[string]any)["valueCm"] != nil {
		t.Fatalf("%v", points)
	}
}

func TestConcurrentClientsDoNotFetch(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 0, 0, 0, time.UTC)
	source := &countingSource{}
	api := API{Service: waterlevel.NewService(source, empty{}, func() time.Time { return now }), Now: func() time.Time { return now }}
	results := make(chan string, 10)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res := httptest.NewRecorder()
			api.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
			results <- res.Body.String()
		}()
	}
	wg.Wait()
	close(results)
	first := ""
	for result := range results {
		if first == "" {
			first = result
		} else if result != first {
			t.Fatal("inconsistent snapshots")
		}
	}
	if source.calls.Load() != 0 {
		t.Fatal("API triggered source fetch")
	}
}

func TestPartialResponseSeparatesFreshnessAndAttempt(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 0, 0, 0, time.UTC)
	obs := &waterlevel.Observation{ValueCM: 32.1, ObservedAt: now.Add(-time.Minute), FetchedAt: now}
	seed := seeded{observation: obs}
	source := &countingSource{}
	service := waterlevel.NewService(source, seed, func() time.Time { return now })
	service.RefreshObservation(context.Background())
	api := API{Service: service, Now: func() time.Time { return now }}
	res := httptest.NewRecorder()
	api.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	var data map[string]any
	if err := json.Unmarshal(res.Body.Bytes(), &data); err != nil {
		t.Fatal(err)
	}
	observation := data["observation"].(map[string]any)
	if observation["status"] != "fresh" || observation["attempt"].(map[string]any)["errorCode"] != "source_error" {
		t.Fatalf("%v", observation)
	}
	obsData := observation["data"].(map[string]any)
	if value, present := obsData["resultTime"]; !present || value != nil {
		t.Fatalf("unknown legacy result time must be null: %v", obsData)
	}
	if data["forecast"].(map[string]any)["data"] != nil || data["summary"] != nil || data["summaryReason"] != "forecast_missing" {
		t.Fatalf("%v", data)
	}
	sourceInfo := data["source"].(map[string]any)
	if sourceInfo["licence"] != "CC BY 4.0" || sourceInfo["forecastUnit"] != "assumed-cm" {
		t.Fatalf("%v", sourceInfo)
	}
}
