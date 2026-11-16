# Projektikuvaus: mökin sääsovellus

Tämä dokumentti kokoaa suunnitteluhaastattelun tulokset. Se on myöhemmän vaatimusmäärittelyn ja OpenSpec-työskentelyn lähtökohta, ei vielä toteutussuunnitelman tai kaikkien teknisten yksityiskohtien lopullinen hyväksyntä.

## 1. Päätösten tila

- **Vahvistettu:** käyttäjän haastattelussa hyväksymä vaatimus tai rajaus.
- **Alustava:** hyväksytty suunnittelusuunta tai ehdotus, jonka yksityiskohdat tai toteutettavuus pitää vielä varmistaa.
- **Avoin:** päätettävä tai selvitettävä asia. Avoimia asioita ei saa ratkaista huomaamatta toteutuksessa.

Haastattelun yksittäiset päätökset on vahvistettu keskustelussa. Viimeiseen koottuun yhteenvetoon ei annettu erillistä kokonaisuuden hyväksyntää ennen tämän dokumentin laatimispyyntöä.

## 2. Tavoite ja käyttöympäristö

**Vahvistettu**

Sovellus kokoaa Ilmatieteen laitoksen (FMI) sää- ja meridatan helposti hahmotettavaksi verkkokäyttöliittymäksi mökin lähiverkkoon. Ensisijaisia käyttölaitteita ovat käyttäjien omat älypuhelimet ja tabletit.

Sovelluksen palvelin toimii kokonaan projektille varatulla Raspberry Pi 3 Model B+:llä. Mökillä on normaalisti jatkuva internetyhteys.

Mökki sijaitsee saaressa itäisellä Suomenlahdella. Lännessä on saaristoa, mutta etelässä ja kaakossa avointa merta yli 40 kilometrin matkalla. Mökkipaikan ja sen edustalla olevan ulapan olosuhteet esitetään erillisinä.

### Käyttötarpeet

1. **Mökkipäivän suunnittelu:** lämpötila, näennäinen lämpötila, tuuli, sade ja sääennuste.
2. **Venematkan olosuhteiden tarkastelu:** erityisesti tuuli, aallokko ja vedenkorkeus.
3. **Rantarakenteiden suojaaminen:** tulevan veden nousun ja sen ajankohdan havaitseminen, jotta lähellä merta olevat rakenteet voidaan tarvittaessa kiinnittää paremmin.

Sovellus näyttää tiedot, mutta ei arvioi venematkan turvallisuutta tai anna toimintaohjeita. Käyttäjä tekee tulkinnan itse. Sovellus ei takaa ennusteiden toteutumista.

## 3. Käyttäjät ja kehittämisen tavoite

**Vahvistettu**

- Käyttäjät tuntevat paikalliset olosuhteet ja osaavat tulkita sää- ja meriennusteita.
- Käyttäjä valitsee itse, minkä vedenkorkeusaseman tiedot ovat kulloinkin hyödyllisimpiä.
- Kaikki mökin lähiverkon laitteet saavat katsella sovellusta ilman kirjautumista.
- MVP:tä käytetään selaimessa. Erillistä sovellusasennusta tai jatkuvasti päällä olevan yhteisen näytön tilaa ei tarvita.

Kehittäjä toteuttaa sovellusta itse ja haluaa oppia erityisesti spec-driven-työskentelyä ja OpenSpecin käyttöä. Taustalla on WordPress-kehitystä sekä PHP-, SCSS-, JavaScript- ja TypeScript-osaamista ja Linux-palvelinten perusylläpitoa.

Teknologioita ei valita vain nykyisen kieliosaamisen perusteella. Tavoitteita ovat kevyt toiminta Pi:llä, nopea käyttöliittymä ja ylläpidettävyys. Toteutuksen vastuut ja tiedonkulku pitää selittää niin, että kehittäjä ymmärtää myös uudella kielellä kirjoitetun koodin.

## 4. MVP:n ominaisuudet

### 4.1 Sijainnit ja asetukset

**Vahvistettu**

- Yksi kiinteä mökkipaikka.
- Yksi kiinteä ulapan ennustepiste.
- Ennalta valittu joukko vedenkorkeusasemia.
- Koordinaatit ja asemavalinnat määritellään asetustiedostossa.
- Ei paikkahakua, karttavalintaa tai asetusten muokkausta selaimessa.

