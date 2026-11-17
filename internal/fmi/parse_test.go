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

func TestUnitMetadata(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 8, 0, 0, time.UTC)
	for _, series := range []struct {
		name, marker, expected, other, href, otherHref string
	}{
		{"observation", `<wml2:MeasurementTimeseries gml:id="obs-obs-1-1-WLEVN2K_PT1S_INSTANT">`, "mm", "cm", "http://qudt.org/vocab/unit/MilliM", "http://qudt.org/vocab/unit/CentiM"},
		{"forecast", `<wml2:MeasurementTimeseries gml:id="mts-1-1-SeaLevelN2000">`, "cm", "mm", "http://qudt.org/vocab/unit/CentiM", "http://qudt.org/vocab/unit/MilliM"},
	} {
		t.Run(series.name, func(t *testing.T) {
			original := sample(t, series.name)
			for _, tc := range []struct {
				name, unit string
				valid      bool
			}{
				{"absent", "", true},
				{"code", `<wml2:uom code="` + series.expected + `"/>`, true},
				{"uom attribute", `<swe:uom uom="` + series.expected + `"/>`, true},
				{"href", `<wml2:uom xlink:href="` + series.href + `"/>`, true},
				{"unitOfMeasure href", `<wml2:unitOfMeasure xlink:href="` + series.href + `"/>`, true},
				{"conflicting code", `<wml2:uom code="` + series.other + `"/>`, false},
				{"conflicting href", `<wml2:uom xlink:href="` + series.otherHref + `"/>`, false},
				{"unknown href", `<wml2:uom xlink:href="https://example.invalid/unit"/>`, false},
				{"conflicting declarations", `<wml2:uom code="` + series.other + `" xlink:href="` + series.href + `"/>`, false},
				{"conflicting elements", `<wml2:uom code="` + series.expected + `"/><swe:uom code="` + series.other + `"/>`, false},
				{"unknown code", `<swe:uom code="unknown"/>`, false},
				{"empty declaration", `<wml2:uom/>`, false},
			} {
				t.Run(tc.name, func(t *testing.T) {
					xml := strings.Replace(original, series.marker, series.marker+tc.unit, 1)
					if tc.unit != "" && xml == original {
						t.Fatal("fixture marker not found")
					}
					var err error
					if series.name == "observation" {
						_, err = ParseObservation([]byte(xml), now)
					} else {
						_, err = ParseForecast([]byte(xml), now)
					}
					if (err == nil) != tc.valid {
						t.Fatalf("valid=%v err=%v", tc.valid, err)
					}
				})
			}
		})
	}
}

func TestUnitMetadataNamespaceMatrix(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 8, 0, 0, time.UTC)
	for _, tc := range []struct {
		name, marker, expected, other string
	}{
		{"observation", `<wml2:MeasurementTimeseries gml:id="obs-obs-1-1-WLEVN2K_PT1S_INSTANT">`, "mm", "cm"},
		{"forecast", `<wml2:MeasurementTimeseries gml:id="mts-1-1-SeaLevelN2000">`, "cm", "mm"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			original := sample(t, tc.name)
			for _, unit := range []struct {
				name, xml string
				valid     bool
			}{
				{"omop expected", `<omop:uom xmlns:omop="http://inspire.ec.europa.eu/schemas/omop/2.9" uom="` + tc.expected + `"/>`, true},
				{"omop conflicting", `<omop:uom xmlns:omop="http://inspire.ec.europa.eu/schemas/omop/2.9" uom="` + tc.other + `"/>`, false},
				{"unknown namespace", `<unknown:uom xmlns:unknown="https://example.invalid/units" code="` + tc.expected + `"/>`, false},
			} {
				t.Run(unit.name, func(t *testing.T) {
					modified := strings.Replace(original, tc.marker, tc.marker+unit.xml, 1)
					if modified == original {
						t.Fatal("fixture marker not found")
					}
					var err error
					if tc.name == "observation" {
						_, err = ParseObservation([]byte(modified), now)
					} else {
						_, err = ParseForecast([]byte(modified), now)
					}
					if (err == nil) != unit.valid {
						t.Fatalf("valid=%v error=%v", unit.valid, err)
					}
				})
			}
		})
	}
}

func TestTimePeriodNestedReferenceMatrix(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 8, 0, 0, time.UTC)
	for _, tc := range []struct{ name, marker string }{
		{"observation", `<gml:TimePeriod  gml:id="time1-1-1">`},
		{"forecast", `<gml:TimePeriod gml:id="time-interval-1-1">`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			original := sample(t, tc.name)
			mutated := strings.Replace(original, tc.marker, strings.TrimSuffix(tc.marker, ">")+` xlink:href="#missing">`, 1)
			if mutated == original {
				t.Fatal("fixture marker not found")
			}
			var err error
			if tc.name == "observation" {
				_, err = ParseObservation([]byte(mutated), now)
			} else {
				_, err = ParseForecast([]byte(mutated), now)
			}
			if err == nil {
				t.Fatal("unresolved nested TimePeriod reference accepted")
			}
		})
	}
}

