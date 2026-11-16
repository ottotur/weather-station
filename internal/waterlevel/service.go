package waterlevel

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrDurabilityUncertain means the file was replaced but directory sync failed.
var ErrDurabilityUncertain = errors.New("cache durability uncertain")

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
	if !s.observationMu.TryLock() {
		return
	}
	defer s.observationMu.Unlock()
	now := s.now().UTC()
	o, err := s.source.Observation(ctx, now)
	code := "source_error"
	if err == nil {
		err = s.store.SaveObservation(o)
		code = "cache_error"
		if errors.Is(err, ErrDurabilityUncertain) {
			code = "cache_durability_uncertain"
		}
	}
	attempt := Attempt{At: &now, Result: "success"}
	if err != nil {
		attempt.Result = "error"
		attempt.ErrorCode = &code
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.ObservationAttempt = attempt
	if err == nil {
		s.state.Observation = o
	}
}

func (s *Service) RefreshForecast(ctx context.Context) {
	if !s.forecastMu.TryLock() {
		return
	}
	defer s.forecastMu.Unlock()
	now := s.now().UTC()
	f, err := s.source.Forecast(ctx, now)
	code := "source_error"
	if err == nil {
		err = s.store.SaveForecast(f)
		code = "cache_error"
		if errors.Is(err, ErrDurabilityUncertain) {
			code = "cache_durability_uncertain"
		}
	}
	attempt := Attempt{At: &now, Result: "success"}
	if err != nil {
		attempt.Result = "error"
		attempt.ErrorCode = &code
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.ForecastAttempt = attempt
	if err == nil {
		s.state.Forecast = f
	}
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
	ticks, stop := newTicker(interval)
	defer stop()
	scheduleTicks(ctx, ticks, refresh)
}

func scheduleTicks(ctx context.Context, ticks <-chan time.Time, refresh func(context.Context)) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticks:
			refresh(ctx)
		}
	}
}