**Avoin:** tarkat koordinaatit, asematunnisteet ja tuulihavaintoihin soveltuva sääasema.

### 4.2 Mökkisää

**Vahvistettu tavoitesisältö**

- Lämpötila.
- Näennäinen eli ”tuntuu kuin” -lämpötila.
- Säätila ja sääsymboli.
- Keskituulen nopeus, puuskanopeus ja tuulen suunta.
- Tuntikohtainen sademäärä ja kalenterivuorokauden sadekertymä.

Etusivun mökkisää perustuu mökkipaikan kuluvan tunnin ennusteeseen. Se nimetään esimerkiksi ”Mökin ennuste klo 14”, ei paikalla mitatuksi nykytilanteeksi.

Sateen todennäköisyys voidaan näyttää vain, jos se saadaan valitusta FMI-aineistosta suoraan ja sen merkitys tunnetaan. Sitä ei päätellä sademäärästä.

**Avoin:** näennäisen lämpötilan saatavuus. Jos se pitää laskea itse, laskentatapa hyväksytään ja dokumentoidaan erikseen.

### 4.3 Tuuli

**Vahvistettu**

- Mökkipaikan ja ulapan tuuliennusteet pidetään erillisinä.
- Esitetään keskituuli ja puuskat yksikössä m/s sekä suunta nuolena ja ilmansuuntana.
- Tuulen suunta tarkoittaa suuntaa, josta tuuli tulee. Nuolen merkitys selitetään.
- Puuskia ei lasketa itse keskituulesta, jos niitä ei ole aineistossa.
- Soveltuvan lähiaseman viimeisimmät tuulihavainnot näytetään erikseen aseman nimellä ja mittausajalla.
- Tuulihistoria ei kuulu MVP:hen.

### 4.4 Aallokko

**Vahvistettu tavoitesisältö, saatavuuden rajoissa**

- Merkitsevä aallonkorkeus.
- Aallokon tulosuunta.
- Aaltojakso lähdeaineiston määritelmän mukaisesti.

Aallonkorkeus on ensisijainen suure. Merkitsevä aallonkorkeus ei tarkoita suurinta yksittäistä aaltoa; käyttöliittymässä tarjotaan lyhyt selite. Ulapan ennustetta ei esitetä mökkirannan paikallisena aallokkona.

### 4.5 Vedenkorkeus

**Vahvistettu**

- Valittujen asemien viimeisimmät havainnot ja ennusteet.
- Yksikkö cm ja korkeusjärjestelmä N2000.
- Erityinen painotus muutokseen viimeisimmästä mitatusta tilanteesta.
- Asemakohtaisesti havaittu korkeus ja mittausaika, ennustekäyrä sekä ennustejakson suurin nousu ja lasku ajankohtineen.
- Muutoksen vertailutaso on kunkin aseman oma viimeisin havainto.
- Havainto ja ennuste näytetään sellaisinaan. Ennustekäyrää ei siirretä havainnon tasolle.
- Jos havainto puuttuu tai on vanhentunut, muutosta ei esitetä luotettavana nykyhetkestä laskettuna lukuna.
- Epävarmuusalue näytetään, jos FMI tarjoaa dokumentoidut rajat.
- Nousu- ja laskuyhteenvedot lasketaan pääennusteesta, eivät epävarmuusalueen äärirajoista.
- Epävarmuusrajoja ei nimetä mahdollisiksi absoluuttisiksi minimi- ja maksimikorkeuksiksi.

Esitysesimerkki: ”Nyt +35 cm → ennustettu huippu +68 cm huomenna klo 04 → nousua 33 cm”.

### 4.6 Ennusteiden aikajänne

**Vahvistettu tavoite, aineistojen kattavuuden rajoissa**

- Seuraavat 48 tuntia mahdollisimman tarkasti, ensisijaisesti tunneittain.
- Mökkisäästä seuraavien seitsemän vuorokauden päivittäinen yleiskuva.
- Kukin aineisto näytetään vain sen todelliselta ennustejaksolta. Puuttuvaa loppuosaa ei jatketa keinotekoisesti.
- Vuorokausikertymä tarkoittaa Suomen paikallista kalenterivuorokautta, ei liukuvaa 24 tunnin jaksoa.

