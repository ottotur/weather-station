# Spec Delta

## Purpose

Auttaa käyttäjää tarkastelemaan Haminan mitattua vedenkorkeutta ja ennustettua muutosta puhelimella ilman, että ennuste, vanha aineisto tai puuttuva tieto sekoittuvat nykyiseen mittaukseen.

## ADDED Requirements

### Requirement: Nimetty havainto ja todellinen ennustejakso
Näkymän SHALL näyttää asema, viimeisin havainto mittausaikoineen, yksikkö cm N2000 sekä pääennusteen kaavio ja avattava tuntitaulukko. Ennuste SHALL rajata nykyhetkestä seuraavaan 48 tuntiin ja todellisiin saatavilla oleviin arvoihin. Käyrää ei SHALL siirtää havaintotasolle eikä aukkoja tai puuttuvaa loppua jatkaa keinotekoisesti.

#### Scenario: Vajaa ennuste
- **WHEN** seuraavalle 48 tunnille on arvoja vain 30 tunniksi ja yksi sisäinen tunti puuttuu
- **THEN** näkymä ilmoittaa lyhyemmän kattavuuden, näyttää aukon sekä kaaviossa että taulukossa eikä piirrä keksittyä jatkoa

### Requirement: Muutos viimeisimmästä havainnosta
Näkymän SHALL esittää suurin nousu ja lasku ajankohtineen vain, kun sekä havainto että ennuste ovat tuoreita ja tulevia arvoja on saatavilla. Vertailutaso SHALL olla saman aseman viimeisin havainto ja ääripisteet SHALL laskea näytetyn ennustejakson pääennusteesta, ei epävarmuusrajoista.

#### Scenario: Nousu ja lasku
- **WHEN** tuore havainto on +32 cm ja tulevan näytetyn jakson suurin arvo on +50 cm ja pienin +20 cm
- **THEN** yhteenveto näyttää nousua 18 cm ja laskua 12 cm kyseisten ennustehetkien kanssa

#### Scenario: Kaikki ennustearvot havaintoa korkeampia
- **WHEN** kaikki tulevat arvot ylittävät havaintotason
- **THEN** näkymä näyttää suurimman nousun ja kertoo ettei havaintotasosta ennusteta laskua; positiivista erotusta ei nimetä laskuksi

#### Scenario: Vertailu ei ole ajantasainen
- **WHEN** havainto tai ennuste puuttuu tai on vanhentunut
- **THEN** nousu- ja laskuluvut piilotetaan ja näkymä kertoo syyn mutta jäljellä oleva aineisto näkyy oman tilansa mukaisesti

### Requirement: Tuoreus ja häiriö näkyvät yhteenvedossa
Näkymän SHALL erottaa tuore, vanhentunut ja puuttuva aineisto sekä viimeisimmän haun epäonnistuminen. Vanhentunut aineisto SHALL säilyttää näkyvissä selvästi merkittynä. Kun tulevia ennustearvoja ei ole, näkymän SHALL ilmoittaa ettei voimassa olevaa ennustetta ole eikä esittää menneitä pisteitä tulevana ennusteena.

#### Scenario: Ennusteen vanhentuminen kesken avoimen sivun
- **WHEN** malliajon ikä ylittää 12 tuntia sivun ollessa auki
- **THEN** viimeistään seuraavassa minuutin päivityksessä ennuste merkitään vanhentuneeksi ja muutossummat poistetaan

#### Scenario: Ensimmäinen käynnistys ilman aineistoa
- **WHEN** palvelin hakee ensimmäistä aineistoa tai haku epäonnistuu ilman välimuistia
- **THEN** näkymä näyttää vastaavasti lataus- tai häiriötilan eikä nollakorkeutta

### Requirement: Paikallinen palvelin ja selainyhteyden häiriö
Näkymän SHALL tarkistaa paikallinen rajapinta noin minuutin välein ja sivulle palattaessa. Selain ei SHALL hakea FMI:ltä tai ulkoisista CDN-palveluista. Jos paikallinen palvelin ei vastaa, viimeinen selaimessa näkyvä tieto SHALL merkitä yhteyskatkolla eikä jäädä varauksetta tuoreeksi.

#### Scenario: Pi ei ole saavutettavissa
- **WHEN** rajapintapyyntö epäonnistuu jo avatulla sivulla
- **THEN** näkyy palvelinyhteyden häiriö, muutossummat piilotetaan ja aiempi sisältö nimetään viimeksi saaduksi tiedoksi

### Requirement: Suomenkielinen kosketuskäyttö
Näkymän SHALL käyttää suomea ja `Europe/Helsinki`-aikaa riippumatta selaimen aikavyöhykkeestä. Päivä ja tarvittaessa UTC-poikkeama SHALL erottaa saman kellonajan toistumiset. Olennainen sisältö SHALL toimia 360 CSS-pikselin leveydessä ilman sivun vaakavieritystä, eikä tietoa SHALL välittää vain värillä tai hiiren osoittimella.

#### Scenario: Eri aikavyöhyke ja talviaikaan siirtyminen
- **WHEN** käyttäjän selain on toisella aikavyöhykkeellä ja ennusteessa on kaksi Suomen klo 03 -ajankohtaa
- **THEN** ajat näkyvät Suomen ajassa toisistaan erotettavina ja kaavion arvot saa kosketuksella tai tuntitaulukosta

### Requirement: Lähde ja tulkinnan rajat
Näkymän SHALL näyttää FMI:n lähdemerkintä, CC BY 4.0 -linkki ja tiedot sovelluksen muunnoksesta ja lasketuista yhteenvedoista. Ennusteen cm-yksikköoletus SHALL olla näkyvästi saatavilla selitteessä. Näkymä ei SHALL esittää keksittyjä epävarmuusrajoja, turvallisuusarvioita tai toimintaohjeita.

#### Scenario: Ennusteen selite
- **WHEN** käyttäjä tarkastelee ennusteen lähde- ja yksikköselitettä
- **THEN** hän näkee, että ennusteen cm on hyväksytty oletus puuttuvan suoran metatiedon vuoksi ja että havaintojen mm-arvot on muunnettu cm-arvoiksi
