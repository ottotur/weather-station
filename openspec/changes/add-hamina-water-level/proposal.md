# Proposal

## Why

Haminan vedenkorkeushavainnon ja N2000-pääennusteen saatavuus on selvitetty. Seuraavaksi tarvitaan ensimmäinen kokonainen toiminto, jolla viimeisin mittaus ja tuleva muutos ovat katsottavissa puhelimella myös FMI-häiriön aikana ilman, että vanha tieto näyttää tuoreelta.

## What Changes

- Lisätään yksi kokonainen ketju: FMI-haku → tarkistus ja normalisointi → pysyvä tiedostovälimuisti → paikallinen JSON-rajapinta → suomenkielinen Svelte-mobiilinäkymä.
- Rajataan asema Hamina Pitäjänsaareen (FMISID 134254, ennusteen geoid -10022823). Havainto muunnetaan millimetreistä senttimetreiksi N2000:ssa; `SeaLevelN2000` tulkitaan senttimetreiksi käyttäjän hyväksymänä oletuksena, ei suoraan varmennettuna metatietona.
- Näytetään viimeisin havainto, enintään seuraavien 48 tunnin todellinen ennuste, kaavio, tuntitaulukko sekä suurin nousu ja lasku ajankohtineen. Käyrää ei siirretä havaintotasolle.
- Toteutetaan käyttäjän hyväksymät vanhentumisrajat: havainto 30 min mittaushetkestä, ennuste 12 h mallin viiteajasta. Vanhentunut tieto merkitään; muutossummat piilotetaan, jos jompikumpi aineisto on vanhentunut tai puuttuu.
- Haetaan havainto 10 min ja ennuste 30 min välein myös virheiden jälkeen; viimeinen tarkistettu aineisto säilytetään. Hakuvirhe ja aineiston ikä esitetään erillisinä.
- Mahdollistetaan paikallinen kehitys ja koko ketjun toistettavat testit Windows 11:n Ubuntu WSL:ssä sekä Windows-selaimessa. Pi-laitehyväksyntä raportoidaan erikseen eikä ristiinkäännöstä tulkita laiteajoksi.

## Capabilities

### New Capabilities

- `water-level-data`: Haminan vedenkorkeusaineistojen keskitetty haku, tarkistus, yksiköt, tuoreus, pysyvä välimuisti ja paikallinen rajapinta.
- `water-level-view`: Vedenkorkeushavainnon, ennusteen ja lasketun muutoksen mobiiliesitys sekä puute- ja häiriötilat.

### Modified Capabilities

Ei olemassa olevia pääspeksifikaatioita.

## Impact

Projektissa ei vielä ole sovelluskoodia tai rakennus-/testiasetuksia. Tuleva toteutus lisää Go-palvelimen, Svelte + TypeScript + Vite -käyttöliittymän ja niiden testaus- ja rakennusohjeet. Palvelin tarjoaa myös käyttöliittymän staattiset tiedostot; selain ei hae FMI:ltä eikä CDN:iltä.

Hyväksytty tuotantokohde on Raspberry Pi 3 Model B+ / Raspberry Pi OS Lite 64-bit / Debian 13 Trixie ja Go `linux/arm64`. Pi:n toimintakunto ja laiteajo ovat testaamatta. Tämä muutos ei kata muita sääaineistoja, asemavalintaa, epävarmuusrajojen selvittämistä, ilmoituksia, täyttä systemd-/julkaisu-/ylläpitoratkaisua tai puhelimen offline-tallennusta.

Tämä ehdotus ei oikeuta toteutuksen aloittamista, Pi:n alustamista tai haarojen yhdistämistä `main`iin. Arkistointi on erillisessä haarassa `docs/archive-water-level-feasibility` (commit `b50b789`); tämä haara alkaa `main`ista eikä sisällä sitä.