## 5. Käyttöliittymä

**Vahvistettu**

Etusivun järjestys:

1. Mökkisää.
2. Ulapan olosuhteet.
3. Vedenkorkeudet.

Aluksi näytetään nopeasti silmäiltävä yhteenveto ja tulevan muutoksen olennaiset tiedot. Sen alta löytyvät tarkemmat ennusteet.

- Kaaviot havainnollistavat kehitystä; avattavat tuntitaulukot tarjoavat tarkat arvot.
- Kaavion ajanhetken tiedot saa esiin koskettamalla, ilman hiiren osoitinta.
- Vedenkorkeusasemille omat nimetyt pienet kaaviot yhteisellä aika-akselilla ja korkeusasteikolla.
- Olennaista tietoa ei esitetä vain värillä.
- Luettavuus huomioidaan myös kirkkaassa ulkovalossa.
- Tiedon tuoreus ja puutteet näkyvät jo yhteenvedossa.

Kieli on suomi. Ajat esitetään Suomen paikallisajassa riippumatta käyttäjän laitteen aikavyöhykkeestä. Yksiköt ovat °C, m/s, sade millimetreinä, aallokko metreinä ja vedenkorkeus senttimetreinä N2000-tasosta.

MVP tukee ajantasaisia iPhonen/iPadin Safari-selaimia sekä Androidin Chromea ja Firefoxia. Erityistä vanhojen selainversioiden tukilupausta ei ole annettu.

**Avoin:** visuaalinen ilme, kaaviokirjasto ja tarkat vuorovaikutusratkaisut.

## 6. Alustavat teknologiavalinnat

| Osa | Suunta ja tila | Perustelu |
| --- | --- | --- |
| Palvelin | **Alustava: Go**, hyväksytty ensisijaiseksi vaihtoehdoksi | Kevyt palvelin, mahdollisuus kääntää kehityskoneella Pi:lle ja toimittaa yhtenä ohjelmatiedostona. Tarkka laite- ja käyttöjärjestelmätuki varmistetaan ensin. |
| Käyttöliittymä | **Vahvistettu: Svelte + TypeScript + Vite** | Sopii kompaktiin, päivittyvään näkymään. Rakennetaan staattisiksi tiedostoiksi; Pi ei tarvitse käyttöliittymän kehityspalvelinta tai palvelinrenderöintiä. |
| Tallennus | **Vahvistettu: tiedostopohjainen välimuisti** | MVP ei tarvitse historiaa tai erillistä tietokantapalvelinta. SQLite arvioidaan mahdollisen RuuviTag-ominaisuuden yhteydessä. |
| Käyttöjärjestelmä | **Alustava: tuettu kevyt Linux ilman työpöytää** | Laite on vain tämän sovelluksen käytössä. Tarkka jakelu, versio ja arkkitehtuuri ovat avoimia. |
| Palvelun ajo | **Vahvistettu toimintatavoite: järjestelmäpalvelu** | Automaattinen käynnistys sähkökatkon jälkeen ja uudelleenkäynnistys kaatumisen jälkeen. Tarkka palvelukonfiguraatio suunnitellaan myöhemmin. |
| Kontit | **Alustava: ei Dockeria oletuksena** | Pienelle laitteelle ei lisätä ylläpitokerrosta ilman osoitettua hyötyä. |

### Harkitut vaihtoehdot

- **Node.js/TypeScript:** tuttu kieli myös palvelimelle, mutta ARMv7-tuen ja ajoversioiden saatavuus pitää huomioida. Ei valittu oletukseksi.
- **Python:** todennäköisesti riittävä suorituskyky ja luonteva tiedonkäsittely, mutta erillinen ajoympäristö ja riippuvuudet ylläpidettävänä. Go valittiin ensisijaisesti kevyen käyttöönoton vuoksi, ei siksi, että Python olisi osoitettu liian hitaaksi.
- **React:** riittävä suorituskyky ja laaja ekosysteemi; Svelte valittiin kompaktiin käyttöliittymään.
- **Pelkkä TypeScript ilman kehystä:** vähentäisi kehysriippuvuuksia, mutta siirtäisi näkymän tilanhallintaa käsin tehtäväksi.
- **Palvelinrenderöinti:** ei valittu; kaaviot tarvitsevat joka tapauksessa selainlogiikkaa.

