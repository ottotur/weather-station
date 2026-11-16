// Package fmi interprets FMI WaterML responses without fetching external XML references.
package fmi

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"weather-station/internal/waterlevel"
)

const (
	wfs   = "http://www.opengis.net/wfs/2.0"
	om    = "http://www.opengis.net/om/2.0"
	omso  = "http://inspire.ec.europa.eu/schemas/omso/3.0"
	gml   = "http://www.opengis.net/gml/3.2"
	wml   = "http://www.opengis.net/waterml/2.0"
	swe   = "http://www.opengis.net/swe/2.0"
	xlink = "http://www.w3.org/1999/xlink"
)

type element struct {
	Name     xml.Name
	Attrs    []xml.Attr
	Text     string
	Children []*element
}

func (e *element) attr(space, local string) string {
	for _, a := range e.Attrs {
		if a.Name.Space == space && a.Name.Local == local {
			return a.Value
		}
	}
	return ""
}
func (e *element) all(space, local string) []*element {
	var out []*element
	var visit func(*element)
	visit = func(n *element) {
		if n.Name.Space == space && n.Name.Local == local {
			out = append(out, n)
		}
		for _, c := range n.Children {
			visit(c)
		}
	}
	visit(e)
	return out
}
func one(e *element, space, local string) (*element, error) {
	found := e.all(space, local)
	if len(found) != 1 {
		return nil, fmt.Errorf("expected one %s, got %d", local, len(found))
	}
	return found[0], nil
}

func document(data []byte) (*element, error) {
	dec := xml.NewDecoder(bytes.NewReader(data))
	var root *element
	var stack []*element
	for {
		token, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("invalid XML: %w", err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			n := &element{Name: t.Name, Attrs: t.Attr}
			if len(stack) == 0 {
				if root != nil {
					return nil, errors.New("multiple roots")
				}
				root = n
			} else {
				p := stack[len(stack)-1]
				p.Children = append(p.Children, n)
			}
			stack = append(stack, n)
		case xml.CharData:
			if len(stack) > 0 {
				stack[len(stack)-1].Text += string(t)
			} else if strings.TrimSpace(string(t)) != "" {
				return nil, errors.New("text outside root")
			}
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		case xml.Directive:
			return nil, errors.New("XML directives not supported")
		}
	}
	if root == nil || root.Name.Space != wfs || root.Name.Local != "FeatureCollection" {
		return nil, errors.New("not a WFS feature collection")
	}
	return root, nil
}

func utc(s string) (time.Time, error) {
	if !strings.HasSuffix(s, "Z") {
		return time.Time{}, errors.New("time must be UTC")
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time: %w", err)
	}
	return t, nil
}

func selectSeries(root *element, property, station, stationCode string) (*element, error) {
	var matches []*element
	for _, series := range root.all(omso, "PointTimeSeriesObservation") {
		prop, err := one(series, om, "observedProperty")
		if err != nil {
			return nil, err
		}
		u, err := url.Parse(prop.attr(xlink, "href"))
		if err != nil {
			return nil, err
		}
		if u.Query().Get("param") != property {
			continue
		}
		ids := series.all(gml, "identifier")
		found := false
		for _, id := range ids {
			if id.attr("", "codeSpace") == "http://xml.fmi.fi/namespace/stationcode/"+stationCode && strings.TrimSpace(id.Text) == station {
				found = true
			}
		}
		if !found {
			return nil, errors.New("wrong station")
		}
		matches = append(matches, series)
	}
	if len(matches) != 1 {
		return nil, fmt.Errorf("expected one %s series", property)
	}
	return matches[0], nil
}

