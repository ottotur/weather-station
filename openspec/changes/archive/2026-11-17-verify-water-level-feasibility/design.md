# Design

## Context

Perustelu: `proposal.md`, kohta Why. Projektissa on projektikuvaus ja OpenSpec-runko mutta ei sovelluskoodia, testejä, rakennusasetuksia tai olemassa olevia pääspeksifikaatioita. Käyttäjä on vahvistanut, ettei vaihetta 0 ole tehty.

Projektikuvauksen §6, §8, §13 ja §15 rajaavat selvitystä: Svelte + TypeScript + Vite on vahvistettu, Go ja käyttöjärjestelmä alustavia, vedenkorkeuden vertailutaso N2000 ja lähteet FMI:n dokumentoituja avoimia rajapintoja. Dokumentaation ennustepituuksien ristiriitaa ei voi ratkaista olettamalla.

Design tarvitaan, jotta erotetaan todistusaineisto, tekninen ehdotus ja käyttäjän hyväksymä päätös ennen ensimmäisen toiminnon suunnittelua.

Käyttäjän tarkentamat lähtötiedot: laite on piirilevymerkinnän mukaan Raspberry Pi 3 Model B+, vuosimerkintä 2017, ei aiemmin oletettu Pi 2 Model B. Laite ei ole tällä hetkellä käyttökunnossa. Käyttäjä aikoo alustaa sen ennen sääaseman asennusta ja sallii tarkoitukseen parhaiten sopivan käyttöjärjestelmän ehdottamisen. Nykyistä asennusta ei tarvitse tutkia. Käyttäjä on valinnut ensimmäiseksi vedenkorkeusasemaksi Haminan. Nämä ovat käyttäjän ilmoittamia tietoja ja päätöksiä, eivät laiteajolla tai FMI-vastauksilla varmennettuja tuloksia.

## Goals / Non-Goals

**Goals:**

- Tuottaa pieni, toistettava näyttöketju ensimmäisen vedenkorkeustoiminnon riippuvuuksista ilman sovellustoteutusta.
- Erotella jokainen keskeinen tulos tilaan varmennettu, alustava tai avoin/estynyt, sekä liittää siihen lähde ja tarkistusaika.
- Antaa yksiselitteiset ehdot jatkosuunnitteluun siirtymiselle; kielteinenkin löydös on kelvollinen selvitystulos.

**Non-Goals:**

- Ei luvata 48 tunnin ennustetta, täsmällistä vanhentumisrajaa tai laitteen suorituskykyä ilman näyttöä.
- Ei rakenneta kokeellistakaan sovelluspalvelua tai käyttöliittymää, vaihdeta käyttöjärjestelmää eikä tehdä Pi-asennusta.
- Ei muuteta projektikuvausta eikä ratkaista sen avoimia valintoja hiljaisesti.

## Decisions

### 1. Selvitys dokumentoidaan, sovelluskäyttäytymistä ei määritellä vielä

Tulosten pääpaikka on `docs/feasibility/water-level.md`. Rakenne:

1. Rajaus, tarkistuspäivä ja tiivis jatkosuunnittelun tila.
2. Pi:n lähtötiedot ja tuen lähteet.
3. Käyttäjän hyväksymä asema: nimi, lähteiden tunnisteet ja valinnan peruste.
4. Havainto- ja ennustelähteiden tarkistusmatriisi.
5. Toistamisohjeet ja esimerkkivastausten alkuperätiedot.
6. Esteet, ehdotukset sekä käyttäjän päätökset ja vielä puuttuvat hyväksynnät.

`docs/feasibility/evidence/` sisältää enintään muutaman tarpeellisen julkisen vastauksen sekä manifestin, jossa on tiedostonimi, kysely, UTC-hakuaika ja lähdedokumentin URL. Vastaukset säilytetään alkuperäisinä; poikkeamat ja tulkinnat kirjataan erikseen. Ei jatkuvaa aineiston keräystä.

Vaihtoehto olisi kirjoittaa nyt sovelluksen vedenkorkeusspeksi. Se hylätään, koska aineiston saatavuus ja tekniset riippuvuudet eivät ole vielä varmistettuja. Muutos on dokumentaatiomuutos, joten delta-speksit ohitetaan metatiedolla `skip_specs: true`.

