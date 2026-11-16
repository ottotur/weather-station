# Proposal

## Why

`docs/PROJECT_BRIEF.md` edellyttää laite- ja lähdeaineistotarkistusta ennen sovelluksen teknisten valintojen lukitsemista. Vaihetta 0 ei ole tehty, joten yhden aseman vedenkorkeustoiminnon toteutus suunniteltaisiin nyt vielä varmentamattomien oletusten varaan.

## What Changes

- Toteutetaan rajattu osa vaiheesta 0: tarkistetaan Raspberry Pi 2 Model B:n ajotuki sekä yhden käyttäjän hyväksymän vedenkorkeusaseman FMI-havainto ja -ennuste.
- Kootaan tulokset tiedostoon `docs/feasibility/water-level.md`: laitteelta saadut tiedot, dokumentaatiolähteet, toistettavat kyselyt, todellisten vastausten havainnot, käyttöehdot ja avoimet päätökset.
- Tallennetaan vain tarvittavat julkiset esimerkkivastaukset ja niiden alkuperätiedot hakemistoon `docs/feasibility/evidence/`. Arkaluonteisia yhteystietoja tai mökkikoordinaatteja ei tallenneta.
- Päätetään näytön perusteella, voidaanko yhden aseman ensimmäinen kokonainen toiminto suunnitella vai tarvitaanko käyttäjän päätös esteestä.
- Ei toteuteta palvelinta, käyttöliittymää, ajastusta, välimuistia, palveluasennusta tai julkaisua. Ei valita mökin tai ulapan pisteitä eikä selvitetä muita sääaineistoja.

## Capabilities

### New Capabilities

Ei uusia sovelluksen toiminnallisia vaatimuksia. Muutos tuottaa selvitysdokumentaation; `.openspec.yaml` sisältää `skip_specs: true`.

### Modified Capabilities

Ei muutoksia. Projektissa ei vielä ole pääspeksifikaatioita.

## Impact

- Muutoksen toteutus lisää vain selvitysdokumentaatiota ja rajattuja julkisia lähdenäytteitä. Sovelluskoodia tai uusia ajonaikaisia riippuvuuksia ei synny.
- Laiteversion, käyttöjärjestelmän ja arkkitehtuurin tarkistus tarvitsee käyttäjän antamat tiedot tai erikseen myönnetyn laiteyhteyden. Käyttöjärjestelmää ei asenneta eikä laitteen asetuksia muuteta.
- FMI-tarkistus käyttää vain dokumentoituja avoimia rajapintoja niiden käyttörajojen mukaisesti. Aseman tarkoituksenmukaisuus hyväksytetään käyttäjällä.
- Selvitys ei hyväksy Go:ta, käyttöjärjestelmää tai puuttuvien tietojen korvaamista automaattisesti eikä tarkoita koko vaiheen 0 tai MVP:n valmistumista.
