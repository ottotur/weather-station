package waterlevel

import (
	"testing"
	"time"
)

func TestFreshnessAndCoverage(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		age  time.Duration
		want Status
	}{
		{29 * time.Minute, Fresh}, {30 * time.Minute, Fresh}, {30*time.Minute + time.Nanosecond, Stale},
	} {
		o := &Observation{ObservedAt: now.Add(-tc.age)}
		if got := ObservationStatus(o, now); got != tc.want {
			t.Errorf("observation age %v: %s", tc.age, got)
		}
	}
	for _, tc := range []struct {
		age  time.Duration
		want Status
	}{
		{11 * time.Hour, Fresh}, {12 * time.Hour, Fresh}, {12*time.Hour + time.Nanosecond, Stale},
	} {
		f := &Forecast{AnalysisTime: now.Add(-tc.age), FetchedAt: now}
		if got := ForecastStatus(f, now); got != tc.want {
			t.Errorf("forecast age %v: %s", tc.age, got)
		}
	}
	v := 42.0
	f := &Forecast{AnalysisTime: now.Add(-13 * time.Hour), FetchedAt: now, Points: []Point{{Time: now.Add(time.Hour), ValueCM: &v}}}
	if ForecastStatus(f, now) != Stale || !HasFutureValues(f, now) {
		t.Fatal("refetch must not change model age or future coverage")
	}
	f.Points[0].Time = now.Add(-time.Second)
	if HasFutureValues(f, now) {
		t.Fatal("ended forecast has no future values")
	}
	if ObservationStatus(nil, now) != Missing || ForecastStatus(nil, now) != Missing {
		t.Fatal("missing status")
	}
}

func TestSummary(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 0, 0, 0, time.UTC)
	o := &Observation{ValueCM: 32, ObservedAt: now}
	a, b, c := 50.0, 20.0, 50.0
	f := &Forecast{AnalysisTime: now, Points: []Point{
		{Time: now.Add(3 * time.Hour), ValueCM: &c}, {Time: now.Add(2 * time.Hour), ValueCM: &b},
		{Time: now.Add(time.Hour), ValueCM: &a}, {Time: now.Add(90 * time.Minute)},
	}}
	s, reason := CalculateSummary(o, f, now)
	if reason != "" || s.Rise.ChangeCM != 18 || s.Fall.ChangeCM != 12 || !s.HasGaps || !s.Rise.At.Equal(now.Add(time.Hour)) {
		t.Fatalf("unexpected summary: %+v %q", s, reason)
	}
	for _, tc := range []struct {
		value      float64
		rise, fall float64
	}{{40, 8, 0}, {20, 0, 12}, {32, 0, 0}} {
		v := tc.value
		f.Points = []Point{{Time: now, ValueCM: &v}}
		s, _ := CalculateSummary(o, f, now)
		if s.Rise.ChangeCM != tc.rise || s.Fall.ChangeCM != tc.fall || (tc.rise == 0 && s.Rise.At != nil) || (tc.fall == 0 && s.Fall.At != nil) {
			t.Fatalf("value %v: %+v", v, s)
		}
	}
	o.ObservedAt = now.Add(-31 * time.Minute)
	if s, reason := CalculateSummary(o, f, now); s != nil || reason != "observation_stale" {
		t.Fatal(s, reason)
	}
	o.ObservedAt = now
	f.AnalysisTime = now.Add(-13 * time.Hour)
	if s, reason := CalculateSummary(o, f, now); s != nil || reason != "forecast_stale" {
		t.Fatal(s, reason)
	}
	f.AnalysisTime = now
	f.Points = nil
	if s, reason := CalculateSummary(o, f, now); s != nil || reason != "no_future_values" {
		t.Fatal(s, reason)
	}
}

func TestSummaryRangeAndTie(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 0, 0, 0, time.UTC)
	o := &Observation{ValueCM: 32, ObservedAt: now}
	low, high := 20.0, 100.0
	f := &Forecast{AnalysisTime: now, Points: []Point{
		{Time: now.Add(-time.Second), ValueCM: &high},
		{Time: now.Add(49 * time.Hour), ValueCM: &high},
		{Time: now.Add(48 * time.Hour), ValueCM: &low},
		{Time: now.Add(time.Hour), ValueCM: &low},
		{Time: now.Add(47 * time.Hour)},
	}}
	s, reason := CalculateSummary(o, f, now)
	if reason != "" || s.Rise.ChangeCM != 0 || s.Fall.ChangeCM != 12 || !s.Fall.At.Equal(now.Add(time.Hour)) || !s.HasGaps {
		t.Fatalf("unexpected summary: %+v %q", s, reason)
	}
	f.Points = []Point{{Time: now, ValueCM: &low}, {Time: now.Add(time.Hour)}}
	s, _ = CalculateSummary(o, f, now)
	if s.HasGaps {
		t.Fatal("missing tail is not an internal gap")
	}
}

func TestAbsentHourIsGap(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 0, 0, 0, time.UTC)
	value := 40.0
	f := &Forecast{AnalysisTime: now, Points: []Point{
		{Time: now, ValueCM: &value},
		{Time: now.Add(2 * time.Hour), ValueCM: &value},
	}}
	p := FuturePoints(f, now)
	if len(p) != 3 || p[1].ValueCM != nil || !p[1].Time.Equal(now.Add(time.Hour)) {
		t.Fatalf("missing hour: %+v", p)
	}
	s, _ := CalculateSummary(&Observation{ObservedAt: now, ValueCM: 32}, f, now)
	if s == nil || !s.HasGaps {
		t.Fatal("missing hour not reported")
	}
}
