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
	omop  = "http://inspire.ec.europa.eu/schemas/omop/2.9"
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
	value, _ := e.attribute(space, local)
	return value
}

func (e *element) attribute(space, local string) (string, bool) {
	for _, a := range e.Attrs {
		if a.Name.Space == space && a.Name.Local == local {
			return a.Value, true
		}
	}
	return "", false
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

func utcPosition(node *element) (time.Time, error) {
	for _, attr := range node.Attrs {
		if attr.Name.Space == xlink && attr.Name.Local == "href" {
			return time.Time{}, errors.New("unsupported reference on time position")
		}
		if attr.Name.Space != "" {
			continue
		}
		switch attr.Name.Local {
		case "indeterminatePosition":
			return time.Time{}, errors.New("indeterminate time position")
		case "frame":
			if attr.Value != "#ISO-8601" {
				return time.Time{}, errors.New("unsupported time reference system")
			}
		}
	}
	return utc(strings.TrimSpace(node.Text))
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

// rawPoint retains the source value without asserting its unit.
type rawPoint struct {
	Time  time.Time
	Value *float64
}

func points(series *element) ([]rawPoint, error) {
	nodes := series.all(wml, "MeasurementTVP")
	if len(nodes) == 0 {
		return nil, errors.New("no points")
	}
	var out []rawPoint
	seen := map[time.Time]bool{}
	valid := 0
	for _, n := range nodes {
		tn, err := one(n, wml, "time")
		if err != nil {
			return nil, err
		}
		t, err := utcPosition(tn)
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
		p := rawPoint{Time: t}
		if text != "" && text != "NaN" {
			v, err := strconv.ParseFloat(text, 64)
			if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
				return nil, errors.New("invalid value")
			}
			p.Value = &v
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
	if href, present := node.attribute(xlink, "href"); present {
		if node.Name.Space == gml && node.Name.Local == "TimeInstant" {
			return time.Time{}, errors.New("reference on TimeInstant")
		}
		if len(node.Children) != 0 || strings.TrimSpace(node.Text) != "" {
			return time.Time{}, errors.New("time reference also contains a value")
		}
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
	instants := node.all(gml, "TimeInstant")
	if len(instants) != 1 {
		return time.Time{}, errors.New("invalid time instant")
	}
	node = instants[0]
	if _, present := node.attribute(xlink, "href"); present {
		return time.Time{}, errors.New("reference on TimeInstant")
	}
	pos, err := one(node, gml, "timePosition")
	if err != nil {
		return time.Time{}, err
	}
	return utcPosition(pos)
}

func periodBoundary(root, period *element, positionName, propertyName string) (time.Time, error) {
	var boundary *element
	for _, child := range period.Children {
		if child.Name.Space == gml && (child.Name.Local == positionName || child.Name.Local == propertyName) {
			if boundary != nil {
				return time.Time{}, errors.New("ambiguous time period boundary")
			}
			boundary = child
		}
	}
	if boundary == nil {
		return time.Time{}, errors.New("missing time period boundary")
	}
	if boundary.Name.Local == propertyName {
		return instant(root, boundary)
	}
	if _, present := boundary.attribute(xlink, "href"); present || len(boundary.Children) != 0 {
		return time.Time{}, errors.New("invalid time position")
	}
	return utcPosition(boundary)
}

func validatePhenomenonTime(root, series *element) error {
	node, err := one(series, om, "phenomenonTime")
	if err != nil {
		return err
	}
	if href, present := node.attribute(xlink, "href"); present {
		if len(node.Children) != 0 || strings.TrimSpace(node.Text) != "" {
			return errors.New("phenomenon time reference also contains a value")
		}
		if !strings.HasPrefix(href, "#") || len(href) == 1 {
			return errors.New("external phenomenon time reference")
		}
		var targets []*element
		for _, name := range []string{"TimePeriod", "TimeInstant"} {
			for _, candidate := range root.all(gml, name) {
				if candidate.attr(gml, "id") == href[1:] {
					targets = append(targets, candidate)
				}
			}
		}
		if len(targets) != 1 {
			return errors.New("unresolved phenomenon time reference")
		}
		node = targets[0]
	}
	var objects []*element
	if node.Name.Space == gml && (node.Name.Local == "TimePeriod" || node.Name.Local == "TimeInstant") {
		objects = append(objects, node)
	} else {
		for _, child := range node.Children {
			if child.Name.Space == gml && (child.Name.Local == "TimePeriod" || child.Name.Local == "TimeInstant") {
				objects = append(objects, child)
			}
		}
	}
	if len(objects) != 1 {
		return errors.New("invalid phenomenon time")
	}
	object := objects[0]
	if object.Name.Local == "TimePeriod" {
		if _, present := object.attribute(xlink, "href"); present {
			return errors.New("unsupported reference on TimePeriod")
		}
		startTime, err := periodBoundary(root, object, "beginPosition", "begin")
		if err != nil {
			return err
		}
		endTime, err := periodBoundary(root, object, "endPosition", "end")
		if err != nil {
			return err
		}
		if endTime.Before(startTime) {
			return errors.New("invalid phenomenon time period")
		}
		return nil
	}
	_, err = instant(root, object)
	return err
}

// validateUnit checks only metadata belonging to the selected series. Missing
// metadata is allowed; declared units must be explicit and known.
func validateUnit(series *element, expected string) error {
	identifiers := map[string]string{
		"http://qudt.org/vocab/unit/MilliM": "mm",
		"http://qudt.org/vocab/unit/CentiM": "cm",
	}
	var units []*element
	var collect func(*element)
	collect = func(node *element) {
		if node.Name.Local == "uom" || node.Name.Local == "unitOfMeasure" {
			units = append(units, node)
		}
		for _, child := range node.Children {
			collect(child)
		}
	}
	collect(series)
	for _, value := range series.all(wml, "value") {
		for _, attr := range value.Attrs {
			if attr.Name.Space == "" && attr.Name.Local == "uom" && strings.TrimSpace(attr.Value) != expected {
				return fmt.Errorf("%w: contradictory or unknown value unit", waterlevel.ErrUnitMetadata)
			}
		}
	}
	for _, unit := range units {
		if unit.Name.Space != wml && unit.Name.Space != swe && unit.Name.Space != omop {
			return fmt.Errorf("%w: unknown unit namespace", waterlevel.ErrUnitMetadata)
		}
		declared := false
		for _, attr := range unit.Attrs {
			value := strings.TrimSpace(attr.Value)
			switch {
			case attr.Name.Space == "" && (attr.Name.Local == "code" || attr.Name.Local == "uom"):
				declared = true
				if value != expected {
					return fmt.Errorf("%w: contradictory or unknown unit code", waterlevel.ErrUnitMetadata)
				}
			case attr.Name.Space == xlink && attr.Name.Local == "href":
				declared = true
				if identifiers[value] != expected {
					return fmt.Errorf("%w: contradictory or unknown unit reference", waterlevel.ErrUnitMetadata)
				}
			}
		}
		if !declared || strings.TrimSpace(unit.Text) != "" || len(unit.Children) > 0 {
			return fmt.Errorf("%w: unsupported unit declaration", waterlevel.ErrUnitMetadata)
		}
	}
	return nil
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
	if err := validatePhenomenonTime(root, s); err != nil {
		return nil, err
	}
	if err := validateUnit(s, "mm"); err != nil {
		return nil, err
	}
	resultTime, err := one(s, om, "resultTime")
	if err != nil {
		return nil, err
	}
	result, err := instant(root, resultTime)
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
		if p[i].Value != nil {
			return &waterlevel.Observation{ValueCM: *p[i].Value / 10, ObservedAt: p[i].Time, ResultTime: &result, FetchedAt: now}, nil
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
	if err := validatePhenomenonTime(root, s); err != nil {
		return nil, err
	}
	if err := validateUnit(s, "cm"); err != nil {
		return nil, err
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
	normalized := make([]waterlevel.Point, len(p))
	for i, point := range p {
		normalized[i] = waterlevel.Point{Time: point.Time, ValueCM: point.Value}
	}
	return &waterlevel.Forecast{AnalysisTime: analysis, ResultTime: result, FetchedAt: now, UnitStatus: "assumed-cm", Points: normalized}, nil
}
