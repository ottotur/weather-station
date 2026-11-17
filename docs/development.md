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

Sovellus on keskeneräinen eikä sen koko ketjua ole vielä hyväksyntätestattu.

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
`WEATHER_CACHE_DIR` välimuistihakemiston. `WEATHER_SOURCE_URL` vaihtaa
FMI-lähteen osoitteen esimerkiksi paikallista testipalvelinta varten.
Oletuksilla käynnistys aloittaa todelliset FMI-haut. Kello voidaan antaa
laskenta- ja palvelukerrokselle erikseen testeissä. Koko ketjun testissä
`WEATHER_TEST_NOW` lukitsee kellon RFC 3339 -aikaan; sitä ei tule käyttää
normaalissa ajossa, koska jäädytetty kello estää tuoreuden oikean arvioinnin.

Vedenkorkeuden laskentamallissa havainnon tuoreus määräytyy mittausajasta
(30 minuuttia) ja ennusteen tuoreus `analysisTime`-mallin viiteajasta
(12 tuntia). Täsmälleen rajalla aineisto on vielä tuore. Mallin viiteaika ei ole
varmennettu julkaisuaika, eikä uusi haku nuorenna vanhaa malliajoa. Ennusteen
tuleva kattavuus tarkistetaan tuoreudesta erikseen.

## Välimuisti ja tallennusvirheet

Välimuistihakemistossa ovat erilliset `observation.json` ja `forecast.json`.
Niiden JSON-kuoressa on `version: 1`, lähde `FMI`, aseman tunniste ja suure
sekä `data`, joka sisältää normalisoidut arvot ja ajat; ennusteessa on lisäksi
`unitStatus: assumed-cm`. Kehityksen aiemman keskeneräisen version tiedostot,
joista nämä lähdetiedot puuttuvat, hylätään latauksessa eikä niitä pidetä
varmennettuna aineistona. Käynnistyksessä
tiedostot tarkistetaan erikseen. Vioittunut tai tuntemattoman version tiedosto
ei estä toisen aineiston käyttöä, ja väliaikaistiedostoja ei ladata.
Vanha havaintotiedosto, jossa on vaaditut lähdetiedot mutta ei tulosaikaa,
voidaan edelleen ladata: sen `resultTime` on API:ssa `null`, kunnes uusi
tarkistettu havainto tallennetaan. Tulosaikaa ei päätellä muista ajoista.

Ennen atomista korvausta tapahtuva virhe säilyttää aiemman tiedoston ja
muistiaineiston. Jos korvaus onnistuu mutta hakemiston synkronointi epäonnistuu,
API:n yritysvirhe on `cache_durability_uncertain`: muistissa on vanha aineisto,
mutta uudelleenkäynnistyksessä voidaan lukea uusi tiedosto. Tarkista levyn tila
ja oikeudet; myöhempi onnistunut haku tallentaa aineiston uudelleen.

Toteutuksen aiemmassa tarkistuksessa `npm audit` raportoi kehitysaikaisesta
`deepmerge`-riippuvuudesta huomautuksen GHSA-gcqx-w48m-w2h4. Uusimmassa
2026-11-17 suoritetussa `npm audit` -ajossa ei raportoitu haavoittuvuuksia.

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
Jos vastaus ilmoittaa yksikkötiedon, valitun sarjan `wml2:uom`- ja
`swe:uom`-, `omop:uom`- sekä `unitOfMeasure`-elementtien `code`-, `uom`- ja
`xlink:href`-tiedot tarkistetaan. Hyväksytyt koodit ovat havainnolle `mm`
ja ennusteelle `cm`; hyväksytyt viitteet ovat täsmälliset QUDT-tunnisteet
[`unit:MilliM`](https://qudt.org/vocab/unit/MilliM) ja
[`unit:CentiM`](https://qudt.org/vocab/unit/CentiM). Viitteitä ei haeta
ajonaikaisesti. Muu tai ristiriitainen ilmoitettu tieto hylätään; puuttuva
metatieto hyväksytään nykyisellä tulkinnalla eikä ennusteen oletus muutu
varmennetuksi tiedoksi.
Myös WaterML-pistearvon `uom` tarkistetaan. Tulkinnan rajaus, skeemalähteet
ja testimatriisi ovat [XML-selvityksessä](feasibility/water-level-xml-schema.md).
FMI:n selvityksessä tarkistetut käyttörajat ja niiden lähteet:
[vedenkorkeusselvitys](feasibility/water-level.md#dokumentoidut-kyselyt-ja-käyttöehdot).
Hakuvälit ovat 10 ja 30 minuuttia. Toteutuksen koko ketjun testaus on kesken.
Rajapinnan kentät ja esimerkit: [api.md](api.md).

## Nykyiset automaattiset tarkistukset

Projektin juuresta: `go test ./...`, `go test -race ./...` ja `go vet ./...`.
Hakemistossa `web`: `npm run check` ja `npm run build`. Nämä eivät korvaa
vielä puuttuvia koko ketjun selain- tai oikean laitteen testejä.
Selaintestit: hakemistossa `web` ensin `npx playwright install chromium` ja
sen jälkeen `npm test`. Chromiumin kosketusemulaatio ei ole oikean laitteen
testi. Testatut selainikkunoiden leveydet ovat 360, 768 ja 1024 CSS-pikseliä.
Testeissä vaihdetaan myös korkeutta pysty- ja vaakasuhteiden välillä ja
tarkistetaan kaavion leveys, yhteenvetojen ja tuntitaulukon näkyminen sekä
sivutason vaakavierityksen puuttuminen.
Koko ketjun testin toistettavat komennot projektin juuresta:

```sh
cd web
npm ci
npx playwright install chromium
npm run build
WEATHER_E2E_TMPDIR=/tmp/opencode npm test
```

Testi kääntää oikean Go-ohjelman, käynnistää paikallisen WFS-valepalvelimen
ja käyttää alkuperäisiä näytevastauksia kiinteällä testikellolla. Se käy läpi
normaalitilan, virheellisen vastauksen, lähdekatkon, vanhenemisen,
uudelleenkäynnistykset ja palautumisen ilman live-FMI-riippuvuutta.
Rakennetun käyttöliittymän Go-palvelintesti ajetaan projektin juuresta
`WEATHER_TEST_WEB_DIST=../../web/dist go test -race ./cmd/weather-station`
sen jälkeen, kun `npm run build` on tehty hakemistossa `web`. Polku tulkitaan
testipaketin `cmd/weather-station` työhakemistosta.
