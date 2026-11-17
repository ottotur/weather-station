# Sää- ja tuuliselvityksen lähdenäytteet

Lähde: Ilmatieteen laitoksen avoin data — HARMONIE (MEPS) -pistesääennuste ja
säähavainnot, [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/).
Ehdot ja kyselyparametrit tarkistettu 2026-11-17 ennen näytteitä.
Tiedostot ovat HTTP-vastausten muuttamattomia tavuja.
Hakuaika tarkoittaa pyynnön aloitusta UTC:ssa, ei lähteen julkaisuaikaa.
Metatietosivujen tarkistukset ja yksi epäonnistunut asemametatietoyritys on
kirjattu [raporttiin](../../weather-wind.md); ne eivät ole sääaineistonäytteitä.

## Testimökkipiste: sääennuste

- Tarkoitus: kuuden ydinsuureen saatavuus ja enintään 48 h kattavuus.
- Tiedosto: [cottage.xml](cottage.xml)
- UTC-hakuaika: `2026-11-17T18:06:12.261681Z`
- HTTP: **200**
- Koko: 11867 tavua
- SHA-256: `c8ad1c887e57e4f13c9ca2b13592bed34d03b59fec9063e8b69f7937725bc68a`
- Täsmällinen URL:

```text
https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=GetFeature&storedquery_id=fmi%3A%3Aforecast%3A%3Aharmonie%3A%3Asurface%3A%3Apoint%3A%3Amultipointcoverage&latlon=60.434318%2C26.687825&starttime=2026-11-17T18%3A00%3A00Z&endtime=2026-11-19T18%3A00%3A00Z&timestep=60&parameters=Temperature%2CWindSpeedMS%2CWindGust%2CWindDirection%2CWeatherSymbol3%2CPrecipitation1h
```

## Testiulappapiste: tuuliennuste

- Tarkoitus: hyväksytyn testipisteen kolme tuulisuuretta, enintään 48 h.
- Tiedosto: [offshore.xml](offshore.xml)
- UTC-hakuaika: `2026-11-17T18:08:55.423601Z`
- HTTP: **200**
- Koko: 10682 tavua
- SHA-256: `3b84a82f393598eb3effa0f372b0a1a54c64e6422cb7ae99b05698a7f99429ce`
- Täsmällinen URL:

```text
https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=GetFeature&storedquery_id=fmi%3A%3Aforecast%3A%3Aharmonie%3A%3Asurface%3A%3Apoint%3A%3Amultipointcoverage&latlon=60.414616%2C26.697395&starttime=2026-11-17T18%3A00%3A00Z&endtime=2026-11-19T18%3A00%3A00Z&timestep=60&parameters=WindSpeedMS%2CWindGust%2CWindDirection
```

## Hyväksytty Kotka Rankki: tuulihavainnot

- Tarkoitus: kolme tuulisuuretta, yksi viimeisin päättynyt tasatuntiväli.
- Hyväksyntä: käyttäjä valitsi Rankin toteutuskeskustelussa 2026-11-17 ennen hakua.
- Tiedosto: [rankki.xml](rankki.xml)
- UTC-hakuaika: `2026-11-17T18:09:56.139618Z`
- HTTP: **200**
- Koko: 6852 tavua
- SHA-256: `c69b5a94e0ad37915792cd2ae97e56a9a82188860ed9c77fb6813b86288150d7`
- Täsmällinen URL:

```text
https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=GetFeature&storedquery_id=fmi%3A%3Aobservations%3A%3Aweather%3A%3Amultipointcoverage&fmisid=101030&timestep=10&parameters=ws_10min%2Cwg_10min%2Cwd_10min&starttime=2026-11-17T17%3A00%3A00Z&endtime=2026-11-17T18%3A00%3A00Z
```

Kolme sääaineistohakua, jokainen HTTP 200; ei epäonnistuneita
sääaineistohakuja eikä uusintahakuja. Kuvaukset, kenttäsanastot ja asemasivut
ovat raportissa URL:ineen; niiden sisältö ei ole näissä kolmessa XML:ssä.
Offline-tarkistukset koskevat tallennettuja näytteitä, eivät verkkodokumenttien
tulevaa muuttumattomuutta. Toistaminen on ensisijaisesti offline:
[verify.py](verify.py). Historiallisen ennuste-URL:n uudelleenhaku ei takaa
samaa malliajoa, koska rajapinta tarjoaa viimeisintä ajoa.