Käyttöliittymän suorituskyky syntyy pääosin käyttäjän laitteessa. Pi:n palvelinkieli ei yksin ratkaise näkymän nopeutta; myös datan koko, kaaviokirjasto ja selainpuolen toteutus vaikuttavat.

## 7. Arkkitehtuuri ja tiedonkulku

**Vahvistettu arkkitehtuurisuunta**

```text
FMI:n dokumentoidut avoimet rajapinnat
                |
                v
Raspberry Pi: ajastettu haku ja aineiston tarkistus
                |
                v
Lähdekohtainen käsittely ja pysyvä tiedostovälimuisti
                |
                v
Kevyt datarajapinta + staattiset käyttöliittymätiedostot
                |
                v
Puhelimen/tabletin selain: Svelte-näkymät ja kaaviot
```

- Pi hakee aineistot keskitetysti, sivun avaamisesta riippumatta.
- Käyttäjän sivulataus ei käynnistä omaa FMI-hakua.
- Selain saa pienen, näkymälle sopivaksi käsitellyn tietopaketin.
- Kaikki käyttöliittymän tarvitsemat tiedostot, myös kirjastot, kuvakkeet ja mahdolliset fontit, tarjotaan Pi:ltä.
- Ei ulkoisia CDN-riippuvuuksia, analytiikkaa tai selaimesta suoraan tehtäviä FMI-hakuja.
- Aineiston haku, käsittely, tallennus ja tarjoilu erotetaan selkeiksi vastuiksi. Tarkkaa moduulirakennetta ei ole päätetty.

### Tallennus ja häiriöt

**Vahvistettu**

- Säilytetään lähdekohtaisesti viimeisin onnistuneesti haettu ja tarkistettu aineisto aikaleimoineen.
- Uusi aineisto korvaa vanhan vasta onnistuneen käsittelyn jälkeen.
- Välimuisti säilyy Pi:n uudelleenkäynnistyksen yli.
- Ei pitkäaikaista sää-, havainto- tai ennustehistoriaa.
- Yhden aineiston häiriö ei estä muiden aineistojen käyttöä.
- Internetkatkossa näytetään viimeksi saatu aineisto selvästi vanhentuneeksi merkittynä.
- Ennustejakson päätyttyä kerrotaan, ettei voimassa olevaa ennustetta ole.
- Puuttuvaa arvoa ei muuteta nollaksi.
- Puhelimen oma offline-tallennus ei kuulu MVP:hen; Pi:n pitää olla saavutettavissa lähiverkossa.

### Päivitysrytmi

**Alustavat, hyväksytyt lähtöarvot**

- Havainnot: Pi tarkistaa uudet tiedot 10 minuutin välein.
- Ennusteet: Pi tarkistaa uudet tiedot 30 minuutin välein.
- Avoin käyttöliittymä tarkistaa Pi:ltä päivitykset noin minuutin välein ja sivulle palattaessa.

Välit tarkennetaan lähteiden julkaisurytmin ja käyttörajojen perusteella. Vanhentuminen arvioidaan havainnon tai ennusteen julkaisuajasta, ei pelkästä onnistuneen haun ajankohdasta.

**Avoin:** aineistokohtaiset vanhentumisrajat, hakuvirheiden uusintakäytäntö ja tietopakettien tarkka rakenne.

## 8. Tietolähteiden periaatteet

**Vahvistettu**

- MVP käyttää vain FMI:n dokumentoituja avoimia rajapintoja.
- Verkkosivuilta ei poimita tietoja.
- Puuttuvia toissijaisia suureita ei korvata huomaamatta toisella lähteellä tai omalla laskennalla.
- Jos ydintarve, kuten valittujen asemien vedenkorkeusennuste, ei toteudu, rajaus palautetaan käyttäjälle päätettäväksi.
- Havainto, ennuste, sijainti ja ajankohta erotetaan toisistaan.

**Selvitettävä ennen toteutuksen lukitsemista**

- Todelliset rajapintavastaukset, suureet, yksiköt ja ennustejaksot.
- Sopivat kyselyt mökki- ja ulappapisteelle sekä asemille.
- Aineistojen käyttöehdot, käyttörajat ja vaaditut lähdemerkinnät.
- Näennäinen lämpötila ja vedenkorkeuden mahdolliset epävarmuusrajat.

