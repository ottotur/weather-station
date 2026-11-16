package waterlevel

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type fakeSource struct {
	calls   int
	mu      sync.Mutex
	obs     *Observation
	fail    bool
	started chan struct{}
	release chan struct{}
}

func (f *fakeSource) Observation(ctx context.Context, now time.Time) (*Observation, error) {
	f.mu.Lock()
	f.calls++
	f.mu.Unlock()
	if f.started != nil {
		close(f.started)
		<-f.release
	}
	if f.fail {
		return nil, errors.New("source failed")
	}
	return f.obs, nil
}
func (f *fakeSource) Forecast(context.Context, time.Time) (*Forecast, error) {
	return nil, errors.New("unavailable")
}

type fakeStore struct {
	obs  *Observation
	fail bool
}

func (s *fakeStore) LoadObservation() (*Observation, error) { return s.obs, nil }
func (s *fakeStore) LoadForecast() (*Forecast, error)       { return nil, errors.New("missing") }
func (s *fakeStore) SaveObservation(o *Observation) error {
	if s.fail {
		return errors.New("disk")
	}
	s.obs = o
	return nil
}
func (s *fakeStore) SaveForecast(*Forecast) error { return nil }

func TestServiceRetainsPublishedData(t *testing.T) {
	now := time.Now().UTC()
	old := &Observation{ValueCM: 1, ObservedAt: now, FetchedAt: now}
	next := &Observation{ValueCM: 2, ObservedAt: now, FetchedAt: now}
	source := &fakeSource{obs: next}
	store := &fakeStore{obs: old, fail: true}
	s := NewService(source, store, func() time.Time { return now })
	s.RefreshObservation(context.Background())
	got := s.Snapshot()
	if got.Observation.ValueCM != 1 || got.ObservationAttempt.Result != "error" || *got.ObservationAttempt.ErrorCode != "cache_error" {
		t.Fatalf("%+v", got)
	}
	store.fail = false
	s.RefreshObservation(context.Background())
	got = s.Snapshot()
	if got.Observation.ValueCM != 2 || got.ObservationAttempt.Result != "success" {
		t.Fatalf("%+v", got)
	}
	got.Observation.ValueCM = 99
	if s.Snapshot().Observation.ValueCM != 2 {
		t.Fatal("snapshot mutable")
	}
	source.fail = true
	s.RefreshObservation(context.Background())
	got = s.Snapshot()
	if got.Observation.ValueCM != 2 || *got.ObservationAttempt.ErrorCode != "source_error" {
		t.Fatalf("%+v", got)
	}
}

func TestRefreshDoesNotOverlap(t *testing.T) {
	now := time.Now().UTC()
	source := &fakeSource{obs: &Observation{ObservedAt: now, FetchedAt: now}, started: make(chan struct{}), release: make(chan struct{})}
	s := NewService(source, &fakeStore{}, func() time.Time { return now })
	done := make(chan struct{})
	go func() { s.RefreshObservation(context.Background()); close(done) }()
	<-source.started
	s.RefreshObservation(context.Background())
	close(source.release)
	<-done
	if source.calls != 1 {
		t.Fatalf("calls: %d", source.calls)
	}
}

func TestScheduleWaitsForNormalTickAfterFailure(t *testing.T) {
	for _, interval := range []time.Duration{10 * time.Minute, 30 * time.Minute} {
		t.Run(interval.String(), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			ticks := make(chan time.Time)
			calls := make(chan struct{}, 2)
			ready := make(chan time.Duration, 1)
			done := make(chan struct{})
			go func() {
				defer close(done)
				scheduleWithTicker(ctx, interval, func(context.Context) { calls <- struct{}{} }, func(d time.Duration) (<-chan time.Time, func()) {
					ready <- d
					return ticks, func() {}
				})
			}()
			<-calls
			if got := <-ready; got != interval {
				t.Fatalf("interval %v", got)
			}
			select {
			case <-calls:
				t.Fatal("immediate retry")
			default:
			}
			ticks <- time.Date(2026, 11, 16, 18, 0, 0, 0, time.UTC).Add(interval)
			<-calls
			cancel()
			<-done
		})
	}
}
