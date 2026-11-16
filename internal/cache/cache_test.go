package cache

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"weather-station/internal/waterlevel"
)

type observationSource struct{ value *waterlevel.Observation }

func (s observationSource) Observation(context.Context, time.Time) (*waterlevel.Observation, error) {
	return s.value, nil
}
func (s observationSource) Forecast(context.Context, time.Time) (*waterlevel.Forecast, error) {
	return nil, errors.New("not used")
}

func TestSyncFailureAfterReplacement(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 0, 0, 0, time.UTC)
	dir := t.TempDir()
	normal := Store{Dir: dir}
	old := &waterlevel.Observation{ValueCM: 10, ObservedAt: now, FetchedAt: now}
	next := &waterlevel.Observation{ValueCM: 20, ObservedAt: now, FetchedAt: now}
	v := 40.0
	forecast := &waterlevel.Forecast{AnalysisTime: now, ResultTime: now, FetchedAt: now, UnitStatus: "assumed-cm", Points: []waterlevel.Point{{Time: now, ValueCM: &v}}}
	if err := normal.SaveObservation(old); err != nil {
		t.Fatal(err)
	}
	if err := normal.SaveForecast(forecast); err != nil {
		t.Fatal(err)
	}
	broken := Store{Dir: dir, syncDir: func(string) error { return errors.New("simulated sync failure") }}
	service := waterlevel.NewService(observationSource{next}, broken, func() time.Time { return now })
	service.RefreshObservation(context.Background())
	snapshot := service.Snapshot()
	if snapshot.Observation.ValueCM != 10 || snapshot.ObservationAttempt.ErrorCode == nil || *snapshot.ObservationAttempt.ErrorCode != "cache_durability_uncertain" {
		t.Fatalf("unexpected in-memory state: %+v", snapshot)
	}
	if snapshot.Forecast == nil || *snapshot.Forecast.Points[0].ValueCM != v {
		t.Fatal("other dataset changed")
	}
	restarted := waterlevel.NewService(observationSource{next}, normal, func() time.Time { return now })
	if got := restarted.Snapshot().Observation; got == nil || got.ValueCM != 20 {
		t.Fatalf("replacement not visible after restart: %+v", got)
	}
}

func TestInterruptedBeforeRenamePreservesPublishedAndDisk(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 0, 0, 0, time.UTC)
	dir := t.TempDir()
	normal := Store{Dir: dir}
	old := &waterlevel.Observation{ValueCM: 10, ObservedAt: now, FetchedAt: now}
	next := &waterlevel.Observation{ValueCM: 20, ObservedAt: now, FetchedAt: now}
	if err := normal.SaveObservation(old); err != nil {
		t.Fatal(err)
	}
	broken := Store{Dir: dir, beforeRename: func(path string) error {
		if err := os.WriteFile(path, []byte(`{"version":1`), 0600); err != nil {
			return err
		}
		return errors.New("interrupted write")
	}}
	service := waterlevel.NewService(observationSource{next}, broken, func() time.Time { return now })
	service.RefreshObservation(context.Background())
	snapshot := service.Snapshot()
	if snapshot.Observation.ValueCM != 10 || snapshot.ObservationAttempt.ErrorCode == nil || *snapshot.ObservationAttempt.ErrorCode != "cache_error" {
		t.Fatalf("%+v", snapshot)
	}
	restarted, err := normal.LoadObservation()
	if err != nil || restarted.ValueCM != 10 {
		t.Fatalf("%+v %v", restarted, err)
	}
}

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