## 9. Tietoturva

**Vahvistettu**

- Vain mökin lähiverkko, ei internetiin avaamista tai porttiohjausta.
- Ei käyttäjätilejä tai kirjautumista.
- Selaimessa vain katselu; asetuksia muokataan palvelimella.
- HTTP hyväksytään luotetun lähiverkon MVP-ratkaisuksi.
- Ylläpito tehdään salatulla SSH-yhteydellä ja tarpeettomat verkkopalvelut suljetaan.

HTTP ei salaa mökkisijaintia tai vastauksia eikä suojaa niitä verkkotasolla muuttamiselta. Tämä on tietoinen rajaus, ei oletus siitä, että lähiverkko olisi automaattisesti turvallinen. Mahdollinen etäkäyttö tai myöhempi anturidata edellyttää tietoturvan uudelleenarviointia.

**Alustava osoiteratkaisu:** paikallinen nimi, esimerkiksi `http://mokinsaa.local`, sekä varayhteys lähiverkon IP-osoitteella. DHCP-varaus ja `.local`-nimen toimivuus eri laitteilla varmistetaan.

## 10. Ylläpito ja julkaiseminen

**Vahvistettu**

- Go-ohjelma ja käyttöliittymä rakennetaan kehityskoneella.
- Sovellus julkaistaan Pi:lle yhtenä versiona käsin käynnistettävällä menettelyllä.
- Edelliseen versioon pitää voida palata.
- Sovellus käynnistyy automaattisesti sähkökatkon jälkeen ja palautuu kaatumisesta.
- Tuettujen pakettilähteiden käyttöjärjestelmän tietoturvapäivitykset asennetaan automaattisesti. Uudelleenkäynnistykset tehdään hallitusti.
- Asetustiedostosta pidetään varmuuskopio muualla kuin Pi:n SD-kortilla.
- Välimuistia ei tarvitse varmuuskopioida.
- Asennuksesta ja palautumisesta laaditaan lyhyt, toistettava ohje.
- Automaattinen julkaisuputki ei kuulu MVP:hen.

**Avoin:** tarkka julkaisumenettely, lokien rajaus, levytilan valvonta ja käyttöjärjestelmän päivityskäytäntö.

## 11. Hyväksymistavoitteet

**Vahvistetut tavoitteet, eivät vielä mitattuja suorituskykylupauksia**

- Yhteenveto tulee käyttökelpoiseksi kahdessa sekunnissa sovitulla testilaitteella mökin lähiverkossa, kun Pi:n välimuisti on valmis.
- Kymmenen yhtäaikaista käyttäjää ei estä normaalia käyttöä.
- Testataan yhden FMI-aineiston häiriö, internetkatko, Pi:n uudelleenkäynnistys ja virheellinen lähdeaineisto.
- Puuttuva tieto ei näy nollana eikä vanha tieto tuoreena.
- Ensimmäinen käynnistys ilman tallennettua aineistoa näyttää selkeän lataus- tai häiriötilan.

**Avoin:** testilaitteet, tarkka mittausmenetelmä ja hyväksyttävän kuormituksen mittarit.

## 12. MVP:n ulkopuolelle jäävät asiat

### Käyttäjän nimeämät jatkokehityskohteet

- Animoitu kartta sadepilville ja salamoille.
- Oman lämpöanturin, esimerkiksi RuuviTagin, datan näyttäminen.
- SQLite-tallennuksen arviointi RuuviTagin yhteydessä.
- Tuulihistoria, jos sille ilmenee tarve.
- Automaattinen julkaisuputki.

### Muut MVP-rajaukset ilman toteutuslupausta myöhemmälle

- Ilmoitukset ja omiin korkeusrajoihin perustuvat korostukset.
- Käyttäjätilit, internetiin avattu palvelu ja etäkäyttö.
- Selaimessa muokattavat paikat tai asetukset.
- Puhelimen itsenäinen offline-käyttö.
- Jatkuvasti päällä olevan yhteisen näytön erikoistila.
- Pitkäaikainen havaintohistoria tai ennusteiden toteutumisen vertailu.

Myöhempi kartta tarvitsee oman tarkastelun karttapohjien, aineistojen, verkkoriippuvuuksien ja suorituskyvyn osalta. RuuviTag tarvitsee oman tarkastelun mittausyhteyden, vastaanottolaitteiston ja historian osalta. Näitä ratkaisuja ei ole päätetty.

