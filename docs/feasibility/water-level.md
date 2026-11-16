# Vedenkorkeustoiminnon toteutettavuus: Hamina

Tarkistuspäivä: **2026-11-16**. Muutos: `verify-water-level-feasibility`.
Tämä on rajattu dokumentaatioselvitys, ei sovellustoteutus, asennus tai MVP:n hyväksyntä.

## Tila ja näytön tulkinta

**Selvityksen tulos:** Haminan havainto ja N2000-pääennuste saadaan
dokumentoidusta WFS-rajapinnasta. Näytteen ennustekattavuus ylittää 48 h.
N2000-ennusteen yksikkö **cm** on käyttäjän hyväksymä oletus `SeaLevel`-kentän
metatiedon perusteella, ei suoraan varmennettu kenttämetatieto.

**Jatkosuunnittelu: estynyt alkuperäisten ehdottomien ehtojen mukaan**:
käyttöjärjestelmä–Go-ehdotusta ei ole vielä hyväksytty lopullisesti, eikä
N2000-kentän suoraa yksikkövarmennusta ole saatu. Käyttäjän hyväksymä
yksikköpoikkeus on kirjattu §3:een; sitä ei tarvitse hyväksyttää uudelleen,
mutta se on säilytettävä näkyvänä myöhemmässä suunnitelmassa.
Varmennettu tarkoittaa alla mainitusta lähteestä tarkistettua tietoa, ei laiteajon
tai jatkuvan saatavuuden todistusta. Alustava tarkoittaa ehdotusta; avoin tarkoittaa
vielä selvittämätöntä tietoa tai hyväksyntää.

## 1. Laite ja tuleva puhdas asennus

