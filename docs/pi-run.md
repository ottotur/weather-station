# Manuaalisen Pi-ajon ohje

Tämä ohje ei ole laitehyväksyntä. Kohde on käyttäjän hyväksymä Raspberry Pi 3
Model B+ ja Raspberry Pi OS Lite 64-bit / Debian 13 Trixie. Pi:n ajo edellyttää,
että laite on valmis ja käyttäjä on erikseen sallinut ajon. Ohje ei alusta,
asenna tai muuta laitteen käyttöjärjestelmää, verkkoa tai palomuuria.

## Rakenna kehityskoneessa

Projektin juuresta, saman Git-version lähteistä:

```sh
cd web
npm ci
npm run build
cd ..
GOOS=linux GOARCH=arm64 GOARM64=v8.0 CGO_ENABLED=0 \
  go build -o /tmp/opencode/weather-station-linux-arm64 ./cmd/weather-station
file /tmp/opencode/weather-station-linux-arm64
```

`file`-tuloksen tulee ilmoittaa 64-bittinen ARM aarch64 -Linux-ohjelma.
Ristiinkäännös ei todista laiteajoa. Siirrä käyttäjän erikseen valitsemalla
menetelmällä ohjelma ja **saman version** `web/dist` laitteelle; älä sekoita
eri versioita.

## Suorita vain erillisellä luvalla

Alla oletetaan, että käyttäjä on luonut versiohakemiston `~/weather-station/v1`,
kopioinut ohjelman sinne nimellä `weather-station` ja saman version rakennetun
käyttöliittymän hakemistoon `~/weather-station/v1/web/dist`. Nämä esimerkkipolut
eivät ole järjestelmäasennus. Tarkista laitteella:

```sh
uname -m
cat /etc/os-release
test -f "$HOME/weather-station/v1/web/dist/index.html"
test -x "$HOME/weather-station/v1/weather-station"
mkdir -p "$HOME/weather-station/cache"
```

Arkkitehtuurin tulee olla `aarch64` ja käyttöjärjestelmän vastata sovittua
Trixie-kohdetta. Pysähdy, jos ne eivät vastaa; älä asenna tai muuta laitetta
ohjeen perusteella. Käynnistä vain käyttäjän sallimassa ympäristössä:

```sh
WEATHER_LISTEN=127.0.0.1:8080 \
WEATHER_STATIC_DIR="$HOME/weather-station/v1/web/dist" \
WEATHER_CACHE_DIR="$HOME/weather-station/cache" \
"$HOME/weather-station/v1/weather-station"
```

Prosessi pysäytetään Ctrl+C:llä. Oletuskuuntelu on vain laitteen localhostissa.
Lähiverkkotesti tarvitsee käyttäjän erikseen valitseman ja hyväksymän
kuunteluosoitteen `WEATHER_LISTEN`-asetukseen; mitään verkkoa ei avata
automaattisesti. Käynnistys aloittaa live-FMI-haut, joten käyttöehdot ja
hakurajat on tarkistettava ensin [kehitysohjeesta](development.md#fmi-lähde).
Välimuisti on erillään version tiedostoista, joten pysäyttäminen ei poista sitä.

Varsinaisessa laitehyväksynnässä kirjataan laite ja OS, käynnistys, rajatut
live-haut, internetkatko, välimuistin säilyminen ja uudelleenkäynnistys.
Mobiili- ja iPad-testit sekä suorituskykymittaukset kirjataan erikseen.
Niitä ei ole vielä tehty eikä WSL:n selainemulaatio korvaa niitä.
