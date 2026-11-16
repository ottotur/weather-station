package fmi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"weather-station/internal/waterlevel"
)

const maxResponseBytes = 2 << 20

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL, HTTP: &http.Client{Timeout: 30 * time.Second}}
}

func (c *Client) get(ctx context.Context, query string, now time.Time, forecast bool) ([]byte, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("service", "WFS")
	q.Set("version", "2.0.0")
	q.Set("request", "GetFeature")
	q.Set("storedquery_id", query)
	q.Set("fmisid", "134254")
	q.Set("timezone", "UTC")
	start, end := now.Add(-time.Hour), now
	if forecast {
		start = now
		end = now.Add(48 * time.Hour)
		q.Set("timestep", "60")
		q.Del("parameters")
	} else {
		q.Set("timestep", "1")
		q.Set("parameters", "WLEVN2K_PT1S_INSTANT")
	}
	q.Set("starttime", start.UTC().Format(time.RFC3339))
	q.Set("endtime", end.UTC().Format(time.RFC3339))
	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("FMI HTTP status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxResponseBytes {
		return nil, fmt.Errorf("FMI response exceeds limit")
	}
	return data, nil
}

func (c *Client) Observation(ctx context.Context, now time.Time) (*waterlevel.Observation, error) {
	data, err := c.get(ctx, "fmi::observations::mareograph::instant::timevaluepair", now, false)
	if err != nil {
		return nil, err
	}
	return ParseObservation(data, now)
}

func (c *Client) Forecast(ctx context.Context, now time.Time) (*waterlevel.Forecast, error) {
	data, err := c.get(ctx, "fmi::forecast::sealevel::point::timevaluepair", now, true)
	if err != nil {
		return nil, err
	}
	return ParseForecast(data, now)
}