## 13. Keskeiset riskit ja riippuvuudet

| Riski | Vaikutus | Sovittu käsittely tai seuraava selvitys |
| --- | --- | --- |
| Pi 2:n laiteversion ja ohjelmistotuen erot | Valittu ajoympäristö tai käyttöjärjestelmä ei sovi laitteelle | Tarkistetaan tarkka laiteversio, arkkitehtuuri ja Go-tuki ennen lukitsemista. |
| FMI-dokumentaatio ja todellinen aineisto poikkeavat toisistaan | Suure tai tavoiteltu ennustejakso puuttuu | Tarkistetaan oikeat rajapintavastaukset. Ydinpuute palautetaan päätettäväksi. |
| Ulappapiste tai sääasema ei kuvaa tavoiteltua aluetta | Käyttäjä tulkitsee väärän alueen olosuhteita | Valitaan pisteet tarkoituksella ja näytetään sijainti/lähde. |
| Ennuste tulkitaan mittaukseksi tai varmaksi tiedoksi | Näkymä antaa harhaanjohtavan käsityksen | Erotetaan havainnot ja ennusteet sekä näytetään dokumentoitu epävarmuus tarvittaessa. |
| Vanhentunut data näyttää tuoreelta | Katko jää käyttäjältä huomaamatta | Aineistokohtaiset aikaleimat, vanhentumisrajat ja häiriötestit. |
| Eri yksiköt, korkeusjärjestelmät tai kertymäjaksot sekoittuvat | Väärä vertailu tai laskettu muutos | N2000, selkeät yksiköt ja dokumentoitu aineiston normalisointi. |
| Sähkökatko, SD-korttivika tai keskeytynyt tallennus | Asetukset tai välimuisti katoavat, palvelu ei palaudu | Hallitut tallennukset, asetusten varmuuskopio ja toistettava asennus. |
| Raskas kaaviototeutus | Käyttöliittymä on hidas puhelimella | Kirjaston arviointi ja suorituskyvyn mittaus oikealla laitteella. |
| HTTP luotetun verkon ulkopuolella | Tietojen paljastuminen tai muuttaminen | Ei internetiin avaamista; käyttörajan muuttuessa uusi tietoturvasuunnitelma. |
| Sovellusta ei avata ennen veden nousua | Rakenteiden suojaamisen tarve jää huomaamatta | Tietoinen rajaus: ei ilmoituksia, käyttäjä tarkistaa ennusteen itse. |
| MVP:n laajuus kasvaa | Ensimmäinen valmis versio viivästyy | Pienet OpenSpec-muutokset ja selkeä jatkokehitysrajaus. |

Pi:n tarkistus ja lähdeaineistojen tarkistus edeltävät lopullista toteutussuunnitelmaa. Aineistojen julkaisurytmi puolestaan edeltää vanhentumisrajojen lopullista määrittelyä.

## 14. Avoimet kysymykset

### Ennen teknisen suunnitelman lukitsemista

1. Mikä Pi 2:n tarkka laiteversio ja suoritinarkkitehtuuri on?
2. Mikä tuettu käyttöjärjestelmä ja Go-versio sopivat sille?
3. Mitkä ovat mökin ja ulappapisteen koordinaatit sekä valitut vedenkorkeus- ja tuulihavaintoasemat?
4. Mitkä FMI-kyselyt tarjoavat tarvittavat tiedot ja kuinka pitkälle?
5. Saadaanko näennäinen lämpötila suoraan, vai tarvitaanko hyväksytty laskentatapa?
6. Onko vedenkorkeuden dokumentoitu epävarmuusaineisto saatavilla?
7. Mitkä käyttöehdot, lähdemerkinnät ja hakurajat koskevat valittuja aineistoja?

### Ennen MVP:n hyväksymistä

8. Mitkä ovat aineistokohtaiset päivitysvälit ja vanhentumisrajat?
9. Mikä kaaviokirjasto täyttää kosketuskäytön, luettavuuden ja suorituskyvyn tarpeet?
10. Millä laitteilla ja menetelmällä hyväksymistavoitteet mitataan?
11. Miten paikallinen nimi, IP-varayhteys, julkaisu, palautus ja lokien hallinta toteutetaan?

