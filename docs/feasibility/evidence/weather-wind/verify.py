"""Offline checks for the fixed feasibility evidence, not application code.

Run from any directory with Python 3.9+; makes no network requests.
"""

import datetime as dt
import hashlib
import math
from pathlib import Path
import re
from urllib.parse import parse_qs, urlsplit
import xml.etree.ElementTree as ET


BASE = Path(__file__).resolve().parent
NS = {
    "wfs": "http://www.opengis.net/wfs/2.0",
    "gml": "http://www.opengis.net/gml/3.2",
    "cov": "http://www.opengis.net/gmlcov/1.0",
    "swe": "http://www.opengis.net/swe/2.0",
    "om": "http://www.opengis.net/om/2.0",
    "target": "http://xml.fmi.fi/namespace/om/atmosphericfeatures/1.1",
}
EXPECTED = {
    "cottage.xml": {
        "sha": "c8ad1c887e57e4f13c9ca2b13592bed34d03b59fec9063e8b69f7937725bc68a",
        "point": (60.43432, 26.68783),
        "fields": ["Temperature", "WindSpeedMS", "WindGust", "WindDirection",
                   "WeatherSymbol3", "Precipitation1h"],
        "start": "2026-11-17T18:00:00Z", "end": "2026-11-19T18:00:00Z",
        "step": 3600, "rows": 49,
        "retrieved": "2026-11-17T18:06:12.261681Z",
        "result": "2026-11-17T15:04:04Z", "analysis": "2026-11-17T12:00:00Z",
    },
    "offshore.xml": {
        "sha": "3b84a82f393598eb3effa0f372b0a1a54c64e6422cb7ae99b05698a7f99429ce",
        "point": (60.41462, 26.69740),
        "fields": ["WindSpeedMS", "WindGust", "WindDirection"],
        "start": "2026-11-17T18:00:00Z", "end": "2026-11-19T18:00:00Z",
        "step": 3600, "rows": 49,
        "retrieved": "2026-11-17T18:08:55.423601Z",
        "result": "2026-11-17T18:03:58Z", "analysis": "2026-11-17T15:00:00Z",
    },
    "rankki.xml": {
        "sha": "c69b5a94e0ad37915792cd2ae97e56a9a82188860ed9c77fb6813b86288150d7",
        "point": (60.37538, 26.95893), "crs": "4258", "fmisid": "101030",
        "fields": ["ws_10min", "wg_10min", "wd_10min"],
        "start": "2026-11-17T17:00:00Z", "end": "2026-11-17T18:00:00Z",
        "step": 600, "rows": 7,
        "retrieved": "2026-11-17T18:09:56.139618Z",
        "result": "2026-11-17T18:00:00Z",
    },
}


def instant(value):
    return dt.datetime.fromisoformat(value.replace("Z", "+00:00"))


def utc(epoch):
    return dt.datetime.fromtimestamp(epoch, dt.timezone.utc).isoformat()


