# Haminan lähdenäytteiden alkuperä

Lähde: **Ilmatieteen laitos, avoin data,
[CC BY 4.0](https://creativecommons.org/licenses/by/4.0/)**.
[FMI:n lisenssi](https://en.ilmatieteenlaitos.fi/open-data-licence).
XML-tiedostot ovat HTTP-vastausten alkuperäisiä runkoja muuttamattomina
(myös alkuperäiset sisennykset ja rivien lopun välilyönnit säilytetty).
Tulkinnat ovat [raportissa](../water-level.md), eivät lähdetiedostoissa.
Hakuajat ovat hakukoneen UTC-kellosta pyynnön alkaessa, eivät julkaisuaikoja.
Kaikki alla luetellut haut palauttivat **HTTP 200**.

## Dokumentaatiolähteet

Tarkistettu 2026-11-16:

- [WFS:n käyttö ja rajat](https://en.ilmatieteenlaitos.fi/open-data-manual-fmi-wfs-services)
- [Kyselyiden muodostaminen ja asematiedot](https://www.ilmatieteenlaitos.fi/latauspalvelun-pikaohje)
- [Asemametatietojen kuvaus](https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=DescribeStoredQueries&storedquery_id=fmi::ef::stations)
- [Havaintokyselyn kuvaus](https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=DescribeStoredQueries&storedquery_id=fmi::observations::mareograph::instant::timevaluepair)
- [Ennustekyselyn kuvaus](https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=DescribeStoredQueries&storedquery_id=fmi::forecast::sealevel::point::timevaluepair)
- [Ennusteen sisällön kuvaus](https://www.ilmatieteenlaitos.fi/avoin-data-meriveden-korkeusennuste)
- [Aineistoluettelo](https://www.ilmatieteenlaitos.fi/avoin-data-avattavat-aineistot)

## Tallennetut vastaukset

### `hamina-station.xml`

- [Tiedosto](hamina-station.xml)
- Hakuaika: `2026-11-16T18:06:24.327153Z`
- Dokumentaatio: asematietojen kuvaus ja latauspalvelun pikaohje yllä.
- Tulos: yksi `EnvironmentalMonitoringFacility`, FMISID 134254,
  nimi Hamina Pitäjänsaari, geoid -10022823, verkko 133.
- SHA-256: `04135bedc71ad3e741a193e9af917e1c234eabc8cfa31cd83270b3558f4f8a47`
- Suoritettu kysely:

```text
https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=GetFeature&storedquery_id=fmi%3A%3Aef%3A%3Astations&fmisid=134254
```

### `hamina-observation.xml`

- [Tiedosto](hamina-observation.xml)
- Hakuaika: `2026-11-16T18:08:05.475845Z`
- Dokumentaatio: havaintokyselyn kuvaus ja pikaohje yllä.
- Tulos: yksi `PointTimeSeriesObservation`, 61 minuuttiarvoa,
  FMISID 134254, geoid -10022823. Ei WFS-virhevastaus.
- SHA-256: `0f9bde6251e9c1c6127b2d9267d8027f4f4d7605a58130b4f70c3cacfe4f3121`
- Suoritettu kysely:

```text
https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=GetFeature&fmisid=134254&timezone=UTC&storedquery_id=fmi%3A%3Aobservations%3A%3Amareograph%3A%3Ainstant%3A%3Atimevaluepair&timestep=1&starttime=2026-11-16T17%3A00%3A00Z&endtime=2026-11-16T18%3A00%3A00Z&parameters=WLEVN2K_PT1S_INSTANT
```

### `hamina-forecast.xml`

- [Tiedosto](hamina-forecast.xml)
- Hakuaika: `2026-11-16T18:08:05.509802Z`
- Dokumentaatio: ennustekyselyn kuvaus, ennustekuvaus ja aineistoluettelo yllä.
- Tulos: kaksi `PointTimeSeriesObservation`-sarjaa, `SeaLevel` ja
  `SeaLevelN2000`, 121 aikaleimaa/sarja. Kummassakin 114 äärellistä
  arvoa ja lopussa seitsemän `NaN`-arvoa. Ei WFS-virhevastaus.
- Vastauksen asematunniste on **geoid -10022823**, ei FMISID. Se yhdistyy
  asemavastauksen kautta FMISID:iin 134254; nimi on Hamina Pitäjänsaari
  ja koordinaatit ovat samat näytön tarkkuuden rajoissa.
- SHA-256: `5005283606e70e5377b94c37ab35d92247f87a40ff72490a5b33a5b51e65ecfd`
- Suoritettu kysely:

```text
https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=GetFeature&fmisid=134254&timezone=UTC&storedquery_id=fmi%3A%3Aforecast%3A%3Asealevel%3A%3Apoint%3A%3Atimevaluepair&timestep=60&starttime=2026-11-16T18%3A00%3A00Z&endtime=2026-11-21T18%3A00%3A00Z
```

### `observation-property.xml`

- [Tiedosto](observation-property.xml)
- Hakuaika: `2026-11-16T18:08:25.629078Z`
- Dokumentaatio: havaintovastauksen `om:observedProperty`-linkki,
  seurattu täsmälleen sellaisenaan, vaikka URL:n luokka on `forecast`.
- Tulos: `Water level (N2000)`, `N2000 quality controlled water level`,
  yksikkö `mm`.
- SHA-256: `cb8acc67453d4aea26e59b64d7a0a517a043f717202584590c514f2951436183`
- Suoritettu kysely:

```text
https://opendata.fmi.fi/meta?observableProperty=forecast&param=WLEVN2K_PT1S_INSTANT&language=eng
```

### `forecast-property.xml`

- [Tiedosto](forecast-property.xml)
- Hakuaika: `2026-11-16T18:08:25.653852Z`
- Dokumentaatio: ennustevastauksen `om:observedProperty`-linkki.
- Tulos: **tyhjä `ObservableProperty`**, ei yksikköä tai suureen kuvausta.
  HTTP-onnistuminen ei tarkoita metatiedon sisällöllistä onnistumista.
  Tämä on tallennettu näyttö dokumentaatiopuutteesta, ei osoitus ennusteen
  puuttumisesta.
- SHA-256: `5ef20b26a9c1c3e36cd2bfe0805c55bc24a9b90d486c65fd2c5576a9fda7b76a`
- Suoritettu kysely:

```text
https://opendata.fmi.fi/meta?observableProperty=forecast&param=SeaLevelN2000&language=eng
```

## Rajattu löytöhaku

Havaintosuureen nimi selvitettiin pikaohjeen mukaisella yhden ajanhetken
haulla, eikä erillistä löytövastausta säilytetä tarpeettomasti. Nimi ja
metatietolinkki ovat yllä säilytetyssä varsinaisessa havaintovastauksessa.
Tämä haku ei ole uusi erillinen havaintonäyte:

```text
https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=GetFeature&storedquery_id=fmi%3A%3Aobservations%3A%3Amareograph%3A%3Ainstant%3A%3Atimevaluepair&fmisid=134254&timezone=UTC&timestep=1&starttime=2026-11-16T18%3A00%3A00Z&endtime=2026-11-16T18%3A00%3A00Z
```