func TestTimePeriodBoundaryReferences(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 8, 0, 0, time.UTC)
	begin := regexp.MustCompile(`<gml:beginPosition>[^<]+</gml:beginPosition>`)
	end := regexp.MustCompile(`<gml:endPosition>[^<]+</gml:endPosition>`)
	for _, tc := range []struct{ name, resultID string }{
		{"observation", "time2-1-1"},
		{"forecast", "time-1-1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			original := sample(t, tc.name)
			for _, boundary := range []struct {
				name, begin, end string
				valid            bool
			}{
				{"inline instants", `<gml:begin><gml:TimeInstant gml:id="boundary-begin"><gml:timePosition>2026-11-16T17:00:00Z</gml:timePosition></gml:TimeInstant></gml:begin>`, `<gml:end><gml:TimeInstant gml:id="boundary-end"><gml:timePosition>2026-11-16T18:00:00Z</gml:timePosition></gml:TimeInstant></gml:end>`, true},
				{"internal references", `<gml:begin xlink:href="#` + tc.resultID + `"/>`, `<gml:end xlink:href="#` + tc.resultID + `"/>`, true},
				{"unresolved reference", `<gml:begin xlink:href="#missing"/>`, `<gml:end xlink:href="#` + tc.resultID + `"/>`, false},
				{"external reference", `<gml:begin xlink:href="https://example.invalid/time"/>`, `<gml:end xlink:href="#` + tc.resultID + `"/>`, false},
				{"empty reference", `<gml:begin xlink:href=""/>`, `<gml:end xlink:href="#` + tc.resultID + `"/>`, false},
				{"reference and inline", `<gml:begin xlink:href="#` + tc.resultID + `"><gml:TimeInstant gml:id="extra-begin"><gml:timePosition>2026-11-16T17:00:00Z</gml:timePosition></gml:TimeInstant></gml:begin>`, `<gml:end xlink:href="#` + tc.resultID + `"/>`, false},
				{"two begin boundaries", `<gml:beginPosition>2026-11-16T17:00:00Z</gml:beginPosition><gml:begin xlink:href="#` + tc.resultID + `"/>`, `<gml:end xlink:href="#` + tc.resultID + `"/>`, false},
				{"reversed boundaries", `<gml:beginPosition>2026-11-16T18:00:00Z</gml:beginPosition>`, `<gml:endPosition>2026-11-16T17:00:00Z</gml:endPosition>`, false},
			} {
				t.Run(boundary.name, func(t *testing.T) {
					if !begin.MatchString(original) || !end.MatchString(original) {
						t.Fatal("fixture boundaries missing")
					}
					mutated := begin.ReplaceAllString(original, boundary.begin)
					mutated = end.ReplaceAllString(mutated, boundary.end)
					var err error
					if tc.name == "observation" {
						_, err = ParseObservation([]byte(mutated), now)
					} else {
						_, err = ParseForecast([]byte(mutated), now)
					}
					if (err == nil) != boundary.valid {
						t.Fatalf("valid=%v error=%v", boundary.valid, err)
					}
				})
			}
		})
	}
}

func TestPointValueUnitAttributes(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 8, 0, 0, time.UTC)
	for _, tc := range []struct{ name, marker, expected, other string }{
		{"observation", `<wml2:MeasurementTimeseries gml:id="obs-obs-1-1-WLEVN2K_PT1S_INSTANT">`, "mm", "cm"},
		{"forecast", `<wml2:MeasurementTimeseries gml:id="mts-1-1-SeaLevelN2000">`, "cm", "mm"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			original := sample(t, tc.name)
			start := strings.Index(original, tc.marker)
			if start < 0 {
				t.Fatal("fixture series missing")
			}
			for _, value := range []struct {
				name, unit string
				valid      bool
			}{
				{"expected", tc.expected, true}, {"conflicting", tc.other, false}, {"unknown", "unknown", false}, {"empty", "", false},
			} {
				t.Run(value.name, func(t *testing.T) {
					head, tail := original[:start], original[start:]
					changed := strings.Replace(tail, "<wml2:value>", `<wml2:value uom="`+value.unit+`">`, 1)
					if changed == tail {
						t.Fatal("fixture value missing")
					}
					var err error
					if tc.name == "observation" {
						_, err = ParseObservation([]byte(head+changed), now)
					} else {
						_, err = ParseForecast([]byte(head+changed), now)
					}
					if (err == nil) != value.valid {
						t.Fatalf("valid=%v error=%v", value.valid, err)
					}
				})
			}
		})
	}
}

