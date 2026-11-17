# Design

## Context

Motivaatio: `proposal.md`, Why. `docs/PROJECT_BRIEF.md` §4 määrittelee tavoitesuureet ja §7 keskitetyn tiedonkulun. Nykyinen Go-palvelin ja Svelte-näkymä käsittelevät vain Haminaa; FMI-asiakkaan asematunniste ja API-polku ovat kiinteitä. Niitä ei yleistetä tässä selvityksessä.

Aiempaa lähdenäyttöä on `docs/feasibility/water-level.md`:ssä ja sen manifestissa. Sama alkuperäisvastauksen, lähdeviitteen, tarkistuspäivän ja tulkinnan erottelu sopii sääselvitykseen. Pääspeksifikaatioita ei vielä ole; vedenkorkeuden vaatimukset ovat aiemman muutoksen delta-speksifikaatioissa, joita tämä työ ei muuta.

Käyttäjä antoi testipisteen 60.434318, 26.687825 (leveysaste, pituusaste). Sitä saa käsitellä tässä suunnitelmassa testisijaintina, ei lopullisen mökin sijaintina tai edustavuuden todisteena. Ulappapistettä ja tuulihavaintoasemaa ei ole valittu. Pi:n kunnostus on lykätty käyttäjän päätöksellä; käytettävissä on Proxmox-palvelin, jonka kokoonpanoa ei ole tutkittu.

## Goals / Non-Goals

**Goals:**

- Tuottaa toistettava lähdenäyttö, joka erottaa käyttäjän päätökset, dokumentoidut merkitykset, näytteestä tarkistetut tiedot ja avoimet kysymykset.
- Ratkaista seuraavan rajatun sää-/tuulitoiminnon suunnittelukelpoisuus, myös silloin kun kaikkia toissijaisia suureita ei saada.
- Pitää havainto, piste-ennuste ja sijainnin edustavuus erillään.

**Non-Goals:**

- Ei sovelluksen, XML-parserin, sääsymbolikuvaston, yhteisen aineistokehyksen tai asetusten toteutusta.
- Ei tuntuu kuin -laskukaavan, puuskakertoimen tai sateen todennäköisyyden omaa johtamista.
- Ei lopullisen mökkipaikan soveltuvuuden hyväksyntää testipisteen perusteella.
- Ei jatkuvaa lähdeseurantaa, tuulihistoriaa, seitsemän päivän yleiskuvan erillistä lähdeselvitystä tai aallokkoa.

## Decisions

### 1. Dokumentaatioselvitys ilman uusia sovellusvaatimuksia

Tuotetaan `docs/feasibility/weather-wind.md` sekä `docs/feasibility/evidence/weather-wind/manifest.md` ja tarpeelliset muuttamattomat lähdevastaukset. Manifesti kirjaa kunkin näytteen tarkoituksen, täsmällisen URL:n, UTC-hakuajan, HTTP-tuloksen, tiedoston ja SHA-256-tiivisteen. Lähdedokumentaation URL ja tarkistuspäivä kirjataan johtopäätösten yhteyteen.

Vaihtoehto olisi lukita uuden sovellustoiminnon vaatimukset heti. Se hylätään, koska tarkat saatavilla olevat suureet ja niiden merkitykset ovat selvityksen kohde. Siksi `skip_specs: true`; hyväksymisehdot ovat tämän dokumentin kohdassa 6, eivät keksittyjä sovellusominaisuuksia.

### 2. Sijaintivalinnat porttina aineistonäytteille

Testimökkipiste on annettu leveys-/pituusasteina; rajapinnan vaatima koordinaattijärjestelmä ja järjestys varmennetaan dokumentaatiosta. Mahdollinen palautettu mallipiste kirjataan erikseen pyydetystä pisteestä.

Selvitetään virallisista sijainti- ja asematiedoista enintään kaksi perusteltua ulappapiste- ja kaksi tuulihavaintoasemaehdokasta. Kirjataan koordinaatit, aseman nimi ja tunniste sekä valintaperuste: etäisyys, meri-/maa-altistus ja saatavilla olevat suureet. Lähin asema ei automaattisesti ole edustavin. Testialueen ominaisuuksia ei päätellä lopullisen mökkipaikan kuvauksesta.

