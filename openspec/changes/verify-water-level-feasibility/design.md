# Design

## Context

Perustelu: `proposal.md`, kohta Why. Projektissa on projektikuvaus ja OpenSpec-runko mutta ei sovelluskoodia, testejä, rakennusasetuksia tai olemassa olevia pääspeksifikaatioita. Käyttäjä on vahvistanut, ettei vaihetta 0 ole tehty.

Projektikuvauksen §6, §8, §13 ja §15 rajaavat selvitystä: Svelte + TypeScript + Vite on vahvistettu, Go ja käyttöjärjestelmä alustavia, vedenkorkeuden vertailutaso N2000 ja lähteet FMI:n dokumentoituja avoimia rajapintoja. Dokumentaation ennustepituuksien ristiriitaa ei voi ratkaista olettamalla.

Design tarvitaan, jotta erotetaan todistusaineisto, tekninen ehdotus ja käyttäjän hyväksymä päätös ennen ensimmäisen toiminnon suunnittelua.

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

### 2. Laitetuki varmennetaan lukevilla tarkistuksilla

Pyydetään käyttäjältä Pi:n mallin/revision, suoritinarkkitehtuurin, käyttöjärjestelmän version ja käyttäjätilan bittisyyden tiedot. Tarvittaessa annetaan lukevat tarkistuskomennot; laiteyhteyttä ei oleteta eikä tunnuksia tallenneta projektiin.

Verrataan tietoja Raspberry Pi:n, käyttöjärjestelmän ylläpitäjän ja Go:n virallisiin, ajantasaisiin tukitietoihin. Dokumentoidaan mahdollinen tuettu käyttöjärjestelmäversio, Go-versio ja käännöskohde tarvittavine ARM-asetuksineen, tukitilanne ja lähteet. Pelkkä prosessorin kyvykkyys ei todista asennetun käyttöjärjestelmän arkkitehtuuria tai tukea.

Käännöskohde on perusteltu ehdotus, ei väite laitteella onnistuneesta ajosta. Oikean ohjelman käännös- ja ajotesti kuuluu erilliseen myöhempään toteutusmuutokseen. Jos tietoja ei saada tai tuettua yhdistelmää ei löydy, merkitään este eikä vaihdeta automaattisesti palvelinkieltä tai käyttöjärjestelmää.

Vaihtoehto olisi asentaa Pi uudelleen tai kirjoittaa testipalvelin tässä muutoksessa. Ne laajentaisivat dokumentaatioselvityksen laite- ja koodimuutokseksi ilman tarvetta.

### 3. FMI-tarkistus kattaa vain yhden hyväksytyn aseman

Pyydetään käyttäjää valitsemaan vedenkorkeusasema. Tarvittaessa esitetään FMI:n asematiedoista löydetyt sopivat vaihtoehdot käyttäjän ilmoittamalta alueelta; ei vaadita mökin tarkkoja koordinaatteja. Asemaa ei päätetä pelkän teknisen saatavuuden perusteella.

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

- Pi:n lähtötiedot ja virallisiin lähteisiin perustuva tuettu ajokohde on dokumentoitu; alustavien teknologiavalintojen hyväksyntä tai tarvittava jatkopäätös näkyy;
- käyttäjä on hyväksynyt aseman ja sen havainto sekä pääennuste on varmennettu todellisista vastauksista;
- yksiköt, N2000, aikatietojen merkitys ja todellinen ennustejakso tunnetaan;
- käyttöehdot, hakurajat ja lähdemerkintä voidaan ottaa huomioon;
- ydintarpeen puutteet on palautettu käyttäjälle eikä korvattu toisella lähteellä, omalla muunnoksella tai rajauksen muutoksella ilman hyväksyntää.

Epävarmuusrajojen puuttuminen ei yksin estä jatkamista, koska ne ovat projektikuvauksessa ehdollisia. Ydinpuutteen, kuten puuttuvan N2000-ennusteen, vaikutus palautetaan käyttäjälle päätettäväksi. Muut mökki- ja ulappapisteet sekä aineistot pysyvät vaiheen 0 jäljellä olevina töinä.

## Risks / Trade-offs

- [Pi ei ole saavutettavissa] → Käyttäjälle tarkat tietopyynnöt; ei arvauksia, ja jatkosuunnittelu merkitään tarvittaessa estyneeksi.
- [Dokumentaatio ja vastaus eroavat] → Molemmat lähteet ja tarkistusajat kirjataan; ristiriita säilytetään näkyvänä.
- [Yksi vastaus ei todista jatkuvaa saatavuutta tai julkaisurytmiä] → Erotetaan hetkellinen havainto dokumentoidusta lupauksesta; myöhempi toteutus tarvitsee häiriötestit.
- [Go-tuki dokumenteissa ei takaa oikean ohjelman toimintaa Pi:llä] → Ei väitetä ajotestiä tehdyksi; testi jää myöhempään toteutusmuutokseen ennen teknologian lopullista lukitsemista.
- [Julkiset näytteet sisältävät tarpeetonta dataa] → Kysely rajataan yhteen asemaan ja tarvittaviin suureisiin sekä lyhyeen ajanjaksoon; salaisuuksia tai yksityisiä sijainteja ei tallenneta.
- [Selvitys tulkitaan koko vaiheen 0 valmistumiseksi] → Loppuarvio luettelee nimenomaisesti ulkopuolelle jääneet sijainnit ja aineistot.

## Migration Plan

Ei sovellusjulkaisua tai tietomigraatiota. Apply-vaihe lisää vain selvitysraportin ja sen lähdenäytteet tarkistettaviksi. Mahdollinen oikaisu tehdään dokumentteja korjaamalla; Pi:n tila ei muutu. Ensimmäinen sovellustoiminto suunnitellaan tämän jälkeen uutena OpenSpec-muutoksena erillisellä käyttäjän pyynnöllä.