### 2. Laitetuki selvitetään mallin ja tulevan puhtaan asennuksen perusteella

Kirjataan käyttäjän ilmoittama Raspberry Pi 3 Model B+ ja piirilevyn vuosimerkintä 2017 sekä tietojen alkuperä. Mallin suoritinarkkitehtuuri ja mahdollinen tuen kannalta olennainen revisiotieto tarkistetaan virallisista laitetiedoista; vuosimerkintää ei tulkita revisioksi. Jos mallimerkintä ei riitä olennaisen tiedon varmistamiseen, pyydetään vain tarvittava lisätieto piirilevystä. Laitteen käynnistämistä, laiteyhteyttä tai nykyisen käyttöjärjestelmän ja käyttäjätilan bittisyyden selvittämistä ei edellytetä. Tunnuksia tai yksityisiä verkkotietoja ei tallenneta projektiin.

Verrataan mallin tietoja Raspberry Pi:n, käyttöjärjestelmän ylläpitäjän ja Go:n virallisiin, ajantasaisiin tukitietoihin. Dokumentoidaan tulevaa puhdasta asennusta varten ehdotettu tuettu käyttöjärjestelmäversio, sen arkkitehtuuri ja käyttäjätilan bittisyys, Go-versio ja käännöskohde tarvittavine ARM-asetuksineen, tukitilanne ja lähteet. Nykyinen asennus merkitään tutkimatta jätetyksi käyttäjän suunnitteleman alustuksen vuoksi; sen tietojen puuttuminen ei yksin estä dokumentoidun asennuskohteen ehdottamista.

Käännöskohde ja käyttöjärjestelmä ovat perusteltuja ehdotuksia, eivät väitteitä tehdystä asennuksesta tai laitteella onnistuneesta ajosta. Käyttäjän lupa ehdottaa sopivaa käyttöjärjestelmää ei tarkoita tietyn version tai Go:n lopullista hyväksyntää. Oikean ohjelman käännös- ja ajotesti kuuluu erilliseen myöhempään toteutusmuutokseen. Jos olennaisia laitetietoja ei saada tai tuettua yhdistelmää ei löydy, merkitään este ja palautetaan jatkopäätös käyttäjälle.

Vaihtoehto olisi asentaa Pi uudelleen tai kirjoittaa testipalvelin tässä muutoksessa. Ne laajentaisivat dokumentaatioselvityksen laite- ja koodimuutokseksi ilman tarvetta.

### 3. FMI-tarkistus kattaa vain yhden hyväksytyn aseman

Ensimmäinen asema on käyttäjän valitsema Hamina. Kirjataan valinta ja hyväksynnän alkuperä raporttiin sekä varmennetaan FMI:n virallisista asematiedoista aseman nimi ja havainto- sekä ennusteaineistojen tunnisteet. Valinnan perusteeksi kirjataan käyttäjän valinta; muita perusteluja ei oleteta eikä mökin tarkkoja koordinaatteja vaadita. Jos aineisto ei ole saatavilla Haminasta, asemaa ei vaihdeta ilman käyttäjän uutta hyväksyntää.

Etsitään dokumentoidut havainto- ja ennustekyselyt, tarkistetaan käyttöehdot ja hakurajat ennen rajattua hakua ja tallennetaan kummastakin todellinen vastaus samalle asemalle. Virhevastauksesta ei päätellä, ettei aineistoa ole: kirjataan virhe, tarkistetaan kysely ja dokumentaatio sekä tarvittaessa kirjataan selvityksen este.

Tarkistusmatriisi sisältää kummallekin aineistolle:

- rajapinnan ja kyselyn tunnisteen, parametrit, asematunnisteen ja vastausmuodon;
- vedenkorkeuskentän merkityksen, yksikön sekä osoituksen N2000-korkeusjärjestelmästä;
- havaintoajat, ennusteen voimassaoloajat ja julkaisu-/malliajon ajan, jos sellainen on saatavilla, sekä aikavyöhykkeen merkityksen;
- viimeisimmän havainnon ja ennusteen ensimmäisen sekä viimeisen ajan, aikavälin ja mahdolliset aukot;
- pääennusteen erottamisen mahdollisista dokumentoiduista epävarmuusrajoista;
- puuttuvien arvojen ja virheiden esitystavan;
- dokumentoidun julkaisurytmin, käyttörajat, lisenssin ja lähdemerkintävelvoitteen.