Esitetään yksi suositeltu yhdistelmä käyttäjälle. Metatietohakuja voi tehdä ehdokkaiden selvittämiseksi, mutta ulapan ja aseman varsinaiset aineistonäytehaut odottavat valinnan hyväksyntää. Testimökkipisteen lähteitä voi tarkistaa erikseen. Valinnan puuttuessa kirjataan täsmällinen esto eikä merkitä koko selvitystä valmiiksi.

### 3. Pieni lähdematriisi, ei kaikkien mallien vertailua

Valitaan dokumentoitu FMI-lähde tavoitellen seuraavaa 48 tuntia ensisijaisesti tunnin askelin. Varmennetaan ensin ehdokaskyselyiden parametrit ja merkitykset; tarkkoja kyselytunnisteita tai kenttänimiä ei arvata suunnittelussa. Lisälähdettä arvioidaan vain, jos valittu lähde jättää seuraavan toiminnon kannalta olennaisen puutteen.

Raportin matriisi kattaa:

| Kohde | Selvitettävät suureet |
| --- | --- |
| Testimökkipisteen ennuste | Lämpötila, keskituuli, puuskat, tuulen suunta, säätila/symbolikoodi ja tuntisade |
| Ehdolliset mökkisuureet | Suora näennäinen lämpötila ja sateen todennäköisyys, jos dokumentoidusti saatavilla |
| Hyväksytty ulappapiste | Keskituuli, puuskat ja suunta |
| Hyväksytty havaintoasema | Viimeisin saatavilla oleva keskituuli, puuskat ja suunta mittausaikoineen |

Jokaisesta suureesta kirjataan tarkka kenttä, yksikkö, määritelmä, mahdollinen keskiarvoistus-/kertymäjakso, puuttuvan arvon esitys sekä saatavuus. Tuulensuunta tulkitaan tulosuunnaksi vain dokumentoidun merkityksen perusteella. Sääsymbolikoodin merkityksen lähde kirjataan; kuvakkeita ei toteuteta.

Sateesta selvitetään, onko kenttä tuntikertymä, intensiteetti vai malliajon alusta kertynyt arvo, sekä mahdollisen Suomen kalenterivuorokausisumman edellytykset. Summia ei lasketa käyttöön eikä puuttuvaa kertymää täydennetä nollilla. Seitsemän päivän kattavuutta ei ryhdytä tavoittelemaan lisähauilla.

### 4. Rajattu, jäljitettävä todellinen näyttö

Ennen näytehakuja tarkistetaan ajantasaiset käyttöehdot, lähdemerkintä, kyselyrajat ja kyselykuvaukset. Näytteet rajataan yhteen hyväksyttyyn kohteeseen kerrallaan ja vain tarvittaviin suureisiin. Ennusteesta pyydetään enintään 48 h; havainnoista aluksi enintään viimeinen tunti. Jos lähde ei tue rajausta, dokumentoidaan sen todellinen käytäntö ennen hakua.

Tavoite on yksi onnistunut rajattu vastaus kutakin valittua aineistoa kohti sekä tarpeelliset metatiedot. Ei jatkuvaa palvelinajoa, rinnakkaisia massahakuja tai automaattista uusintasilmukkaa. Virheellisen vastauksen jälkeen tarkistetaan syy ennen mahdollista perusteltua uutta hakua; myös epäonnistunut yritys kirjataan.

Näytteistä tarkistetaan kohde, suureet, aikajärjestys, askel, äärellisten arvojen alku ja loppu, aukot sekä todellinen tuleva kattavuus. Mittaus-, mallin viite-, tulos-, voimassaolo- ja hakuajat erotetaan; puuttuvaa julkaisuaikaa ei korvata hakuajalla. Dokumentoitu päivitysrytmi erotetaan yhden näytteen havainnoista.

Seuraavalle toteutukselle voidaan ehdottaa lähdekohtaisia haku- ja vanhentumisrajoja perusteluineen, mutta Haminan 30 min / 12 h rajoja ei kopioida sääaineistoille eikä uusia rajoja käsitellä jo hyväksyttyinä.

