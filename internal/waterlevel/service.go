package waterlevel

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrDurabilityUncertain means the file was replaced but directory sync failed.
var ErrDurabilityUncertain = errors.New("cache durability uncertain")

// ErrUnitMetadata means a source unit is contradictory or cannot be interpreted safely.
var ErrUnitMetadata = errors.New("invalid source unit metadata")

type Source interface {
	Observation(context.Context, time.Time) (*Observation, error)
	Forecast(context.Context, time.Time) (*Forecast, error)
}
type Persistence interface {
	LoadObservation() (*Observation, error)
	LoadForecast() (*Forecast, error)
	SaveObservation(*Observation) error
	SaveForecast(*Forecast) error
}
type Attempt struct {
	At        *time.Time `json:"lastAttemptAt"`
	Result    string     `json:"result"`
	ErrorCode *string    `json:"errorCode"`
}

func newAttempt(at time.Time, sourceErr, cacheErr error) Attempt {
	attempt := Attempt{At: &at, Result: "success"}
	var code string
	switch {
	case errors.Is(sourceErr, ErrUnitMetadata):
		code = "unit_metadata_error"
	case sourceErr != nil:
		code = "source_error"
	case errors.Is(cacheErr, ErrDurabilityUncertain):
		code = "cache_durability_uncertain"
	case cacheErr != nil:
		code = "cache_error"
	default:
		return attempt
	}
	attempt.Result = "error"
	attempt.ErrorCode = &code
	return attempt
}

type Snapshot struct {
	Observation        *Observation
	Forecast           *Forecast
	ObservationAttempt Attempt
	ForecastAttempt    Attempt
}
type Service struct {
	mu                        sync.RWMutex
	observationMu, forecastMu sync.Mutex
	source                    Source
	store                     Persistence
	now                       func() time.Time
	state                     Snapshot
}

func NewService(source Source, store Persistence, now func() time.Time) *Service {
	s := &Service{source: source, store: store, now: now, state: Snapshot{ObservationAttempt: Attempt{Result: "pending"}, ForecastAttempt: Attempt{Result: "pending"}}}
	if o, err := store.LoadObservation(); err == nil {
		s.state.Observation = o
	}
	if f, err := store.LoadForecast(); err == nil {
		s.state.Forecast = f
	}
	return s
}

// Snapshot returns copies, so the caller cannot mutate published state.
func (s *Service) Snapshot() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := s.state
	if result.Observation != nil {
		o := *result.Observation
		if o.ResultTime != nil {
			resultTime := *o.ResultTime
			o.ResultTime = &resultTime
		}
		result.Observation = &o
	}
	if result.Forecast != nil {
		f := *result.Forecast
		f.Points = append([]Point(nil), f.Points...)
		for i := range f.Points {
			if f.Points[i].ValueCM != nil {
				v := *f.Points[i].ValueCM
				f.Points[i].ValueCM = &v
			}
		}
		result.Forecast = &f
	}
	return result
}

func (s *Service) RefreshObservation(ctx context.Context) {
	refresh(ctx, s, &s.observationMu, s.source.Observation, s.store.SaveObservation,
		func(state *Snapshot, value *Observation, attempt Attempt) {
			state.ObservationAttempt = attempt
			if attempt.Result == "success" {
				state.Observation = value
			}
		})
}

func (s *Service) RefreshForecast(ctx context.Context) {
	refresh(ctx, s, &s.forecastMu, s.source.Forecast, s.store.SaveForecast,
		func(state *Snapshot, value *Forecast, attempt Attempt) {
			state.ForecastAttempt = attempt
			if attempt.Result == "success" {
				state.Forecast = value
			}
		})
}

func refresh[T any](ctx context.Context, s *Service, lock *sync.Mutex,
	fetch func(context.Context, time.Time) (*T, error),
	save func(*T) error,
	publish func(*Snapshot, *T, Attempt),
) {
	if !lock.TryLock() {
		return
	}
	defer lock.Unlock()
	now := s.now().UTC()
	value, sourceErr := fetch(ctx, now)
	var cacheErr error
	if sourceErr == nil {
		cacheErr = save(value)
	}
	attempt := newAttempt(now, sourceErr, cacheErr)
	s.mu.Lock()
	defer s.mu.Unlock()
	publish(&s.state, value, attempt)
}

func (s *Service) Run(ctx context.Context) {
	go schedule(ctx, 10*time.Minute, s.RefreshObservation)
	go schedule(ctx, 30*time.Minute, s.RefreshForecast)
}
func schedule(ctx context.Context, interval time.Duration, refresh func(context.Context)) {
	scheduleWithTicker(ctx, interval, refresh, func(interval time.Duration) (<-chan time.Time, func()) {
		ticker := time.NewTicker(interval)
		return ticker.C, ticker.Stop
	})
}
func scheduleWithTicker(ctx context.Context, interval time.Duration, refresh func(context.Context), newTicker func(time.Duration) (<-chan time.Time, func())) {
	refresh(ctx)
	for {
		if ctx.Err() != nil {
			return
		}
		ticks, stop := newTicker(interval)
		select {
		case <-ctx.Done():
			stop()
			return
		case <-ticks:
			stop()
			refresh(ctx)
		}
	}
}