| Tieto | Tulos ja tila | Alkuperä |
| --- | --- | --- |
| Käyttäjän laite | Raspberry Pi 3 Model B+, piirilevyn vuosimerkintä 2017; käyttäjän ilmoittama | Tämän muutoksen keskustelu 2026-11-16: ”malli: rasberry pi 3 b+ 2017” |
| Mallin suoritin | Varmennettu mallin teknisistä tiedoista: Broadcom BCM2837B0, Cortex-A53, ARMv8, 64-bittinen, neljä ydintä, 1,4 GHz | [Raspberry Pi 3 Model B+](https://www.raspberrypi.com/products/raspberry-pi-3-model-b-plus/), tarkistettu 2026-11-16 |
| Mallin muisti | Varmennettu mallin teknisistä tiedoista: 1 GB LPDDR2 | Sama tuotesivu, tarkistettu 2026-11-16 |
| Revisio | Tarkkaa piirilevyrevisiota ei ole ilmoitettu. Vuosimerkintää ei tulkita revisioksi. Tuotesivu määrittelee 3 B+:n ARMv8-malliksi ilman tämän selvityksen ajokohteeseen vaikuttavaa revisioerottelua; lisätietopyyntöä ei nyt tarvita. | Käyttäjän ilmoitus ja sama tuotesivu, tarkistettu 2026-11-16 |
| Laitteen tila | Ei tällä hetkellä käyttökunnossa; laiteajoa tai fyysistä tarkastusta ei ole tehty | Käyttäjän ilmoitus tämän muutoksen keskustelussa 2026-11-16 |
| Nykyinen käyttöjärjestelmä ja bittisyys | Ei tutkita: käyttäjä aikoo alustaa laitteen ennen sääaseman asentamista ja sallii sopivan käyttöjärjestelmän ehdottamisen | Käyttäjän päätös tämän muutoksen keskustelussa, vahvistettu suunnitelmapäivityksessä 2026-11-16 |

Suorittimen kyvykkyys ei todista nykyisen asennuksen arkkitehtuuria.
Laitteen toimintakunto ja oikean ohjelman ajo jäävät myöhempään toteutusmuutokseen.
Tunnuksia, sarjanumeroita tai yksityisiä verkkotietoja ei kerätä.

Projektikuvauksen `docs/PROJECT_BRIEF.md` §2, §13 ja §14 puhuvat edelleen Pi 2:sta.
Tämä on näkyvä ristiriita: tässä selvityksessä käytetään käyttäjän tarkentamaa
Pi 3 Model B+:aa hyväksytyn muutossuunnitelman mukaisesti. Projektikuvausta ei
muuteta tämän rajauksen sisällä.

### Dokumentoitu käyttöjärjestelmä–Go-ehdotus

**Alustava ehdotus, tarkistettu 2026-11-16:** Raspberry Pi OS Lite (64-bit),
Debian 13 Trixie, ilman työpöytää; Go **1.27.4** kehityskoneella.
Käännöskohde: `GOOS=linux GOARCH=arm64 GOARM64=v8.0 CGO_ENABLED=0`.
`GOARM` koskee 32-bittistä `arm`-kohdetta, eikä sitä käytetä tässä.
Ei lisätä `lse`- tai muita valinnaisia ARM-laajennuksia oletukseksi.
CGO:n poistaminen on ehdotus riippumattomalle Go-ohjelmalle; tulevien
riippuvuuksien sopivuus on tarkistettava toteutuksessa.

| Näyttö | Varmennettu dokumentaatiosta 2026-11-16 |
| --- | --- |
| OS:n mallituki ja bittisyys | [Raspberry Pi OS downloads](https://www.raspberrypi.com/software/operating-systems/) listaa 3B+:n 64-bit-version yhteensopivaksi malliksi. Lite on ilman työpöytää; sivulla näkyvä julkaisu 2026-10-06, Debian 13, kernel 6.18, 64-bittinen käyttäjätila. |
| Go-version tuki | [Go Release History ja Release Policy](https://go.dev/doc/devel/release): 1.27.4 julkaistu 2026-11-10, uusin listattu korjausversio. Go-haaraa tuetaan, kunnes kaksi uudempaa pääversiota on julkaistu; nyt 1.27 ja 1.26 ovat tuettuja. Tarkistetaan korjausversio uudelleen toteutuksen alussa. |
| OS/CPU-käännöskohde | [Go: Installing from source, environment](https://go.dev/doc/install/source#environment) listaa `linux/arm64`-yhdistelmän sekä `GOARM64`-asetuksen. [Go Minimum Requirements](https://go.dev/wiki/MinimumRequirements) määrittelee ARMv8.0-perustason ja oletuksen `GOARM64=v8.0`; Linuxille vähintään kernel 3.2 (Go 1.24 ja uudemmat). Ehdotetun OS:n kernel 6.18 ylittää tämän. [Go 1.27:n julkaisutiedot](https://go.dev/doc/go1.27) eivät ilmoita poistavansa tätä kohdetta. |

Tämä on lähteisiin perustuva yhteensopivuuspäätelmä, ei tehty käännös, asennus,
suorituskykymittaus tai laiteajotesti. Pi:n toimintakunto on avoin. Ehdotettu
OS-versio ja Go:n lopullinen valinta odottavat käyttäjän hyväksyntää;
käyttöjärjestelmän ehdottamiseen annettu lupa ei ole niiden hyväksyntä.

## 2. Asema ja aineistot

Käyttäjä valitsi Haminan tämän muutoksen keskustelussa 2026-11-16
(”Vedenkorkeusasema: Hamina”). Valinnan peruste on käyttäjän oma valinta;
asemavalinta on hyväksytty. Asemaa ei vaihdeta ilman uutta hyväksyntää.

**Varmennettu 2026-11-16:** [FMI:n asematiedot](https://www.ilmatieteenlaitos.fi/havaintoasemat)
ja dokumentoitu WFS-asemakysely tunnistavat mareografin nimellä **Hamina
Pitäjänsaari**, FMISID **134254**, sijaintitunniste (geoid) **-10022823**,
verkko **133 / Mareografiasema**. Julkiset asemakoordinaatit ovat
60.562767, 27.179200 (EPSG:4258); ne eivät ole mökkikoordinaatteja.
Alkuperäinen asemavastaus: [hamina-station.xml](evidence/hamina-station.xml).

Asiakirjalähde: [FMI:n latauspalvelun pikaohje](https://www.ilmatieteenlaitos.fi/latauspalvelun-pikaohje),
asematietojen haku `fmi::ef::stations` ja kohdistus `fmisid=134254`.
[Havaintokyselyn kuvaus](https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=DescribeStoredQueries&storedquery_id=fmi::observations::mareograph::instant::timevaluepair)
ja [ennustekyselyn kuvaus](https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=DescribeStoredQueries&storedquery_id=fmi::forecast::sealevel::point::timevaluepair)
tukevat molemmat `fmisid`-rajausta. Molempiin käytetään **134254**:ää;
vastausten asematunnisteiden vastaavuus tarkistetaan näytehauissa.
Kaikki tämän kappaleen lähteet tarkistettu 2026-11-16.

### Dokumentoidut kyselyt ja käyttöehdot

Tarkistettu **2026-11-16 ennen aineistohakuja**:

- Rajapinta `https://opendata.fmi.fi/wfs`, WFS **2.0.0**, `request=GetFeature`.
  [Pikaohje](https://www.ilmatieteenlaitos.fi/latauspalvelun-pikaohje) ja
  [WFS-käyttöohje](https://en.ilmatieteenlaitos.fi/open-data-manual-fmi-wfs-services)
  dokumentoivat tallennetut kyselyt ja `DescribeStoredQueries`-metatiedon.
- Havainto: `storedquery_id=fmi::observations::mareograph::instant::timevaluepair`,
  `fmisid=134254`, `timezone=UTC`, `timestep=1`.
  Varsinainen näyte rajataan aikavälille **2026-11-16T17:00:00Z–18:00:00Z**.
  Kuvaus dokumentoi yhden minuutin askeleen ja `parameters`-rajauksen.
  Suureen tarkka nimi selvitetään ensin samalla asemalla vain yhdestä
  ajankohdasta **2026-11-16T18:00:00Z**, ilman `parameters`-rajausta,
  pikaohjeen mukaisesti; sen jälkeen tallennetaan vain N2000-havaintosarja.
- Ennuste: `storedquery_id=fmi::forecast::sealevel::point::timevaluepair`,
  `fmisid=134254`, `timezone=UTC`, `timestep=60`,
  **2026-11-16T18:00:00Z–2026-11-21T18:00:00Z**.
  Viiden vuorokauden yläraja tarkistaa
  [aineistoluettelon](https://www.ilmatieteenlaitos.fi/avoin-data-avattavat-aineistot)
  5 vrk -kuvauksen suhteessa
  [ennustekuvauksen](https://www.ilmatieteenlaitos.fi/avoin-data-meriveden-korkeusennuste)
  noin kahteen vuorokauteen. Ei oleteta pyydetyn aikajakson toteutuvan.
  Ennusteen `DescribeStoredQueries` ei dokumentoi `parameters`-rajausta,
  joten sitä ei käytetä: yksi asema ja vain dokumentoidun kyselyn oletussuureet.
- Kaikki rajaukset kuuluvat yllä linkitettyjen kyselykuvausten parametreihin.
  Tarkat suoritetut URL:t, hakuajat ja onnistumiset kirjataan manifestiin.
  Ei jatkuvaa hakua, ei kaikkien asemien dataa.

[FMI:n lisenssisivu](https://en.ilmatieteenlaitos.fi/open-data-licence)
määrää avoimen datan lisenssiksi **CC BY 4.0**.
[Lisenssi](https://creativecommons.org/licenses/by/4.0/) sallii kopioinnin,
jakamisen ja muokkaamisen, mutta vaatii asianmukaisen lähdemerkinnän,
lisenssilinkin ja muutosten ilmoittamisen; FMI:n hyväksyntää ei saa antaa ymmärtää.
Näytteet tallennetaan muuttamattomina, tulkinnat erilliseen raporttiin.
Lähdemerkintä: **”Lähde: Ilmatieteen laitos, avoin data, CC BY 4.0”**,
linkit alkuperäiseen kyselyyn ja lisenssiin. Myöhempi yksikkömuunnos ja
lasketut yhteenvedot on erotettava alkuperäisarvoista.

WFS-käyttöohjeen dokumentoidut rajat: **20 000 latauspyyntöä/vrk**,
katselupalvelu 10 000/vrk ja molemmille yhteensä **600 pyyntöä/5 min**.
Tämän selvityksen muutamat peräkkäiset metatieto- ja näytehaut alittavat rajat.
Projektikuvauksen alustavat yhden aseman välit (havainto 10 min, ennuste 30 min)
merkitsisivät 144 + 48 = 192 aineistopyyntöä/vrk; tämä on laskelma, ei
tässä toteutettu ajastus. Muut sovellukset ja saman käyttäjän muut haut on
huomioitava myöhemmin. Tarkastetuissa kyselykuvauksissa ei ilmoiteta
erillistä aineistokohtaista aikarajaa tai saatavuustakuuta; mahdolliset
palvelinrajoitukset ja virheet käsitellään näkyvinä rajoitteina, ei ohiteta.

## 3. Tallennetut vastaukset ja tarkistusmatriisi

Alkuperäiset vastaukset ja niiden täsmälliset URL:t, UTC-hakuajat,
SHA-256-tiivisteet sekä dokumentaatiolähteet ovat
[manifestissa](evidence/manifest.md). Havainto- ja ennustehaku onnistuivat
HTTP- ja XML-tasolla. Asemat vastaavat toisiaan: havaintojen FMISID 134254
ja ennusteen geoid -10022823 yhdistyvät tallennetun asemametatiedon kautta.
Ennusteen vastausmuoto on XML/GML + WaterML 2.0, samoin havainnon aikasarja.

Alla olevat näytehavainnot tarkistettu **2026-11-16** tallennetuista vastauksista.
Yksikköä koskeva tarkistus erottaa lähteestä varmennetun tiedon käyttäjän
hyväksymästä oletuksesta: `SeaLevelN2000`-kentän yksikkömetatieto puuttuu
edelleen. Käyttäjä hyväksyi 2026-11-16 jatkoselvityksen vararatkaisuksi saman
yksikön kuin `SeaLevel`-metatiedossa. Tämä päätös ei tee oletuksesta
FMI:n vahvistamaa tietoa.

| Tarkistus | Havainto | Ennuste |
| --- | --- | --- |
| Suure | `WLEVN2K_PT1S_INSTANT`, laadunvalvottu vedenkorkeus | `SeaLevelN2000`; lisäksi `SeaLevel`, jota ei käytetä N2000-arvona |
| N2000-osoitus | Varmennettu: `observation-property.xml` sanoo N2000 | Parametrin nimi viittaa N2000:een, ja ennustekuvaus sanoo aineiston olevan saatavilla N2000:ssa. Parametrin oma metatieto on tyhjä. |
| Yksikkö | Varmennettu `mm` metatiedosta; myöhempi cm-esitys vaatii jaon kymmenellä | **Käyttäjän hyväksymä oletus: cm**, kuten saman vastauksen `SeaLevel`-kentällä, jonka yksikkö on varmennettu `sealevel-property.xml`:stä. N2000-kentän oma `forecast-property.xml` on edelleen tyhjä; sen yksikköä ei väitetä varmennetuksi. |
| Aikaleimat | `wml2:time` mittaushetki UTC:ssa; `phenomenonTime` 17:00–18:00 UTC | `wml2:time` ennusteen voimassaolohetki UTC:ssa. `phenomenonTime` ilmoittaa pyydetyn jakson, ei äärellisten arvojen kattavuutta. |
| Ensimmäinen ja viimeinen arvo | 2026-11-16T17:00:00Z–18:00:00Z, 61 äärellistä arvoa | 2026-11-16T18:00:00Z–2026-11-21T11:00:00Z, 114 äärellistä arvoa kummassakin sarjassa |
| Viimeisin havainto | Rajatun näytteen viimeinen: 2026-11-16T18:00:00Z, 321 mm = 32,1 cm N2000. Ei väitettä koko palvelun viimeisimmästä havainnosta hakuhetkellä. | Ei havainto |
| Askel ja aukot | 60 sekuntia, ei puuttuvia minuutteja tai arvoja näytteessä | 3600 sekuntia, ei sisäisiä aukkoja äärellisessä jaksossa. 121 aikaleimaa/sarja; 2026-11-21T12:00:00Z–18:00:00Z seitsemän `NaN`-arvoa/sarja. |
| Todellinen kattavuus vs. 48 h | Näyte on tarkoituksella vain tunti | Äärellinen jakso 113 h, 48 h tavoite täyttyy ajallisesti myös noin 18:08 UTC hakuhetkestä; ei lupausta tulevien vastausten pituudesta. |
| Puuttuva arvo | Ei havaittu puuttuvia arvoja tässä näytteessä; puuttumisen esitystapaa ei testattu havaintosarjalle | `wml2:value` sisältää tekstin `NaN`; ei nolla. Pyydettyä loppujaksoa ei saa esittää ennusteena. |
| Tuloksen aika | `om:resultTime` 2026-11-16T18:00:00Z; ei erikseen varmennettu julkaisuaika | `om:resultTime` 2026-11-16T17:00:11Z. N2000-sarja viittaa ensimmäisen sarjan aikaan sisäisellä `xlink:href`-viitteellä. |
| Mallin viiteaika | Ei ennustemallia | `analysisTime` 2026-11-16T17:00:00Z. Ei sama asia kuin haku- tai todistettu julkaisuaika. |

Ajat ovat UTC (`Z`, lisäksi kyselyssä `timezone=UTC`).
[FMI:n FAQ](https://en.ilmatieteenlaitos.fi/open-data-questions) dokumentoi
avoimen datan UTC-ajan. Sijainnin `Europe/Helsinki` ei muuta näiden aikojen
tulkintaa. `FeatureCollection.timeStamp` on vastauksen aikaleima, ei ennusteen
julkaisuajaksi varmennettu tieto. [Prosessimetatieto](https://xml.fmi.fi/inspire/process/sealevel)
kuvaa `analysisTime`-ajan malliajojen nimelliseksi viiteajaksi, mutta sisältää
myös vanhan OAAS-HIRLAM-nimen ja `TBA`-kuvauksen; siitä ei päätellä nykyistä
malliketjua tai julkaisurytmiä. Lähteet tarkistettu 2026-11-16.

### Yksikön jatkoselvitys ja käyttäjän hyväksymä oletus

**Varmennettu 2026-11-16:** [SeaLevel-metatieto](https://opendata.fmi.fi/meta?observableProperty=forecast&param=SeaLevel&language=eng)
antaa yksiköksi **cm**. Muuttamaton vastaus on tallennettu
[sealevel-property.xml](evidence/sealevel-property.xml)-tiedostoon.
N2000-kentän metatieto tarkistettiin uudelleen: se on edelleen tyhjä ja
tavuilleen sama kuin alkuperäinen tallennettu vastaus (manifestin ajat).

Täydentävästä virallisesta aineistosta tarkistettiin
[FMI:n ennustekuvaus](https://www.ilmatieteenlaitos.fi/avoin-data-meriveden-korkeusennuste),
[merisääpalvelun selitys](https://en.ilmatieteenlaitos.fi/marine-weather-online)
sekä FMI:n oman SmartMet-kirjaston
[parametriluettelo](https://github.com/fmidev/smartmet-library-newbase/blob/f87b5d3b36bb056f1c2c705836b83273410349b7/newbase/NFmiParameterName.h)
ja [muutoshistoria](https://github.com/fmidev/smartmet-library-newbase/blob/f87b5d3b36bb056f1c2c705836b83273410349b7/smartmet-library-newbase.spec).
Kirjasto tunnistaa `SeaLevelN2000`-parametrin (numero 159) ja historia
kertoo N2000:n tarkoittavan vertailutasoa, mutta ei määritä kentän yksikköä.
Verkkopalvelun selitys käyttää cm-yksikköä molemmille korkeusjärjestelmille,
mutta se ei yksin määritä WFS-kentän yksikköä. WFS:n julkisesta
[esimerkkikonfiguraatiosta](https://github.com/fmidev/smartmet-plugin-wfs/tree/3e790e36dc5081574991000c42b01404e57611bb/cnf/opendata_stored_queries_github)
löytyi vanhoja OAAS-kyselyitä; niitä ei pidetä nykyisen palvelun
konfiguraation tai yksikön todisteena. Tarkistettu 2026-11-16.

**Käyttäjän päätös 2026-11-16:** ”Jos kuitenkaan vastausta ei löydy
käyttäisin samaa mittayksikköä” kuin `SeaLevel`-metatiedossa.
Koska suoraa vahvistusta ei löytynyt tarkistetuista lähteistä,
`SeaLevelN2000` tulkitaan jatkosuunnittelussa **alustavasti senttimetreinä**.
Tämä on nimenomainen käyttäjän hyväksymä poikkeus alkuperäiseen yksikön
varmennustavoitteeseen, ei tehtävässä 3.2 saatu suora yksikkötodiste.
Tarkistusmatriisin selvitys on tehty; lähdepuute jää näkyväksi.
Suureeksi pysyy `SeaLevelN2000`: `SeaLevel`-arvoja ei korvata N2000-arvoiksi
eikä tehdä omaa korkeusjärjestelmämuunnosta.

Jos tätä oletusta käytetään myöhemmässä toteutuksessa, havainto normalisoidaan
millimetreistä senttimetreiksi, ennusteen arvo säilytetään sellaisenaan ja
yksikköoletus sekä sen riski on tuotava toteutussuunnitelmassa näkyviin.
Tässä selvityksessä ei lasketa käyttöön tarkoitettuja nousu- ja
laskuyhteenvetoja eikä toteuteta muunnosta. FMI:n ristiriitainen vahvistus
edellyttää tulkinnan ja siihen perustuvien laskelmien korjaamista.
[FMI:n tukipyyntö](https://en.ilmatieteenlaitos.fi/open-data-support-request)
on edelleen mahdollinen tapa saada suora vahvistus; yhteydenottoa ei ole tehty.

### Pääennuste, epävarmuus ja julkaisurytmi

Tarkistettu **2026-11-16** yllä mainitusta ennustevastauksesta,
[ennustekuvauksesta](https://www.ilmatieteenlaitos.fi/avoin-data-meriveden-korkeusennuste),
[aineistoluettelosta](https://www.ilmatieteenlaitos.fi/avoin-data-avattavat-aineistot)
ja [merisääpalvelun selityksestä](https://en.ilmatieteenlaitos.fi/marine-weather-online):

- **Pääennuste:** valitun WFS-kyselyn N2000-ennustesarja on
  `SeaLevelN2000`. FMI kuvaa vedenkorkeusennusteen automaattiseksi,
  havainnoilla korjatuksi merimallitulokseksi. Sarjaa ei tämän perusteella
  nimetä mediaaniksi, keskiarvoksi tai todennäköisyysrajaksi. Korjaus
  on lähteen menetelmä; sovellus ei saa itse siirtää käyrää havaintotasolle.
- **Epävarmuusrajat:** tallennettu vastaus sisältää vain `SeaLevel`- ja
  `SeaLevelN2000`-sarjat. Ne ovat eri korkeusjärjestelmien arvoja, eivät
  ala- ja yläraja. [ListStoredQueries](https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=ListStoredQueries)
  tarkistettiin 2026-11-16T18:47:27Z: vedenkorkeudelle löytyi kolme
  vastausmuotoa (`simple`, `timevaluepair`, `multipointcoverage`), ei
  erillistä vedenkorkeuden todennäköisyyskyselyä.
  Tarkistetuista avoimen datan kuvauksista ei löytynyt dokumentoitua
  epävarmuusrajojen koneellista kyselyä. **Saatavuus avoimesta rajapinnasta
  jää avoimeksi**, ei väitettä ettei FMI tuottaisi epävarmuusennusteita.
- **Todennäköisyysennusteen merkitys:** verkkopalvelun selitys kertoo
  erillisestä 10 vuorokauden parviennusteesta, mediaanista sekä 50 % ja
  80 % ennusteväleistä. Nämä eivät ole absoluuttisia minimi- tai
  maksimivedenkorkeuksia. Verkkopalvelua ei käytetä datalähteenä,
  eikä sen käyrästä tai sisäisestä rajapinnasta poimita arvoja.
  Jos dokumentoitu rajapinta myöhemmin löytyy, se tarvitsee erillisen
  tarkistuksen. Ehdollisten epävarmuusrajojen puuttuminen ei tämän
  muutoksen designin mukaan yksin estä ensimmäisen toiminnon suunnittelua.
- **Dokumentoitu julkaisurytmi:** ennustekuvaus (sivulla päiväys 2.9.2022)
  sanoo noin **04, 10, 16 ja 22 UTC**, aineistoluettelo **4 kertaa/vrk**.
  Verkkopalvelun selitys sanoo pääennusteelle neljä päivitystä/vrk ja
  erilliselle todennäköisyysennusteelle kaksi. Näitä rytmejä ei sekoiteta.
  Havaintokyselyn kuvaus antaa yhden minuutin aika-askeleen;
  se ei takaa minuutin julkaisulatenssia. Verkkopalvelun havaintojen
  10 minuutin päivitys ei sekään ole WFS:n julkaisulupaus.
- **Ristiriidat jäävät näkyviin:** 113 h:n äärellinen näyte on pidempi
  kuin erillisen kuvauksen noin 2 vrk ja sopii paremmin 5 vrk:n
  aineistoluetteloon. Näytteen `analysisTime=17:00Z` ja
  `resultTime=17:00:11Z` eivät yksin vahvista vanhan sivun kellonaikoja
  eivätkä todista uutta julkaisurytmiä. Vanha prosessikuvaus mainitsee
  OAAS-HIRLAMin ja `TBA`:n. Mikään näistä ristiriidoista ei oikeuta
  keksimään täsmällistä aikataulua tai saatavuuslupausta.
- **Vaikutus tuoreuden arviointiin:** myöhemmän toteutuksen on säilytettävä
  havaintohetki, ennusteen viite-/tulosaika ja viimeinen äärellinen
  voimassaolohetki erillään hakuajasta. Sisäiset `xlink:href`-aikaviitteet
  pitää ratkaista. Onnistunut uusintahaku ei nuorenna vanhaa malliajoa.
  Vanhentumisraja ja viiteajan/tulosajan valinta ikälaskentaan tarvitsevat
  myöhemmän nimenomaisen määrittelyn; tässä ei aseteta numerorajaa.

## 4. Toistamisohjeet

### Tallennetun näytön tarkistus ilman verkkohakuja

Aja projektin juuresta Python 3:lla (vain standardikirjasto). Tämä ei
asenna riippuvuuksia, muuta tiedostoja tai ota yhteyttä Pi:hin/FMI:hin.
Tiiviste kertoo vastaako tiedosto manifestiin kirjattua runkoa; se ei ole
FMI:n digitaalinen allekirjoitus tai lupaus vastauksen sisällön oikeellisuudesta.

```bash
python3 - <<'PY'
from pathlib import Path
import hashlib, math, re, xml.etree.ElementTree as ET
from datetime import datetime

base = Path("docs/feasibility/evidence")
manifest = (base / "manifest.md").read_text()
sections = manifest.split("### `")[1:]
assert len(sections) == 6
for section in sections:
    name = section.split("`", 1)[0]
    expected = re.search(r"SHA-256: `([0-9a-f]{64})`", section)[1]
    raw = (base / name).read_bytes()
    assert hashlib.sha256(raw).hexdigest() == expected, name
    ET.fromstring(raw)
    print("XML + SHA-256 OK:", name)

ns = {
    "wfs": "http://www.opengis.net/wfs/2.0",
    "gml": "http://www.opengis.net/gml/3.2",
    "om": "http://www.opengis.net/om/2.0",
    "wml": "http://www.opengis.net/waterml/2.0",
}
for name, count, finite, step, station in [
    ("hamina-observation.xml", 61, 61, 60, "134254"),
    ("hamina-forecast.xml", 121, 114, 3600, "-10022823"),
]:
    root = ET.parse(base / name).getroot()
    members = root.findall("wfs:member", ns)
    assert len(members) == (1 if "observation" in name else 2)
    for member in members:
        assert member.find(".//gml:identifier", ns).text == station
        pairs = member.findall(".//wml:MeasurementTVP", ns)
        times = [datetime.fromisoformat(p.find("wml:time", ns).text.replace("Z", "+00:00")) for p in pairs]
        values = [float(p.find("wml:value", ns).text) for p in pairs]
        assert len(pairs) == count
        assert all(math.isfinite(v) for v in values[:finite])
        assert all(math.isnan(v) for v in values[finite:])
        assert all((b - a).total_seconds() == step for a, b in zip(times, times[1:]))
    print("Sarjat, asemat, askeleet ja puuttuvat arvot OK:", name)
PY
```

Raaka-XML:n rivien lopun välilyöntejä tai sisennyksiä ei siivota:
muotoilukin kuuluu tallennettuun alkuperäisrunkoon. `manifest.md` linkittää
jokaiseen kuuteen XML-tiedostoon ja raportti takaisin manifestiin.
Matriisin arvojen määrä ja askel on laskettava kustakin sarjasta, ei koko
tiedoston kaikkien aikapisteiden yhteismäärästä. Aikojen sisäiset
`xlink:href`-viitteet on huomioitava käsin luettaessakin.

### Valinnainen rajattu uusintahaku

1. Tarkista ennen uusintahakua §2:n lähteistä ajantasaiset käyttöehdot,
   hakurajat ja `DescribeStoredQueries`-parametrit.
2. Käytä [manifestin](evidence/manifest.md) `hamina-observation.xml`- ja
   `hamina-forecast.xml`-kohdissa olevia **täsmällisiä suoritettuja URL:eja**.
   Ne vastaavat tallennettujen sarjojen asemaa, suureita, UTC-aikoja ja askelia.
   Yksi pyyntö kumpaankin riittää; ei silmukkaa, automaattisia uusintoja
   tai kaikkien asemien hakua.
3. Tallenna uudet vastaukset erillisiin tiedostoihin, älä ylikirjoita näitä
   näytteitä. Esimerkiksi `curl --fail --show-error --max-time 60 --output
   /tmp/opencode/hamina-observation-uusi.xml 'MANIFESTIN_TARKKA_URL'`,
   kun olet ensin varmistanut ettei kohdetiedosto ole jo olemassa.
   Kirjaa todellinen UTC-hakuaika ja HTTP-tulos. Tarkista lisäksi XML:
   HTTP 200 ei takaa sarjojen tai metatiedon olemassaoloa.
4. Historiallinen havaintohaku saattaa edelleen onnistua, mutta ennusteiden
   historiaa tai alkuperäisen malliajon palautumista ei luvata.
   Nykyhetken näytteelle vaihda aikarajaus käsin: havainto enintään tunti
   menneisyyttä ja ennuste enintään viisi vuorokautta, sama asema ja samat
   suureet/askeleet. Kirjaa muutetut URL:t, uudet hakuajat ja uudet tulokset
   erikseen, älä käytä tämän manifestin aikoja uusille hauille.
5. Uuden haun arvot, viimeisin havainto ja malliajo vaihtuvat ajan mukana.
   Se ei ole toistettavuusvirhe: alkuperäinen johtopäätös tarkistetaan
   säilytetystä näytteestä; uusi haku kuvaa uutta tilannetta.
   `NaN`, tyhjä sarja tai virhevastaus arvioidaan dokumentaation ja
   aikarajauksen kanssa, ei päätellä automaattisesti ettei aineistoa ole.

Ohjeiden URL-parametrit ja paikallisten linkkien kohteet tarkistettu
manifestia sekä tallennettuja vastauksia vasten 2026-11-16.

## 5. Loppuarvio, päätökset ja rajauksen ulkopuoliset työt

Arvioitu **2026-11-16** muutoksen `design.md`:n §4 ehtoja vasten.
Dokumentaatioselvityksen valmistuminen ei tarkoita, että kaikki
jatkosuunnitteluehdot olisi varmennettu tai teknologia lukittu.

| Jatkosuunnitteluehto | Näyttö ja tila 2026-11-16 | Päätös tai jatkotarve |
| --- | --- | --- |
| Pi:n lähtötiedot ja tuettu ajokohde | Malli käyttäjän ilmoittama; mallin ARMv8-tuki, OS:n yhteensopivuus ja Go:n kohde dokumentoidusti tarkistettu (§1) | **Alustava:** Raspberry Pi OS Lite 64-bit / Debian 13 Trixie + Go 1.27.4 / `linux/arm64`, `GOARM64=v8.0`. Käyttäjän lopullinen hyväksyntä puuttuu. |
| Käyttäjän hyväksymä asema, havainto ja pääennuste | **Varmennettu:** käyttäjän valitsema Hamina Pitäjänsaari, FMISID 134254 ja geoid -10022823; aidot vastaukset (§2–3, manifesti) | Ei asemavaihdosta tai muuta datalähdettä. Yksi onnistunut vastaus ei takaa jatkuvaa saatavuutta. |
| Yksiköt ja N2000 | Havainnon N2000 ja mm varmennettu, ennusteen N2000-parametri tunnistettu; suora yksikkömetatieto puuttuu | **Käyttäjän hyväksymä oletus:** ennuste cm kuten `SeaLevel`. Alkuperäinen suoran yksikkövarmennuksen tavoite ei täyty; tämä on sallittu riskipäätös, ei varmennettu lähdefakta. |
| Aikojen merkitys ja ennustejakso | **Varmennettu näytteestä:** UTC, havaintohetki, ennusteen viite-/tulosaika, 113 h äärellinen jakso ja seitsemän puuttuvaa lopputuntia (§3) | Vanhentumisraja sekä ajastus määritellään myöhemmin. Julkaisuajaksi ei oleteta hakuaikaa. |
| Käyttöehdot, hakurajat ja lähdemerkintä | **Dokumentoitu:** CC BY 4.0, 20 000 latauspyyntöä/vrk, yhteinen 600/5 min, lähdemerkintä ja muutosten ilmoittaminen (§2) | Huomioitava toteutuksessa keskitetysti; ei tässä tehtyä ajastusta tai käyttöliittymää. |
| Ydinpuutteiden palauttaminen käyttäjälle | Yksikköpuute esitettiin käyttäjälle, joka hyväksyi `SeaLevel`-yksikön vararatkaisuksi; muita korvauksia ei tehty (§3) | Ei piilotettua N2000-muunnosta, käyrän siirtoa tai puuttuvan ennusteen jatkamista. |
| Epävarmuusrajat | **Avoin:** tarkistetusta WFS-kyselystä ei saatu rajoja eikä löytynyt dokumentoitua erillistä kyselyä (§3) | Ei yksin jatkamisen este: projektikuvaus edellyttää näyttämistä vain jos dokumentoidut rajat ovat saatavilla. |

### Käyttäjälle palautettavat päätökset ja seuraukset

1. **Hyväksyttävä vielä:** yllä oleva käyttöjärjestelmä–Go-yhdistelmä tulevan
   toteutuksen suunnittelukohteeksi, tai valittava muu erikseen selvitettävä
   yhdistelmä. Käyttäjän aiempi lupa ehdottaa sopivaa käyttöjärjestelmää
   ei ole automaattinen tietyn version tai Go:n hyväksyntä.
2. **Jo hyväksytty:** N2000-ennusteelle cm-oletus, jos suoraa vahvistusta
   ei löydy. Oletus ja sen lähdepuute seuraavat mahdolliseen myöhempään
   toteutussuunnitelmaan; niitä ei muuteta varmennetuksi yksikkötiedoksi.
   Mahdollinen FMI:n täydentävä vahvistus tai ristiriitainen tieto
   kirjataan uutena näyttönä. Tässä ei lähetetty ulkopuolista tukipyyntöä.
3. **Myöhemmin määriteltävä, ei nyt päätetty:** havaintojen ja ennusteiden
   täsmälliset vanhentumisrajat, aikaviitteiden käyttö, virheiden käsittely,
   puuttuvien arvojen näkymä sekä yksikköoletuksen näkyvä dokumentointi.
   Tämä selvitys ei hyväksy sovelluksen laskentaa, käyttöliittymää tai
   turvallisuuteen liittyviä tulkintoja.

### Vaiheen 0 jäljellä olevat työt ja myöhemmät toteutustyöt

- Mökkipaikan ja ulappapisteen valinta sekä niiden koordinaatit;
  yksityisiä sijainteja ei ole kerätty.
- Muut vedenkorkeusasemat, sopiva tuulihavaintoasema sekä mökki- ja
  ulappapisteiden sää-, tuuli-, sade- ja aallokkoaineistot.
- Näennäisen lämpötilan ja mahdollisen sateen todennäköisyyden saatavuus,
  seitsemän päivän sääyleiskuva ja muiden aineistojen käyttöehdot.
- Pi:n toimintakunnon tarkistus, puhdas asennus, oikean ohjelman käännös-
  ja laiteajotesti, suorituskyky sekä mobiililaitteiden testaus.
- Palvelin, käyttöliittymä, ajastus, välimuisti, palveluasennus, tietoturvan
  toteutus, julkaisu, palautuminen ja häiriötestit.
- Projektikuvauksen vanhan Pi 2 -merkinnän korjaaminen myöhemmin erillisellä
  hyväksynnällä. Tätä selvitystä ohjaa käyttäjän tarkentama Pi 3 B+.

Nämä eivät muutu tehdyiksi tämän muutoksen tehtävien valmistuessa.
Seuraava sovellustoiminto suunnitellaan omana OpenSpec-muutoksenaan vasta
käyttäjän uudesta pyynnöstä ja yllä olevat päätökset huomioiden.

## 6. Selvityskokonaisuuden tarkistus

Tarkistettu **2026-11-16**:

- Kaikki kuusi lähde-XML:ää jäsentyvät, ja niiden SHA-256-tiivisteet
  vastaavat manifestia. Aiemmin tallennettuja viittä lähdenäytettä
  ei muutettu; jatkoselvitys lisäsi vain `SeaLevel`-suuremetatiedon.
- Raportin oma offline-toistamisohje suoritettiin onnistuneesti.
  Lisäksi tarkistettiin manifestin kyselyparametrien ja UTC-hakuaikojen
  vastaavuus, sarjojen asematunnisteet, aikarajat, yhden minuutin/
  tunnin askeleet, puuttuvat arvot sekä 113 h:n äärellinen kattavuus.
- Kaikki raportin ja manifestin paikalliset linkit löytyvät. Keskeisissä
  johtopäätöksissä on lähde, tarkistuspäivä ja tieto siitä, mikä on
  varmennettu, alustava, käyttäjän hyväksymä oletus tai avoin.
- Ristiriidat (Pi 2 / Pi 3 B+, 2 vrk / 5 vrk, viiteaika /
  dokumentoitu päivitysrytmi ja vanha prosessikuvaus) sekä puuttuva
  N2000-yksikkömetatieto säilyvät näkyvinä.
- Muutos sisältää vain dokumentaatiota, julkisia alkuperäisvastauksia
  ja tehtävien seurantaa. Ei sovelluskoodia, uusia ajonaikaisia
  riippuvuuksia, salaisuuksia tai yksityisiä sijainti-/verkkotietoja.
  Dokumentin Python-ohje on paikallisen näytön tarkistus, ei sovellus.
- Pi:llä ei ole tehty asennusta, alustusta, käännöstä tai ajotestiä.
  Koko vaihetta 0 tai MVP:tä ei väitetä valmiiksi tai hyväksytyksi.
  Selvityksen valmistuminen sisältää §3:n käyttäjän nimenomaisesti
  hyväksymän yksikköoletuksen, ei väitettä tämän lähdepuutteen ratkeamisesta.
