# Sääasema

Sääasema hakee Haminan vedenkorkeushavainnon ja ennusteen sekä näyttää ne
paikallisessa selainnäkymässä. Toteutus on keskeneräinen: koko ketjun
hyväksyntätestit puuttuvat, eikä sovellusta ole testattu Raspberry Pi:llä
tai oikeilla mobiililaitteilla. Ennusteen cm on käyttäjän hyväksymä oletus,
ei suoraan varmennettu metatieto.

## Paikallinen ajo

Ubuntu WSL:n Linux-tiedostojärjestelmässä:

```sh
cd web
npm ci
npm run build
cd ..
go run ./cmd/weather-station
```

Avaa `http://localhost:8080`. Oletusajo aloittaa oikeat FMI-haut. Ohjelman voi
pysäyttää painamalla Ctrl+C. Kehitysympäristö, asetukset ja testausohje:
[docs/development.md](docs/development.md).
