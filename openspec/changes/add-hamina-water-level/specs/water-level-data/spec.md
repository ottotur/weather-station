# Spec Delta

## Purpose

Tarjota Haminan tarkistettu vedenkorkeushavainto ja pääennuste paikallisille asiakkaille siten, että yksiköt, aineiston ikä ja häiriöt ovat yksiselitteisiä.

## ADDED Requirements

### Requirement: Keskitetty asemakohtainen haku
Palvelimen SHALL hakea Haminan (FMISID 134254, geoid -10022823) havainnot FMI:n dokumentoidulla mareografikyselyllä 10 minuutin ja ennusteet vedenkorkeuskyselyllä 30 minuutin välein sivukatseluista riippumatta. Virheen jälkeen SHALL odottaa normaalia hakuväliä; tiheää uusintasilmukkaa ei käytetä.

#### Scenario: Useita selaimia ja hakuvirhe
- **WHEN** kymmenen selainta avaa näkymän ja ajastettu ennustehaku epäonnistuu
- **THEN** selaimet eivät käynnistä FMI-hakuja ja seuraava ennusteyritys tehdään normaalin 30 minuutin välin jälkeen

#### Scenario: Ensimmäinen käynnistys
- **WHEN** palvelin käynnistyy ilman välimuistia
- **THEN** molempien aineistojen haku aloitetaan odottamatta ensimmäistä täyttä hakuväliä ja rajapinta voi sillä aikaa palvella puuttuvan aineiston tilaa

### Requirement: Yksiköt ja korkeusjärjestelmä
Palvelimen SHALL käyttää havaintosuuretta `WLEVN2K_PT1S_INSTANT`, muuntaa mm-arvot cm-arvoiksi jakamalla kymmenellä ja käyttää ennusteena vain `SeaLevelN2000`-sarjaa. Ennusteen cm-yksikkö SHALL merkitä käyttäjän hyväksymäksi oletukseksi, ei varmennetuksi metatiedoksi. `SeaLevel` ei SHALL korvata N2000-sarjaa.

#### Scenario: Näyteaineiston normalisointi
- **WHEN** havainto on 321 mm ja ennusteen N2000-arvo on 40
- **THEN** rajapinta tarjoaa 32,1 cm ja 40 cm N2000 sekä ennusteen yksikköoletuksen tilan

### Requirement: Tarkistus ennen julkaisua
Palvelimen SHALL tarkistaa asema, suure, aikaleimat ja aikaviitteet sekä numeeriset arvot ennen aineiston julkaisua. `NaN` ja puuttuvat arvot SHALL käsitellä puuttuvina, ei nollina. Rakenteellisesti virheellinen tai kokonaan käyttökelvoton vastaus ei SHALL korvata viimeistä tarkistettua aineistoa.

#### Scenario: Virheellinen lähdeaineisto
- **WHEN** HTTP 200 -vastaus sisältää WFS-virheen, väärän aseman, ratkaisemattoman sisäisen aikaviitteen tai vain puuttuvia arvoja
- **THEN** yritys merkitään epäonnistuneeksi ja aiempi tarkistettu aineisto säilyy

#### Scenario: Ennusteen puuttuvat tunnit
- **WHEN** ennusteen lopussa tai sisällä on `NaN`-arvoja
- **THEN** rajapinta säilyttää tiedon aukoista eikä tuota näille tunneille nollia tai keksittyjä arvoja

### Requirement: Tuoreuden arviointi lähdeajasta
Palvelimen SHALL merkitä havainto vanhentuneeksi, kun mittaushetken ikä ylittää 30 minuuttia, ja ennuste, kun `analysisTime`-viiteajan ikä ylittää 12 tuntia. Haku-, tulos- ja voimassaoloajat SHALL säilyttää erillisinä. Onnistunut uusintahaku ei SHALL nuorentaa samaa havaintoa tai malliajoa.

#### Scenario: Vanhentumisrajan reunat
- **WHEN** havainto on tasan 30 minuuttia tai ennuste tasan 12 tuntia vanha
- **THEN** kyseinen aineisto on vielä tuore; rajan ylityttyä se on vanhentunut myös ilman uutta FMI-hakua