def check_sample(filename, expected):
    data = (BASE / filename).read_bytes()
    digest = hashlib.sha256(data).hexdigest()
    assert digest == expected["sha"], (filename, "SHA-256 mismatch")
    assert digest in (BASE / "manifest.md").read_text(), "Missing manifest digest"
    root = ET.fromstring(data)
    assert root.tag == f"{{{NS['wfs']}}}FeatureCollection"
    assert root.attrib["numberReturned"] == "1"
    assert len(root.findall("wfs:member", NS)) == 1
    fields = [node.attrib["name"] for node in root.findall(".//swe:field", NS)]
    assert fields == expected["fields"], fields
    domain = root.find(".//cov:SimpleMultiPoint", NS)
    assert domain.attrib["srsDimension"] == "3"
    crs = expected.get("crs", "4326")
    assert domain.attrib["srsName"].endswith(f"crs={crs}&time=unixtime")
    positions = [list(map(float, row.split())) for row in
                 domain.find("cov:positions", NS).text.strip().splitlines()]
    rows = [list(map(float, row.split())) for row in root.find(
        ".//gml:doubleOrNilReasonTupleList", NS).text.strip().splitlines()]
    assert len(positions) == len(rows) == expected["rows"]
    assert all(len(row) == 3 for row in positions)
    assert all(len(row) == len(fields) for row in rows)
    assert all(tuple(row[:2]) == expected["point"] for row in positions)
    times = [row[2] for row in positions]
    assert times[0] == instant(expected["start"]).timestamp()
    assert times[-1] == instant(expected["end"]).timestamp()
    assert all(b - a == expected["step"] for a, b in zip(times, times[1:]))
    assert root.find(".//gml:beginPosition", NS).text == expected["start"]
    assert root.find(".//gml:endPosition", NS).text == expected["end"]
    result = root.find(".//om:resultTime//gml:timePosition", NS).text
    assert result == expected["result"]
    if "analysis" in expected:
        analysis = root.find(".//om:parameter//gml:timePosition", NS).text
        assert analysis == expected["analysis"]
    if "fmisid" in expected:
        station = root.find(".//target:Location/gml:identifier", NS)
        assert station.text == expected["fmisid"]
        assert station.attrib["codeSpace"].endswith("/fmisid")
    print(filename, "SHA/XML/fields/location/times OK;",
          len(rows), "rows; step", expected["step"], "seconds; result", result)
    for index, field in enumerate(fields):
        finite = [(time, row[index]) for time, row in zip(times, rows)
                  if math.isfinite(row[index])]
        assert len(finite) == len(rows), (filename, field, "non-finite value")
        values = [value for _, value in finite]
        print(" ", field, "finite", len(finite), "missing", len(rows)-len(finite),
              "range", min(values), max(values), "first", utc(finite[0][0]),
              "last", utc(finite[-1][0]), "latest value", finite[-1][1])
        if "analysis" in expected:
            horizon = instant(expected["end"]) - instant(expected["retrieved"])
            print("   finite future horizon from retrieval:", horizon)
    if "WeatherSymbol3" in fields:
        codes = sorted({row[fields.index("WeatherSymbol3")] for row in rows})
        assert codes == [3.0, 31.0, 32.0]
        print("  WeatherSymbol3 codes", codes)


def check_links():
    for path in (BASE / "manifest.md", BASE.parent.parent / "weather-wind.md"):
        for target in re.findall(r"\[[^\]]*\]\(([^)]+)\)", path.read_text()):
            if "://" not in target and not target.startswith("#"):
                assert (path.parent / target.split("#")[0]).exists(), (path, target)
    print("Report and manifest local links OK")


def check_requests():
    urls = re.findall(r"^https://opendata\.fmi\.fi/wfs\?[^\n]+",
                      (BASE / "manifest.md").read_text(), re.MULTILINE)
    assert len(urls) == len(EXPECTED) == 3
    assert {path.name for path in BASE.glob("*.xml")} == set(EXPECTED)
    for url, (filename, expected) in zip(urls, EXPECTED.items()):
        query = parse_qs(urlsplit(url).query)
        assert query["service"] == ["WFS"] and query["version"] == ["2.0.0"]
        assert query["request"] == ["GetFeature"]
        assert query["parameters"] == [",".join(expected["fields"])]
        assert query["starttime"] == [expected["start"]]
        assert query["endtime"] == [expected["end"]]
        assert query["timestep"] == [str(expected["step"] // 60)]
        window = instant(expected["end"]) - instant(expected["start"])
        if "fmisid" in expected:
            assert query["fmisid"] == [expected["fmisid"]]
            assert window <= dt.timedelta(hours=1)
        else:
            requested = tuple(map(float, query["latlon"][0].split(",")))
            assert all(abs(a - b) < 0.00001
                       for a, b in zip(requested, expected["point"]))
            assert window <= dt.timedelta(hours=48)
        producer = ("fmi::observations::weather::multipointcoverage"
                    if "fmisid" in expected else
                    "fmi::forecast::harmonie::surface::point::multipointcoverage")
        assert query["storedquery_id"] == [producer]
        print(filename, "manifest request/point/window OK")


if __name__ == "__main__":
    check_requests()
    for name, expectation in EXPECTED.items():
        check_sample(name, expectation)
    check_links()
