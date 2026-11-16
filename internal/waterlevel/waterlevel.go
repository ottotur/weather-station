// Package waterlevel contains the source-independent water level rules.
package waterlevel

import (
	"math"
	"sort"
	"time"
)

type Observation struct {
	ValueCM    float64   `json:"valueCm"`
	ObservedAt time.Time `json:"observedAt"`
	FetchedAt  time.Time `json:"fetchedAt"`
}

type Point struct {
	Time    time.Time `json:"time"`
	ValueCM *float64  `json:"valueCm"`
}

type Forecast struct {
	AnalysisTime time.Time `json:"analysisTime"`
	ResultTime   time.Time `json:"resultTime"`
	FetchedAt    time.Time `json:"fetchedAt"`
	UnitStatus   string    `json:"unitStatus"`
	Points       []Point   `json:"points"`
}

type Status string

const (
	Missing Status = "missing"
	Fresh   Status = "fresh"
	Stale   Status = "stale"
)

func ObservationStatus(o *Observation, now time.Time) Status {
	if o == nil {
		return Missing
	}
	if now.Sub(o.ObservedAt) > 30*time.Minute {
		return Stale
	}
	return Fresh
}

func ForecastStatus(f *Forecast, now time.Time) Status {
	if f == nil {
		return Missing
	}
	if now.Sub(f.AnalysisTime) > 12*time.Hour {
		return Stale
	}
	return Fresh
}

// FuturePoints includes missing values so callers can display gaps rather than
// draw a line across them. Its result is ordered independently of input order.
func FuturePoints(f *Forecast, now time.Time) []Point {
	if f == nil {
		return nil
	}
	var result []Point
	for _, p := range f.Points {
		if !p.Time.Before(now) && !p.Time.After(now.Add(48*time.Hour)) {
			result = append(result, p)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Time.Before(result[j].Time) })
	return result
}

func HasFutureValues(f *Forecast, now time.Time) bool {
	for _, p := range FuturePoints(f, now) {
		if p.ValueCM != nil && !math.IsNaN(*p.ValueCM) && !math.IsInf(*p.ValueCM, 0) {
			return true
		}
	}
	return false
}

type Extremum struct {
	ChangeCM float64    `json:"changeCm"`
	At       *time.Time `json:"at"`
}

type Summary struct {
	Rise    Extremum `json:"rise"`
	Fall    Extremum `json:"fall"`
	HasGaps bool     `json:"hasGaps"`
}

// CalculateSummary returns a reason when a reliable comparison is unavailable.
func CalculateSummary(o *Observation, f *Forecast, now time.Time) (*Summary, string) {
	if ObservationStatus(o, now) != Fresh {
		return nil, "observation_" + string(ObservationStatus(o, now))
	}
	if ForecastStatus(f, now) != Fresh {
		return nil, "forecast_" + string(ForecastStatus(f, now))
	}
	points := FuturePoints(f, now)
	var summary Summary
	found := false
	gapAfterValue := false
	for _, p := range points {
		if p.ValueCM == nil || math.IsNaN(*p.ValueCM) || math.IsInf(*p.ValueCM, 0) {
			if found {
				gapAfterValue = true
			}
			continue
		}
		if gapAfterValue {
			summary.HasGaps = true
		}
		found = true
		delta := *p.ValueCM - o.ValueCM
		if delta > summary.Rise.ChangeCM {
			at := p.Time
			summary.Rise = Extremum{ChangeCM: delta, At: &at}
		}
		if -delta > summary.Fall.ChangeCM {
			at := p.Time
			summary.Fall = Extremum{ChangeCM: -delta, At: &at}
		}
	}
	if !found {
		return nil, "no_future_values"
	}
	return &summary, ""
}