func TestIndeterminateTimesRejected(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 8, 0, 0, time.UTC)
	for _, tc := range []struct{ name, marker, replacement string }{
		{"observation value time", "<wml2:time>", `<wml2:time indeterminatePosition="before">`},
		{"forecast value time", "<wml2:time>", `<wml2:time indeterminatePosition="before">`},
		{"observation result time", "<gml:timePosition>", `<gml:timePosition indeterminatePosition="before">`},
		{"forecast period beginning", "<gml:beginPosition>", `<gml:beginPosition indeterminatePosition="before">`},
		{"observation empty indeterminate", "<wml2:time>", `<wml2:time indeterminatePosition="">`},
		{"forecast foreign time frame", "<wml2:time>", `<wml2:time frame="https://example.invalid/frame">`},
		{"observation value time reference", "<wml2:time>", `<wml2:time xlink:href="#missing">`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			name := "observation"
			if strings.HasPrefix(tc.name, "forecast") {
				name = "forecast"
			}
			original := sample(t, name)
			mutated := ""
			if tc.name == "forecast value time" || tc.name == "forecast foreign time frame" {
				start := strings.Index(original, `<wml2:MeasurementTimeseries gml:id="mts-1-1-SeaLevelN2000">`)
				if start < 0 {
					t.Fatal("N2000 series missing")
				}
				mutated = original[:start] + strings.Replace(original[start:], tc.marker, tc.replacement, 1)
			} else {
				mutated = strings.Replace(original, tc.marker, tc.replacement, 1)
			}
			if mutated == original {
				t.Fatal("fixture marker missing")
			}
			var err error
			if name == "observation" {
				_, err = ParseObservation([]byte(mutated), now)
			} else {
				_, err = ParseForecast([]byte(mutated), now)
			}
			if err == nil {
				t.Fatal("indeterminate time accepted as exact")
			}
		})
	}
}

func TestExplicitISOFrameAccepted(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 8, 0, 0, time.UTC)
	original := sample(t, "observation")
	mutated := strings.Replace(original, "<wml2:time>", `<wml2:time frame="#ISO-8601">`, 1)
	if mutated == original {
		t.Fatal("fixture time missing")
	}
	if _, err := ParseObservation([]byte(mutated), now); err != nil {
		t.Fatal(err)
	}
}

func TestWaterMLUnitMetadataLocations(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 8, 0, 0, time.UTC)
	for _, tc := range []struct{ name, marker, expected, other string }{
		{"observation", `<wml2:MeasurementTimeseries gml:id="obs-obs-1-1-WLEVN2K_PT1S_INSTANT">`, "mm", "cm"},
		{"forecast", `<wml2:MeasurementTimeseries gml:id="mts-1-1-SeaLevelN2000">`, "cm", "mm"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			original := sample(t, tc.name)
			start := strings.Index(original, tc.marker)
			if start < 0 {
				t.Fatal("fixture series missing")
			}
			head, tail := original[:start], original[start:]
			for _, placement := range []struct {
				name  string
				apply func(string, string) string
			}{
				{"series default", func(text, unit string) string {
					metadata := `<wml2:defaultPointMetadata><wml2:DefaultTVPMeasurementMetadata><wml2:uom code="` + unit + `"/></wml2:DefaultTVPMeasurementMetadata></wml2:defaultPointMetadata>`
					return strings.Replace(text, tc.marker, tc.marker+metadata, 1)
				}},
				{"point metadata", func(text, unit string) string {
					metadata := `</wml2:value><wml2:metadata><wml2:TVPMeasurementMetadata><wml2:uom code="` + unit + `"/></wml2:TVPMeasurementMetadata></wml2:metadata>`
					return strings.Replace(text, "</wml2:value>", metadata, 1)
				}},
			} {
				for _, value := range []struct {
					name, code string
					valid      bool
				}{{"expected", tc.expected, true}, {"conflicting", tc.other, false}, {"unknown", "unknown", false}} {
					t.Run(placement.name+"/"+value.name, func(t *testing.T) {
						mutated := placement.apply(tail, value.code)
						if mutated == tail {
							t.Fatal("fixture marker missing")
						}
						var err error
						if tc.name == "observation" {
							_, err = ParseObservation([]byte(head+mutated), now)
						} else {
							_, err = ParseForecast([]byte(head+mutated), now)
						}
						if (err == nil) != value.valid {
							t.Fatalf("valid=%v error=%v", value.valid, err)
						}
					})
				}
			}
		})
	}
}

func TestOtherForecastSeriesUnitDoesNotAffectN2000(t *testing.T) {
	now := time.Date(2026, 11, 16, 18, 8, 0, 0, time.UTC)
	original := sample(t, "forecast")
	marker := `<wml2:MeasurementTimeseries gml:id="mts-1-1-SeaLevel">`
	modified := strings.Replace(original, marker, marker+`<wml2:uom code="mm"/>`, 1)
	if modified == original {
		t.Fatal("SeaLevel fixture marker missing")
	}
	f, err := ParseForecast([]byte(modified), now)
	if err != nil || f.UnitStatus != "assumed-cm" {
		t.Fatalf("other series affected N2000: %+v %v", f, err)
	}
}
