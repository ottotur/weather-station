# Hyväksyntätilanne

Tarkistettu 2026-11-16. Kehitysympäristö: Windows 11, Ubuntu 26.04.1 LTS
WSL 2.7.14.0, Linux x86_64. Go 1.27.4, Node 24.21.0, npm 11.19.0,
Playwright 1.65.1 / Chromium 159.0.8154.6.

## Kehityskoneen automaattiset tarkistukset

| Tarkistus | Tulos |
| --- | --- |
| `go test ./...` | Hyväksytty |
| `go test -race ./...` | Hyväksytty WSL/x86_64 |
| `go vet ./...` | Hyväksytty |
| `npm run check` ja `npm run build` | Hyväksytty |
| `WEATHER_E2E_TMPDIR=/tmp/opencode npm test` (`web/`) | 9 Chromium-testiä hyväksytty; mukana oikea Go-palvelin ja paikallinen WFS-valepalvelin |
| `WEATHER_TEST_WEB_DIST=../../web/dist go test -race ./cmd/weather-station` | Rakennetut staattiset tiedostot Go-palvelimelta: hyväksytty |
| `go build` hostille | Linux/amd64-käännös onnistui |
| `GOOS=linux GOARCH=arm64 GOARM64=v8.0 CGO_ENABLED=0 go build` | Ristiinkäännös onnistui; `file` tunnisti aarch64 ELF -ohjelman |

Selaintestit käsittävät 360, 768 ja 1024 CSS-pikselin leveydet,
koonmuutokset ja emuloidun kosketuksen. Emulaatio ei todista oikean laitteen
kosketuksen, selaimen tai suorituskyvyn toimintaa.

2026-11-17 katselmointikorjausten jälkeen suoritettiin uudelleen Go:n
race-testit, `go vet`, Svelte-tyyppitarkistus ja rakennus sekä 12
Playwright-testiä; kaikki onnistuivat. Uudet testit kattavat vastaamatta
jäävän rajapintapyynnön, ennusteesta poistuneen kaaviovalinnan ja päättyneen
ennusteen yhteenvedon syyn. Tämä ei muuta laitehyväksynnän tilaa.

## Windows-selaimen rajattu automaattitarkistus

Windowsin Google Chrome 157.0.8103.52 avasi WSL-palvelimen osoitteessa
`http://localhost:18092/` headless-tilassa. Paikallisella WFS-valepalvelimella
ja kiinteällä testikellolla selaimen renderöity DOM sisälsi havainnon 32,1 cm,
tuoreustilan ja ennustekaavion. WSL:n localhost-välitys siis toimi tässä
kokeessa. Live-FMI:tä, palomuurimuutoksia tai lähiverkkokuuntelua ei käytetty.
Headless-ajo ei tarkista manuaalisesti ulkoasua tai vuorovaikutusta eikä korvaa
oikean laitteen kosketusta.

## Käyttäjän Windows-selaintarkistus

Käyttäjä ilmoitti 2026-11-16 testanneensa näkymää Windowsissa selaimella
Brave 1.99.57 (64-bittinen, Chromium 157.0.8103.79): hänen mukaansa kaikki
näytti toimivan. Käyttäjä vahvisti erikseen, että tuntitaulukon avaaminen,
kaavion pisteen valinta ja näkymä kapeassa selainikkunassa toimivat oikein.
Ikkunan tarkkaa leveyttä tai kosketusemulaation käyttöä ei ilmoitettu.

Käyttäjän ilmoittamat häiriötestien tulokset:

- Paikallisen Go-palvelimen pysäytyksen jälkeen yhteyden katkeamisesta näkyi
  ilmoitus, aiempi aineisto säilyi näkyvissä ja yhteenvedot poistuivat.
  Palvelimen uudelleenkäynnistyksen jälkeen näkymä palautui normaaliksi.
- Kun lähdeosoitteeksi asetettiin `http://127.0.0.1:1/wfs`, näkyi
  ilmoitus viimeisimmän haun epäonnistumisesta, mutta aiempi havainto ja
  ennuste säilyivät näkyvissä.

Nämä ovat käyttäjän Windows-Bravessa raportoimia tuloksia; WSL:n Playwrightin
kosketusemulaatio on erillinen automaattitesti eikä oikean laitteen testi.
WSL:n localhost-välitys toimi testissä; ohjeistetut vaiheet eivät edellyttäneet
palomuuri- tai verkkoasetusten muuttamista.

## Avoimet hyväksynnät

- Raspberry Pi:n käyttöjärjestelmää ei ole asennettu eikä laitetta ole ajettu
  tai testattu. Ristiinkäännös ei ole Pi-ajotesti.
- Oikeiden puhelinten Safari-, Android Chrome- ja Firefox-testit sekä iPadin
  Safari- ja asennonvaihtotestit ovat tekemättä. Myös Pi:n suorituskyky- ja
  kymmenen asiakkaan laitemittaukset ovat avoinna.
