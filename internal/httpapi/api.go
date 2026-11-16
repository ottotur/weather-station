package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"weather-station/internal/waterlevel"
)

type API struct {
	Service *waterlevel.Service
	Now     func() time.Time
}

type forecastMetadata struct {
	AnalysisTime time.Time  `json:"analysisTime"`
	ResultTime   time.Time  `json:"resultTime"`
	FetchedAt    time.Time  `json:"fetchedAt"`
	UnitStatus   string     `json:"unitStatus"`
	LastValidAt  *time.Time `json:"lastValidAt"`
}

func (a API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	now := a.Now().UTC()
	s := a.Service.Snapshot()
	summary, reason := waterlevel.CalculateSummary(s.Observation, s.Forecast, now)
	points := waterlevel.FuturePoints(s.Forecast, now)
	if points == nil {
		points = []waterlevel.Point{}
	}
	var metadata *forecastMetadata
	var lastShown *time.Time
	if s.Forecast != nil {
		metadata = &forecastMetadata{AnalysisTime: s.Forecast.AnalysisTime, ResultTime: s.Forecast.ResultTime, FetchedAt: s.Forecast.FetchedAt, UnitStatus: s.Forecast.UnitStatus}
		for _, point := range s.Forecast.Points {
			if point.ValueCM != nil && (metadata.LastValidAt == nil || point.Time.After(*metadata.LastValidAt)) {
				t := point.Time
				metadata.LastValidAt = &t
			}
		}
	}
	for _, point := range points {
		if point.ValueCM != nil {
			t := point.Time
			lastShown = &t
		}
	}
	response := struct {
		SchemaVersion int            `json:"schemaVersion"`
		ServerTime    time.Time      `json:"serverTime"`
		Station       map[string]any `json:"station"`
		Unit          string         `json:"unit"`
		Datum         string         `json:"datum"`
		Observation   any            `json:"observation"`
		Forecast      any            `json:"forecast"`
		Summary       any            `json:"summary"`
		SummaryReason string         `json:"summaryReason"`
		Source        any            `json:"source"`
	}{
		SchemaVersion: 1, ServerTime: now, Station: map[string]any{"name": "Hamina Pitäjänsaari", "fmisid": 134254, "geoid": -10022823}, Unit: "cm", Datum: "N2000",
		Observation: struct {
			Status  waterlevel.Status       `json:"status"`
			Data    *waterlevel.Observation `json:"data"`
			Attempt waterlevel.Attempt      `json:"attempt"`
		}{waterlevel.ObservationStatus(s.Observation, now), s.Observation, s.ObservationAttempt},
		Forecast: struct {
			Status          waterlevel.Status  `json:"status"`
			Data            *forecastMetadata  `json:"data"`
			Points          []waterlevel.Point `json:"points"`
			HasFutureValues bool               `json:"hasFutureValues"`
			LastShownAt     *time.Time         `json:"lastShownAt"`
			Attempt         waterlevel.Attempt `json:"attempt"`
		}{waterlevel.ForecastStatus(s.Forecast, now), metadata, points, waterlevel.HasFutureValues(s.Forecast, now), lastShown, s.ForecastAttempt},
		Summary: summary, SummaryReason: reason,
		Source: map[string]string{"name": "Ilmatieteen laitos", "licence": "CC BY 4.0", "licenceUrl": "https://creativecommons.org/licenses/by/4.0/", "forecastUnit": "assumed-cm", "processing": "Havainto muunnettu millimetreistä senttimetreiksi; muutokset laskettu ennusteesta."},
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(response)
}
