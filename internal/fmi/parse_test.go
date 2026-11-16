package fmi

import (
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

func sample(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile("../../docs/feasibility/evidence/hamina-" + name + ".xml")
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestObservationSample(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 8, 0, 0, time.UTC)
	data := sample(t, "observation")
	o, err := ParseObservation([]byte(data), now)
	if err != nil || o.ValueCM != 32.1 || !o.ObservedAt.Equal(time.Date(2026, 11, 16, 18, 0, 0, 0, time.UTC)) {
		t.Fatalf("%+v %v", o, err)
	}
	for _, bad := range []string{
		strings.ReplaceAll(data, "134254", "999999"),
		strings.ReplaceAll(data, "WLEVN2K_PT1S_INSTANT", "OTHER"),
		strings.ReplaceAll(data, "2026-11-16T18:00:00Z", "not-a-time"),
		strings.ReplaceAll(data, "<wml2:value>321.0</wml2:value>", "<wml2:value>broken</wml2:value>"),
		regexp.MustCompile(`<wml2:value>[^<]*</wml2:value>`).ReplaceAllString(data, "<wml2:value></wml2:value>"),
	} {
		if _, err := ParseObservation([]byte(bad), now); err == nil {
			t.Fatal("invalid observation accepted")
		}
	}
	if _, err := ParseObservation([]byte(data), now.Add(-10*time.Minute)); err == nil {
		t.Fatal("future observation accepted")
	}
}

func TestForecastSample(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 8, 0, 0, time.UTC)
	data := sample(t, "forecast")
	f, err := ParseForecast([]byte(data), now)
	if err != nil {
		t.Fatal(err)
	}
	n, missing := 0, 0
	for _, p := range f.Points {
		if p.ValueCM == nil {
			missing++
		} else {
			n++
		}
	}
	if n != 114 || missing != 7 || f.UnitStatus != "assumed-cm" || f.ResultTime.Format(time.RFC3339) != "2026-11-16T17:00:11Z" || f.AnalysisTime.Format(time.RFC3339) != "2026-11-16T17:00:00Z" {
		t.Fatalf("%d %d %+v", n, missing, f)
	}
	for _, bad := range []string{
		strings.ReplaceAll(data, "SeaLevelN2000", "Other"),
		strings.ReplaceAll(data, "analysisTime", "otherTime"),
		strings.Replace(data, `xlink:href="#time-1-1"`, `xlink:href="#missing"`, 1),
		"<not-xml",
		`<ows:ExceptionReport xmlns:ows="http://www.opengis.net/ows/1.1"/>`,
		regexp.MustCompile(`<wml2:value>[^<]*</wml2:value>`).ReplaceAllString(data, "<wml2:value>NaN</wml2:value>"),
		strings.Replace(data, `<wml2:MeasurementTimeseries gml:id="mts-1-1-SeaLevelN2000">`, `<wml2:MeasurementTimeseries gml:id="mts-1-1-SeaLevelN2000"><wml2:uom code="mm"/>`, 1),
		strings.Replace(data, `<wml2:MeasurementTimeseries gml:id="mts-1-1-SeaLevelN2000">`, `<wml2:MeasurementTimeseries gml:id="mts-1-1-SeaLevelN2000"><swe:uom code="mm"/>`, 1),
	} {
		if _, err := ParseForecast([]byte(bad), now); err == nil {
			t.Fatal("invalid forecast accepted")
		}
	}
	marker := `<wml2:MeasurementTimeseries gml:id="mts-1-1-SeaLevelN2000">`
	start := strings.Index(data, marker)
	if start < 0 {
		t.Fatal("missing N2000 series")
	}
	tail := data[start:]
	re := regexp.MustCompile(`<wml2:value>[^<]*</wml2:value>`)
	loc := re.FindStringIndex(tail)
	badPoint := data[:start] + tail[:loc[0]] + "<wml2:value>not-a-number</wml2:value>" + tail[loc[1]:]
	if _, err := ParseForecast([]byte(badPoint), now); err == nil {
		t.Fatal("invalid number accepted")
	}
	loc2 := re.FindStringIndex(tail[loc[1]:])
	pos := loc[1] + loc2[0]
	gap := data[:start] + tail[:pos] + re.ReplaceAllString(tail[pos:pos+loc2[1]-loc2[0]], "<wml2:value>NaN</wml2:value>") + tail[loc[1]+loc2[1]:]
	g, err := ParseForecast([]byte(gap), now)
	if err != nil || g.Points[1].ValueCM != nil {
		t.Fatalf("internal gap: %v", err)
	}
}

func TestUnresolvedPhenomenonTimeRejected(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 8, 0, 0, time.UTC)
	t.Run("forecast", func(t *testing.T) {
		forecast := sample(t, "forecast")
		bad := strings.Replace(forecast, `xlink:href="#time-interval-1-1"`, `xlink:href="#missing"`, 1)
		if bad == forecast {
			t.Fatal("fixture reference missing")
		}
		if _, err := ParseForecast([]byte(bad), now); err == nil {
			t.Fatal("unresolved forecast phenomenonTime accepted")
		}
	})
	t.Run("observation", func(t *testing.T) {
		observation := sample(t, "observation")
		bad := strings.Replace(observation, "<om:phenomenonTime>", `<om:phenomenonTime xlink:href="#missing">`, 1)
		if bad == observation {
			t.Fatal("fixture phenomenonTime missing")
		}
		if _, err := ParseObservation([]byte(bad), now); err == nil {
			t.Fatal("unresolved observation phenomenonTime accepted")
		}
	})
	t.Run("observation result time", func(t *testing.T) {
		observation := sample(t, "observation")
		bad := strings.Replace(observation, "<om:resultTime>", `<om:resultTime xlink:href="#missing">`, 1)
		if bad == observation {
			t.Fatal("fixture resultTime missing")
		}
		if _, err := ParseObservation([]byte(bad), now); err == nil {
			t.Fatal("unresolved observation resultTime accepted")
		}
	})
}