### 5. Proxmox on myöhempi testiympäristö, ei uusi kohdealusta

Kirjataan testauksen työnjako: Linux-VM voi myöhemmin palvella integraatio-, katko-, uudelleenkäynnistys-, julkaisu- ja oikeiden mobiiliselainten testejä. VM:n käyttöjärjestelmä, resurssit, arkkitehtuuri, osoitteet ja asennus ratkaistaan erillisessä käyttöönottomuutoksessa, ei tässä selvityksessä.

Oikeilla puhelimilla VM:n palvelua vasten tehdyt testit ovat näyttöä kyseisten selainten toiminnasta, eivät Pi:n suorituskyvystä. VM-tulokset, arm64-ristiinkäännös ja oikea Pi-ajo raportoidaan erillisinä. Issue #1 jää avoimeksi; Pi-tavoitetta tai hyväksymisehtoja ei poisteta.

Vaihtoehtona olisi odottaa Pi:tä ennen muuta kehitystä. Sitä ei tarvita lähteiden selvittämiseen, eikä käyttäjä halua kunnostaa laitetta vielä. Toisaalta VM:n onnistumisesta ei päätellä Pi:n toimivuutta.

### 6. Selvityksen hyväksymisehdot

- Testimökkipiste on merkitty testiksi; ulappapisteen ja havaintoaseman valinnat sekä käyttäjän hyväksyntä on kirjattu alkuperineen.
- Jokaiselle matriisin suureelle on lähteistetty saatavuustulos, merkitys ja yksikkö tai nimenomainen puute. Ei hiljaisia lähde- tai laskentakorvauksia.
- Jokaisesta valitusta aineistosta on kohdistettu todellinen vastaus, toistettavat kyselyt ja tarkistettu aikakattavuus. HTTP 200 yksin ei ole hyväksyntä.
- Näytteiden XML/rakenteellinen luettavuus, manifestin tiivisteet ja raportin paikalliset linkit on tarkistettu offline-toistamisohjeella. Alkuperäisiä vedenkorkeusnäytteitä ei ole muutettu.
- Loppuarvio erottaa varmennetun tiedon, ehdotukset ja avoimet päätökset. Se suosittelee yhtä rajattua seuraavaa toimintoa, esimerkiksi testimökkipisteen lämpötila- ja tuuliennustetta, tai palauttaa ydinaineiston puutteen käyttäjälle.
- Puuttuva toissijainen suure ei automaattisesti estä rajatumpaa jatkoehdotusta. Ydinpuutetta ei ratkaista hyväksymättömällä korvaavalla laskennalla.
- Raportti ei väitä koko MVP:tä, lopullista sijaintia, VM-ajoa tai Pi-/mobiilihyväksyntää tehdyksi. Dokumentaatioselvityksen valmistuminen ei käynnistä sovellustoteutusta.

## Risks / Trade-offs

- [Testipaikka ei vastaa lopullista mökkiä] → Sijaintiin ja aseman edustavuuteen liittyvät päätökset tarkistetaan uudelleen ennen varsinaista käyttöönottoa.
- [Mallin maapiste ja meripiste kuvaavat eri olosuhteita] → Pyydetty ja palautettu sijainti sekä pisteiden roolit dokumentoidaan erikseen.
- [Kenttien nimet peittävät erilaiset keskiarvot tai kertymät] → Yksikön lisäksi tarkistetaan aikajakso ja dokumentoitu semantiikka.
- [Kaikki tavoitesuureet eivät ole saatavilla] → Matriisi säilyttää puutteet näkyvinä; seuraavan toiminnon rajaus palautetaan käyttäjälle.
- [Yksi onnistunut vastaus ei todista jatkuvaa saatavuutta] → Ei saatavuuslupauksia; myöhemmät häiriötestit tehdään erillisessä toteutuksessa.
- [VM antaa liian optimistisen kuvan Pi:stä] → Suorituskykytulokset sidotaan mitattuun ympäristöön ja Pi-varmennus säilytetään.
