# Proposal

## Why

`docs/PROJECT_BRIEF.md` edellyttää laite- ja lähdeaineistotarkistusta ennen sovelluksen teknisten valintojen lukitsemista. Vaihetta 0 ei ole tehty, joten yhden aseman vedenkorkeustoiminnon toteutus suunniteltaisiin nyt vielä varmentamattomien oletusten varaan.

## What Changes

- Toteutetaan rajattu osa vaiheesta 0: tarkistetaan käyttäjän ilmoittaman Raspberry Pi 3 Model B+:n (piirilevyn vuosimerkintä 2017) dokumentoitu ajotuki sekä käyttäjän hyväksymän Haminan vedenkorkeusaseman FMI-havainto ja -ennuste.
- Kootaan tulokset tiedostoon `docs/feasibility/water-level.md`: käyttäjän ilmoittamat laitetiedot, dokumentaatiolähteet, toistettavat kyselyt, todellisten vastausten havainnot, käyttöehdot ja avoimet päätökset.
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
- Laite tunnistetaan käyttäjän antamista piirilevymerkinnöistä ja mallin tekniset tiedot tarkistetaan virallisista lähteistä. Laite ei ole tällä hetkellä käyttökunnossa ja käyttäjä aikoo alustaa sen ennen sääaseman asennusta, joten nykyistä käyttöjärjestelmää tai käyttäjätilan bittisyyttä ei tarvitse selvittää. Selvitetään mallille sopiva tuettu käyttöjärjestelmä–Go-yhdistelmä tulevaa puhdasta asennusta varten. Laitetta ei alusteta, käyttöjärjestelmää ei asenneta eikä laitteen asetuksia muuteta tässä muutoksessa.
- FMI-tarkistus käyttää vain dokumentoituja avoimia rajapintoja niiden käyttörajojen mukaisesti. Käyttäjä on valinnut Haminan; aseman virallinen nimi ja aineistojen tunnisteet varmennetaan selvityksessä.
- Selvitys ei hyväksy Go:ta, käyttöjärjestelmää tai puuttuvien tietojen korvaamista automaattisesti eikä tarkoita koko vaiheen 0 tai MVP:n valmistumista.