Tavoiteltua 48 tuntia verrataan todelliseen vastaukseen; puuttuvaa loppuosaa ei oleteta saatavaksi. Julkaisuajaksi ei tulkita hakuaikaa ilman lähteen tukea. Vanhentumisen määrittelyä varten kirjataan käytettävissä olevat aikatiedot ja niiden rajoitteet, mutta ei keksitä numeerista rajaa yhdestä vastauksesta.

Vaihtoehto olisi tarkistaa kaikki MVP:n FMI-aineistot kerralla. Se hylätään, jotta tämä muutos pysyy itsenäisenä yhden toiminnon edellytystarkistuksena.

### 4. Loppuarvio erottaa selvityksen valmistumisen toteutuksen valmiudesta

Selvitys voidaan hyväksyä dokumentaationa myös silloin, kun se löytää esteen. Jatkosuunnittelun tilaksi merkitään kuitenkin **estynyt**, jos olennaiset lähtötiedot tai hyväksynnät puuttuvat.

Yhden aseman toiminnon jatkosuunnittelu on mahdollista vasta, kun:

- Pi:n malliin perustuvat lähtötiedot ja virallisiin lähteisiin perustuva tuettu ajokohde tulevalle puhtaalle asennukselle on dokumentoitu; alustavien teknologiavalintojen hyväksyntä tai tarvittava jatkopäätös näkyy. Nykyisen asennuksen tietoja ei vaadita;
- käyttäjä on hyväksynyt aseman ja sen havainto sekä pääennuste on varmennettu todellisista vastauksista;
- yksiköt, N2000, aikatietojen merkitys ja todellinen ennustejakso tunnetaan;
- käyttöehdot, hakurajat ja lähdemerkintä voidaan ottaa huomioon;
- ydintarpeen puutteet on palautettu käyttäjälle eikä korvattu toisella lähteellä, omalla muunnoksella tai rajauksen muutoksella ilman hyväksyntää.

Epävarmuusrajojen puuttuminen ei yksin estä jatkamista, koska ne ovat projektikuvauksessa ehdollisia. Ydinpuutteen, kuten puuttuvan N2000-ennusteen, vaikutus palautetaan käyttäjälle päätettäväksi. Muut mökki- ja ulappapisteet sekä aineistot pysyvät vaiheen 0 jäljellä olevina töinä.

## Risks / Trade-offs

- [Pi ei ole käyttökunnossa] → Selvitys perustuu käyttäjän ilmoittamaan mallimerkintään ja virallisiin tukitietoihin. Nykyistä asennusta ei tarvitse tutkia, mutta mallin olennaiset epäselvyydet kirjataan esteiksi ja laiteajo jää myöhempään toteutusmuutokseen.
- [Dokumentaatio ja vastaus eroavat] → Molemmat lähteet ja tarkistusajat kirjataan; ristiriita säilytetään näkyvänä.
- [Yksi vastaus ei todista jatkuvaa saatavuutta tai julkaisurytmiä] → Erotetaan hetkellinen havainto dokumentoidusta lupauksesta; myöhempi toteutus tarvitsee häiriötestit.
- [Go-tuki dokumenteissa ei takaa oikean ohjelman toimintaa Pi:llä] → Ei väitetä ajotestiä tehdyksi; testi jää myöhempään toteutusmuutokseen ennen teknologian lopullista lukitsemista.
- [Julkiset näytteet sisältävät tarpeetonta dataa] → Kysely rajataan yhteen asemaan ja tarvittaviin suureisiin sekä lyhyeen ajanjaksoon; salaisuuksia tai yksityisiä sijainteja ei tallenneta.
- [Selvitys tulkitaan koko vaiheen 0 valmistumiseksi] → Loppuarvio luettelee nimenomaisesti ulkopuolelle jääneet sijainnit ja aineistot.

## Migration Plan

Ei sovellusjulkaisua tai tietomigraatiota. Apply-vaihe lisää vain selvitysraportin ja sen lähdenäytteet tarkistettaviksi. Mahdollinen oikaisu tehdään dokumentteja korjaamalla; Pi:n tila ei muutu. Ensimmäinen sovellustoiminto suunnitellaan tämän jälkeen uutena OpenSpec-muutoksena erillisellä käyttäjän pyynnöllä.
