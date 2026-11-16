package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"weather-station/internal/waterlevel"
)

func TestIndependentRestartAndCorruption(t *testing.T) {
	dir := t.TempDir()
	s := Store{Dir: dir}
	now := time.Date(2026, 11, 16, 18, 0, 0, 0, time.UTC)
	v := 40.0
	o := &waterlevel.Observation{ValueCM: 32.1, ObservedAt: now, FetchedAt: now}
	f := &waterlevel.Forecast{AnalysisTime: now.Add(-time.Hour), ResultTime: now, FetchedAt: now, UnitStatus: "assumed-cm", Points: []waterlevel.Point{{Time: now.Add(time.Hour), ValueCM: &v}, {Time: now.Add(2 * time.Hour)}}}
	if err := s.SaveObservation(o); err != nil {
		t.Fatal(err)
	}
	if err := s.SaveForecast(f); err != nil {
		t.Fatal(err)
	}
	got, err := (Store{Dir: dir}).LoadForecast()
	if err != nil || got.UnitStatus != "assumed-cm" || !got.AnalysisTime.Equal(f.AnalysisTime) || got.Points[1].ValueCM != nil {
		t.Fatalf("%+v %v", got, err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".forecast.json-interrupted.tmp"), []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := s.LoadForecast(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "observation.json"), []byte("{"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := s.LoadObservation(); err == nil {
		t.Fatal("accepted corrupt observation")
	}
	if _, err := s.LoadForecast(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "observation.json"), []byte(`{"version":42,"data":{}}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := s.LoadObservation(); err == nil {
		t.Fatal("accepted unknown version")
	}
}

func TestFailedWritePreservesPrevious(t *testing.T) {
	dir := t.TempDir()
	s := Store{Dir: dir}
	now := time.Now().UTC()
	o := &waterlevel.Observation{ValueCM: 32.1, ObservedAt: now, FetchedAt: now}
	if err := s.SaveObservation(o); err != nil {
		t.Fatal(err)
	}
	invalid := *o
	invalid.ObservedAt = now.Add(time.Hour)
	if err := s.SaveObservation(&invalid); err == nil {
		t.Fatal("invalid value accepted")
	}
	old, err := s.LoadObservation()
	if err != nil || old.ValueCM != 32.1 {
		t.Fatalf("%+v %v", old, err)
	}
	bad := Store{Dir: filepath.Join(dir, "not-a-directory")}
	if err := bad.SaveObservation(o); err == nil {
		t.Fatal("write unexpectedly succeeded")
	}
}