#### Scenario: Sama malliajo uudelleen
- **WHEN** yli 12 tuntia vanha malliajo haetaan onnistuneesti uudelleen
- **THEN** ennuste pysyy vanhentuneena vaikka viimeisin hakuyritys onnistui

### Requirement: Riippumaton pysyvä välimuisti
Palvelimen SHALL säilyttää havainto ja ennuste erillisinä viimeisinä tarkistettuina aineistoina uudelleenkäynnistyksen yli. Uusi aineisto SHALL julkaista muistissa vasta onnistuneen tarkistuksen, atomisen korvauksen ja hakemiston synkronoinnin jälkeen. Ennen korvausta tapahtuva tallennusvirhe SHALL säilyttää aiempi kohdetiedosto ja muistissa julkaistu aineisto. Yhden aineiston haku-, luku- tai kirjoitusvirhe ei SHALL estää toisen tarjoilua.

#### Scenario: Katkennut tallennus ja uudelleenkäynnistys
- **WHEN** ennusteen tallennus keskeytyy ennen korvausta ja palvelin käynnistetään uudelleen
- **THEN** viimeinen kokonaan tallennettu ennuste voidaan lukea eikä osittaista tiedostoa julkaista

#### Scenario: Vioittunut havaintovälimuisti
- **WHEN** havaintovälimuisti ei ole luettavissa mutta ennustevälimuisti on kelvollinen
- **THEN** havainto ilmoitetaan puuttuvaksi ja ennuste tarjotaan oman ikänsä mukaisessa tilassa

### Requirement: Korvauksen jälkeinen epävarma tallennustila
Jos korvaus onnistuu mutta hakemiston synkronointi epäonnistuu, palvelimen SHALL säilyttää aiempi muistissa julkaistu aineisto ja ilmoittaa erillisellä syykoodilla tallennuksen epävarma tila. Uusi tiedosto voi jo olla levyllä ja tulla ladatuksi uudelleenkäynnistyksessä.

#### Scenario: Hakemiston synkronointi epäonnistuu korvauksen jälkeen
- **WHEN** tarkistettu uusi aineisto on atomisesti korvannut kohdetiedoston mutta hakemiston synkronointi epäonnistuu
- **THEN** muistissa tarjotaan aiempi aineisto ja raportoidaan tallennuksen epävarma tila; uudelleenkäynnistyksessä voidaan lukea uusi tiedosto eikä aiemman tiedoston säilymistä luvata

### Requirement: Paikallinen rajapinta erottaa datan ja häiriön
Palvelimen SHALL tarjota `GET /api/water-level/hamina` JSON-muodossa myös osittaisessa tai tyhjässä datatilassa. Vastauksen SHALL sisältää palvelinaika, asema, cm/N2000, lähdeajat, tuoreus, todellinen kattavuus, puutteet ja viimeisen hakuyrityksen tulos kummallekin aineistolle. Ei-äärellisiä JSON-lukuja ei SHALL lähettää.

#### Scenario: Katko tuoreen välimuistin aikana
- **WHEN** FMI-haku epäonnistuu mutta aiempi havainto on alle 30 minuuttia vanha
- **THEN** rajapinta ilmoittaa havainnon tuoreeksi ja hakuyrityksen epäonnistuneeksi erillisinä tietoina

#### Scenario: Tyhjä käynnistys ja internetkatko
- **WHEN** kumpaakaan aineistoa ei ole ja molemmat haut epäonnistuvat
- **THEN** rajapinta palauttaa selkeät puuttumis- ja virhetilat ilman tekaistuja nolla-arvoja

### Requirement: Rajattu ja jäljitettävä lähdekäyttö
Palvelimen SHALL käyttää vain dokumentoituja FMI-kyselyitä yhdelle asemalle, noudattaa käyttörajoja ja estää päällekkäiset saman aineiston haut. Vastauksen SHALL tunnistaa FMI lähteeksi, CC BY 4.0 -lisenssi sekä sovelluksen tekemä yksikkömuunnos ja laskenta.

#### Scenario: Hidas lähde ja monta asiakasta
- **WHEN** FMI-vastaus viipyy ja paikallisia rajapintapyyntöjä tulee samanaikaisesti
- **THEN** uusia päällekkäisiä saman aineiston FMI-hakuja ei synny ja asiakkaat saavat välimuistivastauksen lähde- ja käsittelytietoineen
