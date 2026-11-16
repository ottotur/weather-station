// Package cache stores the two independently versioned, validated datasets.
package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"weather-station/internal/waterlevel"
)

type Store struct{ Dir string }

type envelope[T any] struct {
	Version int `json:"version"`
	Data    T   `json:"data"`
}

func (s Store) write(name string, value any) (err error) {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(s.Dir, "."+name+"-*.tmp")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer func() { _ = f.Close(); _ = os.Remove(tmp) }()
	if _, err = f.Write(data); err != nil {
		return err
	}
	if err = f.Sync(); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	if err = os.Rename(tmp, filepath.Join(s.Dir, name)); err != nil {
		return err
	}
	dir, err := os.Open(s.Dir)
	if err != nil {
		return err
	}
	defer dir.Close()
	return dir.Sync()
}

func read[T any](s Store, name string) (T, error) {
	var zero T
	b, err := os.ReadFile(filepath.Join(s.Dir, name))
	if err != nil {
		return zero, err
	}
	var e envelope[T]
	if err = json.Unmarshal(b, &e); err != nil {
		return zero, err
	}
	if e.Version != 1 {
		return zero, errors.New("unsupported cache version")
	}
	return e.Data, nil
}

func validObservation(o *waterlevel.Observation) error {
	if o == nil || o.ObservedAt.IsZero() || o.FetchedAt.IsZero() || o.ObservedAt.After(o.FetchedAt) || math.IsNaN(o.ValueCM) || math.IsInf(o.ValueCM, 0) {
		return errors.New("invalid cached observation")
	}
	return nil
}
func validForecast(f *waterlevel.Forecast) error {
	if f == nil || f.AnalysisTime.IsZero() || f.ResultTime.IsZero() || f.FetchedAt.IsZero() || f.AnalysisTime.After(f.FetchedAt) || f.UnitStatus != "assumed-cm" || len(f.Points) == 0 {
		return errors.New("invalid cached forecast")
	}
	var previous time.Time
	finite := false
	for _, p := range f.Points {
		if p.Time.IsZero() || (!previous.IsZero() && !p.Time.After(previous)) {
			return errors.New("invalid cached forecast time")
		}
		previous = p.Time
		if p.ValueCM != nil {
			if math.IsNaN(*p.ValueCM) || math.IsInf(*p.ValueCM, 0) {
				return errors.New("invalid cached value")
			}
			finite = true
		}
	}
	if !finite {
		return errors.New("no cached forecast values")
	}
	return nil
}

func (s Store) SaveObservation(o *waterlevel.Observation) error {
	if err := validObservation(o); err != nil {
		return err
	}
	return s.write("observation.json", envelope[*waterlevel.Observation]{Version: 1, Data: o})
}
func (s Store) SaveForecast(f *waterlevel.Forecast) error {
	if err := validForecast(f); err != nil {
		return err
	}
	return s.write("forecast.json", envelope[*waterlevel.Forecast]{Version: 1, Data: f})
}
func (s Store) LoadObservation() (*waterlevel.Observation, error) {
	o, err := read[*waterlevel.Observation](s, "observation.json")
	if err != nil {
		return nil, fmt.Errorf("observation cache: %w", err)
	}
	if err = validObservation(o); err != nil {
		return nil, err
	}
	return o, nil
}
func (s Store) LoadForecast() (*waterlevel.Forecast, error) {
	f, err := read[*waterlevel.Forecast](s, "forecast.json")
	if err != nil {
		return nil, fmt.Errorf("forecast cache: %w", err)
	}
	if err = validForecast(f); err != nil {
		return nil, err
	}
	return f, nil
}
