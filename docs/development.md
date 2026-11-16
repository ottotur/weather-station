# Kehitysympäristö

Käytä projektia WSL:n Linux-tiedostojärjestelmässä, ei Windowsin `/mnt/c`-polulla.
Toteutuksen alussa todettu ympäristö: Windows 11, WSL 2.7.14.0, Ubuntu
26.04.1 LTS, Linux x86_64; `wsl.exe --version`, `cat /etc/os-release` ja
`uname -m` näyttävät paikalliset tiedot.

Työkalut: Go 1.27.4 (Go-jakelun uusin vakaa versio tarkistushetkellä;
`go version`, tarvittaessa `GOTOOLCHAIN=go1.27.4 go version`), Node 24.21.0
ja npm 11.19.0 (`node --version`, `npm --version`). Käyttöliittymän lukitut
versiot ovat Svelte 5.57.3, Vite 8.3.8, Svelte Vite -liitännäinen 7.3.1,
TypeScript 5.9.3, Node-tyypit 24.19.2 ja svelte-check 4.7.8; versiot voi tarkistaa
komennolla `npm ls --depth=0` hakemistossa `web`.

Tämä on vasta keskeneräinen runko, ei toimiva vedenkorkeussovellus.

Rakennetun rungon voi käynnistää projektin juuresta:

```sh
cd web
npm ci
npm run build
cd ..
go run ./cmd/weather-station
```

Oletusosoite on `http://127.0.0.1:8080`. `WEATHER_LISTEN` vaihtaa
kuunteluosoitteen, `WEATHER_STATIC_DIR` staattisen hakemiston ja
`WEATHER_CACHE_DIR` välimuistihakemiston. `WEATHER_SOURCE_URL` on varattu
FMI-lähteen osoitteelle testejä varten; lähdehakua ei ole vielä toteutettu.
Kello annetaan laskentafunktioille erikseen, ja käynnistysasetukset sisältävät
testattavan kellofunktion. Palvelin ei vielä käytä kumpaakaan hakuihin.

Vedenkorkeuden laskentamallissa havainnon tuoreus määräytyy mittausajasta
(30 minuuttia) ja ennusteen tuoreus `analysisTime`-mallin viiteajasta
(12 tuntia). Täsmälleen rajalla aineisto on vielä tuore. Mallin viiteaika ei ole
varmennettu julkaisuaika, eikä uusi haku nuorenna vanhaa malliajoa. Ennusteen
tuleva kattavuus tarkistetaan tuoreudesta erikseen.

`npm audit` raportoi kehitysaikaisen Svelte Vite -liitännäisen `deepmerge`-
riippuvuudesta matalan vakavuuden haavoittuvuuden GHSA-gcqx-w48m-w2h4.
Tarkistushetkellä npm ei tarjonnut korjausta; riippuvuus ei ole selaimen
ajonaikainen riippuvuus.

## FMI-lähde

Hakulogiikka käyttää [FMI:n WFS-latauspalvelua](https://www.ilmatieteenlaitos.fi/latauspalvelun-pikaohje)
ja yhden aseman dokumentoituja kyselyitä:
[mareografihavainto](https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=DescribeStoredQueries&storedquery_id=fmi::observations::mareograph::instant::timevaluepair)
ja [vedenkorkeusennuste](https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=DescribeStoredQueries&storedquery_id=fmi::forecast::sealevel::point::timevaluepair).
Havainto rajataan viimeiseen tuntiin ja ennuste korkeintaan seuraavaan 48 tuntiin.
Vastaus rajataan 2 MiB:iin ja HTTP-aikakatkaisu on 30 sekuntia.

Lähde: Ilmatieteen laitos, avoin data, [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/).
Havainnon millimetrit muunnetaan senttimetreiksi. `SeaLevelN2000`-ennusteen cm
on käyttäjän hyväksymä oletus, ei suoraan varmennettu yksikkömetatieto.
FMI:n selvityksessä tarkistetut käyttörajat ja niiden lähteet:
[vedenkorkeusselvitys](feasibility/water-level.md#dokumentoidut-kyselyt-ja-käyttöehdot).
Suunnitellut hakuvälit ovat 10 ja 30 minuuttia, mutta ajastusta ei vielä ole.