Nämä ovat selvitys- ja päätöskohteita, eivät lupa laajentaa MVP:tä. Faktoja selvitetään dokumentaatiosta, laitteelta ja rajapinnoista; merkittävät vaihtoehdot palautetaan käyttäjälle päätettäväksi.

## 15. Ehdotus toteutusvaiheista

**Vahvistettu etenemissuunta:** pieniä OpenSpec-muutoksia yhden suuren toteutuspaketin sijaan. Tarkat muutosrajat ja tehtävät ovat alustavia.

### Vaihe 0: toteutettavuuden tarkistus

- Selvitetään laite- ja käyttöjärjestelmätuki.
- Valitaan sijainnit ja asemat.
- Tarkistetaan FMI:n todelliset vastaukset, kattavuus ja käyttöehdot.
- Ratkaistaan mahdolliset ydintarpeiden esteet käyttäjän kanssa.

Tuloksena varmennetut lähtötiedot, ei vielä koko sovelluksen toteutus.

### Vaihe 1: ensimmäinen kokonainen toiminto

- Yhden aseman vedenkorkeushavainto ja -ennuste.
- Koko ketju: haku → tarkistus → tallennus → datarajapinta → mobiilinäkymä.
- N2000, aikaleimat, muutoksen laskenta ja häiriötilanteet.

Välitavoite ei vielä ole valmis MVP. Toteutusjärjestys ei muuta etusivun sovittua sisältöjärjestystä.

### Vaihe 2: muu MVP

- Muut vedenkorkeusasemat.
- Mökkisää ja lähiaseman tuulihavainnot.
- Ulapan tuuli ja aallokko.
- Seitsemän päivän yleiskuva.
- Yhteenvedot, kaaviot ja tuntitaulukot.

Kokonaisuus voidaan jakaa useampaan pieneen OpenSpec-muutokseen.

### Vaihe 3: käyttöönoton viimeistely

- Oikean Pi:n ja mobiililaitteiden suorituskykytestit.
- Katko-, virhe- ja uudelleenkäynnistystestit.
- Palveluasennus, käsin tehtävä julkaisu ja palautus.
- Varmuuskopiointi- ja ylläpito-ohjeet.

### Työskentelyperiaatteet

- Jokainen muutos yhdistää vaatimukset, hyväksymistestit ja toteutuksen.
- Koodin vastuut ja tiedonkulku käydään ymmärrettävästi läpi.
- Määräpäivää ei ole; työmäärä arvioidaan toteutettavuuden tarkistuksen jälkeen.
- Jos laajuutta pitää pienentää, seitsemän päivän yleiskuva on ehdotettu ensimmäiseksi karsittavaksi, ei tietojen tuoreus tai virheenkäsittely. Karsintaa ei ole päätetty.
- Tämä dokumentti ei käynnistä sovelluskoodausta tai OpenSpec Propose -vaihetta. Niistä sovitaan erikseen.

## 16. Haastattelun aikana tarkistetut lähteet

Lähteet tukevat alustavaa suunnittelua. Ne eivät korvaa rajapintojen tai käytettävän laitteen tarkistusta.

- [Raspberry Pi 3 Model B+: laiteversiot ja tekniset tiedot](https://www.raspberrypi.com/products/raspberry-pi-3-model-b-plus/)
- [Node.js: alustatuki](https://github.com/nodejs/node/blob/main/BUILDING.md) — muuttuva dokumentti; käytettävän version tuki tarkistettava erikseen.
- [FMI: avatut tietoaineistot](https://www.ilmatieteenlaitos.fi/avoin-data-avattavat-aineistot)
- [FMI: meriveden korkeusennuste](https://www.ilmatieteenlaitos.fi/avoin-data-meriveden-korkeusennuste)
- [FMI: aaltoennustedata WAM](https://www.ilmatieteenlaitos.fi/avoin-data-aaltoennustedata-wam)

FMI:n eri kuvaussivuilla havaittiin ristiriita vedenkorkeusennusteen pituudessa. Siksi todellinen kattavuus jätetään avoimeksi, kunnes rajapintavastaus on tarkistettu. WAM-kuvaus puolestaan tukee tavoitetta kolmesta aallokkosuureesta tunneittain, mutta ajantasainen saatavuus on silti varmistettava.