func points(series *element) ([]waterlevel.Point, error) {
	nodes := series.all(wml, "MeasurementTVP")
	if len(nodes) == 0 {
		return nil, errors.New("no points")
	}
	var out []waterlevel.Point
	seen := map[time.Time]bool{}
	valid := 0
	for _, n := range nodes {
		tn, err := one(n, wml, "time")
		if err != nil {
			return nil, err
		}
		t, err := utc(strings.TrimSpace(tn.Text))
		if err != nil {
			return nil, err
		}
		if seen[t] {
			return nil, errors.New("duplicate time")
		}
		seen[t] = true
		vn, err := one(n, wml, "value")
		if err != nil {
			return nil, err
		}
		text := strings.TrimSpace(vn.Text)
		p := waterlevel.Point{Time: t}
		if text != "" && text != "NaN" {
			v, err := strconv.ParseFloat(text, 64)
			if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
				return nil, errors.New("invalid value")
			}
			p.ValueCM = &v
			valid++
		}
		out = append(out, p)
	}
	if valid == 0 {
		return nil, errors.New("no finite values")
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Time.Before(out[j].Time) })
	return out, nil
}

func instant(root, node *element) (time.Time, error) {
	if node.attr(xlink, "href") == "" {
		instants := node.all(gml, "TimeInstant")
		if len(instants) == 1 && instants[0] != node && instants[0].attr(xlink, "href") != "" {
			node = instants[0]
		}
	}
	if href := node.attr(xlink, "href"); href != "" {
		if !strings.HasPrefix(href, "#") || len(href) == 1 {
			return time.Time{}, errors.New("external time reference")
		}
		var matches []*element
		for _, n := range root.all(gml, "TimeInstant") {
			if n.attr(gml, "id") == href[1:] {
				matches = append(matches, n)
			}
		}
		if len(matches) != 1 {
			return time.Time{}, errors.New("unresolved time reference")
		}
		node = matches[0]
	}
	pos, err := one(node, gml, "timePosition")
	if err != nil {
		return time.Time{}, err
	}
	return utc(strings.TrimSpace(pos.Text))
}

// ParseObservation returns the latest finite observation, converted from mm to cm.
func ParseObservation(data []byte, now time.Time) (*waterlevel.Observation, error) {
	root, err := document(data)
	if err != nil {
		return nil, err
	}
	s, err := selectSeries(root, "WLEVN2K_PT1S_INSTANT", "134254", "fmisid")
	if err != nil {
		return nil, err
	}
	p, err := points(s)
	if err != nil {
		return nil, err
	}
	for i := len(p) - 1; i >= 0; i-- {
		if p[i].Time.After(now) {
			return nil, errors.New("future observation")
		}
		if p[i].ValueCM != nil {
			return &waterlevel.Observation{ValueCM: *p[i].ValueCM / 10, ObservedAt: p[i].Time, FetchedAt: now}, nil
		}
	}
	return nil, errors.New("no observation")
}

// ParseForecast accepts only the N2000 series and records its unverified cm assumption.
func ParseForecast(data []byte, now time.Time) (*waterlevel.Forecast, error) {
	root, err := document(data)
	if err != nil {
		return nil, err
	}
	s, err := selectSeries(root, "SeaLevelN2000", "-10022823", "geoid")
	if err != nil {
		return nil, err
	}
	for _, namespace := range []string{wml, swe} {
		for _, unit := range s.all(namespace, "uom") {
			code := unit.attr("", "code")
			if code != "" && code != "cm" {
				return nil, errors.New("conflicting forecast unit")
			}
		}
	}
	p, err := points(s)
	if err != nil {
		return nil, err
	}
	var analysis time.Time
	count := 0
	for _, n := range s.all(om, "NamedValue") {
		names := n.all(om, "name")
		if len(names) == 1 && strings.HasSuffix(names[0].attr(xlink, "href"), "/analysisTime") {
			value, err := one(n, om, "value")
			if err != nil {
				return nil, err
			}
			analysis, err = instant(root, value)
			if err != nil {
				return nil, err
			}
			count++
		}
	}
	if count != 1 || analysis.After(now) {
		return nil, errors.New("missing or future analysis time")
	}
	rt, err := one(s, om, "resultTime")
	if err != nil {
		return nil, err
	}
	result, err := instant(root, rt)
	if err != nil {
		return nil, err
	}
	return &waterlevel.Forecast{AnalysisTime: analysis, ResultTime: result, FetchedAt: now, UnitStatus: "assumed-cm", Points: p}, nil
}
