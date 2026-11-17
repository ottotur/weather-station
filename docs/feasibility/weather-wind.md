# Sää- ja tuulilähteiden toteutettavuusselvitys

Muutos: `verify-weather-wind-feasibility`. Aloitettu 2026-11-17.

## Tila ja rajaus

Selvitys on toteutettu 2026-11-17. Testipisteiden kuusi sääennustesuuretta /
kolme ulappatuulisuuretta ja hyväksytyn Rankin kolme tuulihavaintoa saatiin
rajatuilla kyselyillä. Seuraavaa sovellustoimintoa tai lopullisen sijainnin
soveltuvuutta ei ole hyväksytty. Tämä muutos tuottaa
vain raportin ja rajatut lähdenäytteet: ei sovelluskoodia, asennuksia,
verkkopalvelun avaamista, aallokkoa tai seitsemän päivän yleiskuvaa.
Rankin havaintonäyte haettiin vasta käyttäjän asemavalinnan hyväksynnän jälkeen.
Jatkosuunnittelu on lähdenäytön perusteella mahdollista §9:n rajauksella;
puutteet ja ehdotukset eivät ole uusia hyväksyttyjä sovellusvaatimuksia.

## 1. Käyttäjän valitsemat testikohteet

Koordinaatit ovat järjestyksessä leveysaste, pituusaste.

| Kohde | Pyydetty piste | Valinnan alkuperä |
| --- | --- | --- |
| Testimökkipiste | 60.434318, 26.687825 | Käyttäjän suunnitteluvaiheessa antama piste; kirjattu muutoksen proposal.md:ssä ja design.md:n Context-kohdassa |
| Testiulappapiste | 60.414616, 26.697395 | Käyttäjän myöhempi nimenomainen valinta; kirjattu samoihin suunnitteluasiakirjoihin |

Kumpikaan ei ole lopullinen sijainti. Testiulappapiste on jo hyväksytty tätä
selvitystä varten eikä edellytä uutta sijaintivalintaa. Pyydetty ja palautettu
mallipiste erotetaan näytteiden yhteydessä. Mahdollinen sijaintieste palautetaan
käyttäjälle; pistettä ei vaihdeta hiljaisesti. Havainto, piste-ennuste ja
edustavuusarvio eivät ole sama asia.

Alkuperäviitteet:
[proposal](../../openspec/changes/verify-weather-wind-feasibility/proposal.md),
[design](../../openspec/changes/verify-weather-wind-feasibility/design.md).

## 2. Tuulihavaintoaseman valinta — Rankki hyväksytty

Tarkistettu 2026-11-17 FMI:n
[asemataulukosta ja ympäristökuvauksista](https://www.ilmatieteenlaitos.fi/havaintoasemat)
Playwrightilla (taulukon haku `Kotka`, kummankin aseman infopainike).
Taulukon koordinaatit pyöristyvät yhteen desimaaliin; alla käytetään virallisen
WFS-asemametatiedon tarkempia EPSG:4258 / `Lat Long` -koordinaatteja.
FMI kuvaa EUREF-FIN-koordinaatit käytännössä WGS84:ää vastaaviksi.

| Ehdokas | FMISID / WMO / geoid | Leveysaste, pituusaste | Etäisyys mökkiin / ulapalle |
| --- | --- | --- | --- |
| **Suositus: Kotka Rankki** | 101030 / 02976 / -16000076 | 60.375377, 26.958926 | 16,27 / 15,01 km |
| Vaihtoehto: Kotka Haapasaari | 101042 / 02967 / -16000147 | 60.286758, 27.184825 | 31,88 / 30,35 km |

Etäisyydet ovat itse laskettuja linnuntie-etäisyyksiä haversine-kaavalla
(pallon säde 6371,0088 km), eivät FMI:n edustavuusarvioita.

**Rankin perustelu:** FMI:n mukaan etelän suunta on lähes avomerta, muualla
hajanaista saaristoa ja luoteessa manner. Tuulimittari on noin 15 m tornin
päällä noin 5 m mastossa; saarella on 10–20 m metsää. Se on näistä kahdesta
lähempänä testipisteitä ja tarjoaa perustellun meri-/saaristotuulen
vertailuhavainnon. Lähimmän aseman ei väitetä automaattisesti olevan edustavin,
eikä tässä ole kartoitettu kaikkia alueen asemia.

**Haapasaaren perustelu:** FMI kuvaa aseman ulkosaaristoon noin 20 km
keskimääräisestä rannikkoviivasta. Se on vaihtoehto painotettaessa avoimempaa
merialtistusta pidemmän etäisyyden kustannuksella. FMI:n mukaan metsäiset
kalliokohoumat heikentävät hieman itäkoillisen tuulta. Kumpikaan asema ei
todista testipisteiden paikallistuulta tai lopullisen mökkipaikan edustavuutta.

Molemmat ovat asemametatiedossa automaattisia sääasemia (verkko 121),
toimintajakson loppu `now`. Käyttäjän hyväksymän Rankin havaintonäytteen
mittausajat, tuoreus ja tuulisuureiden äärelliset arvot on tarkistettu;
tulokset ovat luvussa 6. Haapasaaren havaintonäytettä ei haettu.

Metatietolähteet (tarkistettu 2026-11-17):

- [Asemakyselyn kuvaus](https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=DescribeStoredQueries&storedquery_id=fmi::ef::stations).
  `fmisid` on yksi kokonaisluku.
- [Rankin metatieto](https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=GetFeature&storedquery_id=fmi::ef::stations&fmisid=101030):
  HTTP 200, vastausleima 2026-11-17T18:03:39Z, yksi asema.
- [Haapasaaren metatieto](https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=GetFeature&storedquery_id=fmi::ef::stations&fmisid=101042):
  HTTP 200, sama vastausleima, yksi asema.
- Ennen yksittäishakuja tehty metatietoyritys samalla URL-pohjalla
  `fmisid=101030,101042` palautti HTTP 400. Kuvaus tarkistettiin: listan sijaan
  käytettiin yllä olevia kahta yksittäistä kokonaislukutunnistetta.
  Ei havaintohakua, automaattista uusintasilmukkaa tai muita asemaehdokkaita.

**Hyväksytty valinta:** Kotka Rankki, FMISID 101030. Käyttäjä valitsi
2026-11-17 tässä toteutuskeskustelussa kysymykseen vastauksen
”Kotka Rankki (Recommended)”. Kysymyksessä hyväksyttiin rajattu havaintonäytehaku
ja rajattiin valinta vain selvitykseen, ei lopullisen mökkipaikan edustavuuteen.

## 3. Lähdekäytön edellytykset

Kaikki tämän luvun lähteet tarkistettu **2026-11-17**, ennen sääaineistonäytteitä.

- [FMI:n lisenssi ja ehdot](https://www.ilmatieteenlaitos.fi/avoin-data-lisenssi):
  CC BY 4.0; mainitaan lisenssinantaja ja aineiston nimi. Lähdemerkintä tässä:
  **Ilmatieteen laitoksen avoin data — HARMONIE (MEPS) -pistesääennuste ja
  säähavainnot**, [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/).
  Alkuperäisvastauksia ei muokata; laskettu etäisyys ja raportin tulkinnat
  erotetaan alkuperäisdatasta. FMI:n organisaatiologoa ei käytetä.
- [WFS-palveluohje](https://en.ilmatieteenlaitos.fi/open-data-manual-fmi-wfs-services)
  ja ehdot: latauspalvelu 20 000 pyyntöä/vrk/käyttäjä, katselupalvelu
  10 000, yhteensä 600/5 min/käyttäjä. Käyttäjän tunnistustapaa tai
  rinnakkaisuus-/uusintatakuuta ei päätellä. Myös metatietohaut budjetoidaan
  varovaisesti pyynnöiksi. Ei jatkuvaa hakua eikä uusintasilmukkaa.
- [Pikaohje](https://www.ilmatieteenlaitos.fi/latauspalvelun-pikaohje):
  rekisteröitymistä ei tarvita. Rajataan yksi ennustenäyte/testipiste, enintään
  48 h ja vain tarvittavat suureet; yksi hyväksytyn aseman näyte, enintään 1 h.

### Varmennetut kyselyt ja parametrit

[Ennustekuvaus](https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=DescribeStoredQueries&storedquery_id=fmi::forecast::harmonie::surface::point::multipointcoverage):
`fmi::forecast::harmonie::surface::point::multipointcoverage`.
Sallitut parametrit: `starttime,endtime,timestep,parameters,place,latlon,geoid,fmisid,wmo`.
`latlon` on leveysaste,pituusaste; ajat ISO 8601 UTC, `timestep` minuutteina,
`parameters` pilkuin eroteltuna. Ei dokumentoitua `crs`-, `order`- tai
`origintime`-parametria tähän kyselyyn. Oletusikkuna 50 h ei ole saatavuustakuu.

[Havaintokuvaus](https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=DescribeStoredQueries&storedquery_id=fmi::observations::weather::multipointcoverage):
`fmi::observations::weather::multipointcoverage`.
Parametrit: `starttime,endtime,timestep,parameters,crs,bbox,place,fmisid,maxlocations,geoid,wmo`.
Käytetään hyväksyttyä `fmisid=101030`; ei lähimmän aseman automaattivalintaa.
`bbox` olisi pituusaste,leveysaste,pituusaste,leveysaste, ei `latlon`-järjestys.
`latlon` ja `order` eivät kuulu tämän kyselyn kuvaukseen. Oletusikkuna 12 h
rajataan yhteen tuntiin. Tuloksen `srsName` ja `axisLabels` tarkistetaan;
koordinaattijärjestelmää ei pakoteta arvauksella.

[WFS-käyttöohje](https://en.ilmatieteenlaitos.fi/open-data-manual-wfs-examples-and-guidelines)
varoittaa epätasaisista aikakohdistuksista. Käytetään tasatuntia, ennusteessa
60 min ja havainnoissa 10 min askelta. `timestep` valitsee ajanhetket,
ei laske havaintojen tuntikeskiarvoa.

### Kenttä- ja aikamerkitysten lähteet

- [Ennustesuureiden metatieto](https://opendata.fmi.fi/meta?observableProperty=forecast&param=Temperature,WindSpeedMS,WindGust,WindDirection,WeatherSymbol3,Precipitation1h,FeelsLike,ProbabilityOfPrecipitation):
  kuusi ydinkenttää; `FeelsLike` ja `ProbabilityOfPrecipitation` eivät palauta
  määritelmiä. Sanasto ei yksin todista tuottajan kenttätukea.
- [Havaintotuulen metatieto](https://opendata.fmi.fi/meta?observableProperty=observation&param=ws_10min,wg_10min,wd_10min)
  ja [havaintojen tulkinta](https://en.ilmatieteenlaitos.fi/guidance-to-observations):
  edeltävän 10 min keskituuli, saman jakson suurin 3 sekunnin keskituuli
  puuskana ja keskisuunta. 360° on pohjoisesta, 0 tyyni; puuttuva ei ole tyyni.
- [Ennustemalliohje](https://en.ilmatieteenlaitos.fi/open-data-manual-forecast-models)
  ja [prosessikuvaus](https://xml.fmi.fi/inspire/process/harmonie_scandinavia_surface):
  `analysisTime` on nimellinen malliajon viiteaika, `phenomenonTime`
  ennustejakso. `resultTime` kirjataan tulosaikana, ei todistettuna
  julkaisuaikana. HTTP-hakuaika on erillinen.
- [UTC-ohje](https://en.ilmatieteenlaitos.fi/open-data-questions):
  avoin data käyttää UTC:tä. Sateen havaintokoodeja ei yleistetä ennusteeseen.
- Pikaohjeen **WeatherSymbol3**-taulukko määrittelee koodit; SmartSymbolin
  yökoodeja ei siirretä tähän kuvastoon. Kuvakkeita ei toteuteta.

## 4. Testimökkipisteen ennustenäyte

[Alkuperäinen vastaus](evidence/weather-wind/cottage.xml) ja
[URL, hakuaika, HTTP-tulos ja SHA-256](evidence/weather-wind/manifest.md).
Haku rajattiin 2026-11-17T18:00Z–2026-11-19T18:00Z, 60 min.
XML on luettavissa, yksi `GridSeriesObservation`, kuusi nimettyä
`swe:field`-kenttää, 49 aika-/arvoriviä.

Pyydetty piste: 60.434318, 26.687825. Vastaus palauttaa
**60.43432, 26.68783**, EPSG:4326, nimi Källarbottnen, geoid 844637,
alue Pyhtää, aikavyöhyke Europe/Helsinki. Ero on alle metrin ja vastaa
koordinaattien pyöristystä. Samat koordinaatit ovat kattavuuden jokaisella
rivillä. Vastauksesta ei saada erillisen lähimmän mallihilaruudun keskikohtaa
tai maa-/meriluokitusta: palautettua pistettä ei väitetä hilaruudun
keskipisteeksi eikä geoid-nimeä sijainnin vaihtamiseksi.

### Suurematriisi

Merkitykset ja yksiköt luvun 3 ennustemetatiedosta; saatavuus
`cottage.xml`:n nimetyistä kentistä ja arvoriveistä. Kaikki kuusi kenttää
sisältävät 49/49 äärellistä arvoa. Näytteessä ei ole puuttuvaa arvoa:
mahdollisen puuttuvan arvon tarkkaa XML-esitystä ei siten ole varmennettu.
Myöhempi toteutus ei saa muuttaa NaN-/nil-/tuntematonta arvoa nollaksi.

| Suure / tarkka kenttä | Yksikkö | Dokumentoitu merkitys / jakso | Näytteen saatavuus |
| --- | --- | --- | --- |
| Lämpötila `Temperature` | degC | Ilman lämpötila; metatieto ei ilmoita keskiarvojaksoa. Mallisivun mukaan 2 m korkeus. | 49 äärellistä, 2,8–6,3 °C |
| Keskituuli `WindSpeedMS` | m/s | `avg`, PT1H; mallin 10 m tuuli, ei havaintoaseman 10 min keskiarvo | 49 äärellistä, 4,36–7,21 |
| Puuska `WindGust` | m/s | `max`, PT1H; sisempi `avg`, PT3S | 49 äärellistä, 8,0–12,7 |
| Suunta `WindDirection` | deg | `avg`, PT1H; tulosuunta alla mainitun FMI:n termiohjeen perusteella | 49 äärellistä, 62–206° |
| Sääsymboli `WeatherSymbol3` | index | Kategorinen koodi, ei jatkuva mittaluku | 49 äärellistä; esiintyvät 3 = pilvistä, 31 = heikkoa vesisadetta, 32 = vesisadetta |
| Tuntisade `Precipitation1h` | **mm/h** lähteen merkintänä | **sum**, PT1H, ”Precipitation amount 1 hour”; kertymä, ei hetkellinen intensiteetti tai ajon alusta kertynyt sade | 49 äärellistä, 0,0–1,6 |
| Suora näennäinen lämpötila | Ei varmennettu | `FeelsLike` ei saanut määritelmää ennustesanastosta; ei mukana valitun lähteen näytteessä | Saatavuus jää varmentamatta; ei omaa tuntuu kuin -laskentaa |
| Sateen todennäköisyys | Ei varmennettu | `ProbabilityOfPrecipitation` ei saanut määritelmää; myöskään erillinen sanastotarkistus `PoP`-nimellä ei tuottanut määritelmää | Saatavuus jää varmentamatta; ei todennäköisyyden johtamista sateesta |

[FMI:n Tuulet-termiohje](https://www.ilmatieteenlaitos.fi/tuulet), tarkistettu
2026-11-17: ”Tuulen suunnalla tarkoitetaan suuntaa, josta tuuli puhaltaa”;
etelätuuli on noin 180°. Tätä käytetään `WindDirection`-kentän
”Wind direction” -määritelmän tulkintaan. Ennusteen tyynen erikoiskoodia ei
varmennettu; havaintojen 0/360-erottelua ei siirretä malliin oletuksena.

Valitun tuottajan puuttuvia toissijaisia kenttiä ei väitetä mahdottomiksi
kaikissa FMI-tuotteissa. Lisälähdettä ei oteta vain näiden täydentämiseksi,
koska seuraavan pienen lämpötila-/tuulitoiminnon ydin saadaan tästä lähteestä.

### Ajat, kattavuus ja sade

- `analysisTime`: **2026-11-17T12:00:00Z** (nimellinen malliajo).
- `resultTime`: **2026-11-17T15:04:04Z** (lähteen tulosaika).
- Vastauksen `timeStamp`: **2026-11-17T18:06:10Z**; erillinen
  pyyntöaloitus **18:06:12.261681Z**. Vastauksen leima ei ole
  luotettava julkaisuaika; sen muutaman sekunnin eroa paikalliseen hakuaikaan
  ei korjata eikä käytetä viivelupauksena.
- Voimassaoloajat kaikille kuudelle suureelle:
  **2026-11-17T18:00Z–2026-11-19T18:00Z**, 49 riviä,
  tasainen 3600 s askel, ei duplikaatteja, aikajärjestysaukkoja tai
  ei-äärellisiä arvoja. Nimellinen ikkuna 48 h.
- Todellinen viimeiseen äärelliseen pisteeseen ulottuva tuleva kattavuus
  hakuaikahetkestä on **47 h 53 min 47,738319 s**, ei täysi 48 h.
  Ensimmäinen rivi oli hakuhetkellä jo menneisyydessä; se ei pidennä tulevaa
  kattavuutta. Yksi näyte ei todista tulevien hakujen kattavuutta.

`Precipitation1h` on dokumentoitu **tunnin summa**, vaikka sanaston yksikkö
on `mm/h`. Yksikkömerkintä säilytetään näkyvänä; lukua ei nimetä hetkelliseksi
sateen intensiteetiksi eikä ajon alusta kertyneeksi sateeksi. Metatieto ei
kerro yksiselitteisesti, sijoittuuko kertymä ennen vai jälkeen aikaleiman.
Siksi tästä ei vielä hyväksytä Suomen kalenterivuorokauden sadesummaa.

Vuorokausisumma edellyttäisi tuottajalta varmennettua aikaleiman
kertymäankkuria, toisensa poissulkevia yhden tunnin jaksoja,
kaikkien tarvittavien rajajaksojen äärellisiä arvoja sekä
`Europe/Helsinki`-kalenteripäivän muuntamista UTC-rajoiksi.
Kesäaikasiirtymäpäivä on 23 tai 25 h, ei aina 24 h. Puuttuvaa jaksoa ei
täytetä nollalla; osittaista päivää ei esitetä kokonaisena.
Tässä ei lasketa käyttöönotettavaa summaa.

### Offline-toistaminen

Repossa, ilman verkkohakuja tai riippuvuuksien asennusta:

```sh
python3 docs/feasibility/evidence/weather-wind/verify.py
```

[Tarkistusskripti](evidence/weather-wind/verify.py) käyttää vain Pythonin
standardikirjastoa. Se tarkistaa muuttamattomat tavut SHA-256:lla, XML:n
luettavuuden (ei ulkoista XSD-validointia), yhden kohteen, kenttäjärjestyksen,
koordinaatit, aika-/arvorivien määrät, ensimmäisen ja viimeisen ajan,
tarkan askelluksen, lähdeajat ja suurekohtaiset äärelliset arvot.
Se tulostaa minimit/maksimit, viimeiset arvot ja tulevan kattavuuden sekä
tarkistaa raportin ja manifestin paikalliset linkit. Se ei päättele
kertymäankkuria tai yksikköä numeroista eikä hae näytteitä uudelleen.
Suoritettu mökkinäytteelle 2026-11-17: kaikki tarkistukset läpi.

## 5. Testiulappapisteen tuuliennuste

[Vastaus](evidence/weather-wind/offshore.xml),
[alkuperätiedot](evidence/weather-wind/manifest.md). Käyttäjän valitsema piste
**60.414616, 26.697395** palautuu **60.41462, 26.69740**, EPSG:4326.
Nimi Lilla Krokö, geoid 844587; ero alle metri ja vastaa pyöristystä.
Nimi ei todista mallin maa-/meriluokitusta. Kuten mökillä, erillistä
mallihilaruudun keskikohtaa tai sen maa-/merimaskia ei palauteta.
Näyte ei osoita pistettä vaihtavaa sijaintiestettä; merialtistuksen
lopullinen validointi jää puutteeksi, ei hiljaiseksi pisteenvaihdoksi.

| Kenttä | Yksikkö ja merkitys | Tulos |
| --- | --- | --- |
| `WindSpeedMS` | m/s, PT1H keskiarvo | 49/49 äärellistä, 5,67–8,95 |
| `WindGust` | m/s, PT1H maksimi 3 s keskiarvoista | 49/49 äärellistä, 8,3–13,5 |
| `WindDirection` | deg, PT1H keskisuunta, tulosuunta luvun 4 termilähteellä | 49/49 äärellistä, 102–207° |

Lähdemerkitykset ovat samat kuin mökkiennusteessa, eivät aseman
10 minuutin tilastoja. Puuttuvan arvon esitystä tai tyynen mallikoodia ei
tästäkään näytteestä voida vahvistaa; ei puuskien korvauslaskentaa.

`analysisTime` **2026-11-17T15:00Z**, `resultTime`
**2026-11-17T18:03:58Z**. Näyte palautti eri malliajon kuin noin kolme
minuuttia aiemmin haettu mökkinäyte. Niiden eroja ei siis voi selittää
yksin sijainnilla. Ajot kirjataan erikseen, eikä vanhaa näytettä haeta
uudelleen vain vertailun tasoittamiseksi.

Kaikilla suureilla voimassaolo **2026-11-17T18:00Z–2026-11-19T18:00Z**,
49 riviä tunnin välein; ei aika-aukkoja, duplikaatteja tai puuttuvia arvoja.
Todellinen tuleva kattavuus hakuhetkestä **47 h 51 min 4,576399 s**.
Luvun 4 offline-komento suoritettu myös tälle näytteelle 2026-11-17:
SHA-256, XML-rakenne, kentät, piste ja ajallinen kattavuus tarkistettu.

## 6. Hyväksytyn Rankin tuulihavainto

[Vastaus](evidence/weather-wind/rankki.xml),
[alkuperätiedot](evidence/weather-wind/manifest.md). Asematunniste
**101030**, nimi Kotka Rankki, geoid -16000076 ja WMO 2976 vastaavat
hyväksyttyä asemaa. Vastauspiste **60.37538, 26.95893**, EPSG:4258,
vastaa metatietopisteen pyöristystä.

Hakuikkuna on viimeisin **päättynyt tasatuntiväli 17:00–18:00Z**,
2026-11-17; kesto 1 h, ei 12 h oletusikkunaa.
Haku aloitettiin 18:09:56.139618Z. Näytteen viimeinen mittaus on siten
9 min 56,139618 s vanha suhteessa hakuaikaan. Tämä on rajatun näytteen
viimeinen arvo, ei todistus siitä, ettei ikkunan jälkeen olisi uusia arvoja.
Mittaus-/tilastoaika luetaan rivistä, ei HTTP-haun tai `resultTime`:n perusteella.
`resultTime` on 18:00Z; havainnoilla ei ole malliajon viiteaikaa.

| Kenttä | Yksikkö ja jakso (§3:n havaintolähteet) | Näytteen viimeisin äärellinen |
| --- | --- | --- |
| `ws_10min` | m/s, edeltävän 10 min keskinopeus | **7,6 m/s**, 2026-11-17T18:00Z |
| `wg_10min` | m/s, edeltävän 10 min suurin 3 s keskinopeus | **12,0 m/s**, sama aika |
| `wd_10min` | deg, edeltävän 10 min keskisuunta, mistä tuulee | **182°**, sama aika |

Jokaisella suureella 7/7 äärellistä arvoa, 17:00Z–18:00Z,
600 s askel, ei puuttuvia aikaleimoja tai duplikaatteja. Alueet
keskituuli 7,1–8,8 m/s, puuska 12,0–13,6 m/s ja suunta 176–182°.
Näyte ei sisällä puuttuvia arvoja eikä tyventä; puuttuvan XML-koodausta ei
varmenneta tästä. Puuska on lähteen havainto, ei laskettu kerroin.
Jakson ensimmäisen rivin 10 min tilasto ulottuu hakuikkunaa edeltävään aikaan;
seitsemän aikaleimaa ei tarkoita seitsemää erillistä tulevaa tuntia.

Offline-komento suoritettu kaikille kolmelle näytteelle: myös aseman
tunniste, kenttäkohtaiset viimeiset äärelliset arvot ja aikavälit tarkistettu.

## 7. Päivitysrytmit ja ehdotetut tuoreusrajat

Dokumentaatiot tarkistettu 2026-11-17:

| Aineisto | Dokumentoitu rytmi / rajoitus | Mitä ei ole varmennettu |
| --- | --- | --- |
| Molempien pisteiden HARMONIE-ennuste | [WFS-kuvaus](https://opendata.fmi.fi/wfs?service=WFS&version=2.0.0&request=DescribeStoredQueries&storedquery_id=fmi::forecast::harmonie::surface::point::multipointcoverage): uusi aineisto 6 h välein, oletushaku 50 h. [Operatiivisen mallin sivu](https://www.ilmatieteenlaitos.fi/avoin-data-saaennustedata): 8 ajoa/vrk, 66 h ennuste, 2,5 km hila; vain viimeisin ajo, saatavuutta lyhentävät laskenta- ja siirtoviiveet. Sivun kellonaikaluettelossa on kuitenkin vain seitsemän aikaa (21 UTC puuttuu). [Prosessikuvaus](https://xml.fmi.fi/inspire/process/harmonie_scandinavia_surface): 8 ajoa/vrk, 54 h, vanhempi malliversio. | WFS:n täsmällinen nykyinen julkaisurytmi ristiriitainen; 3 h laskentarytmi, 6 h julkaisu, 50 h oletusikkuna, 54/66 h mallikattavuus eivät ole sama lupaus. Näytteiden eri ajot osoittavat päivittymisen, eivät vakioritmiä. |
| Rankin havainto | [Pikaohje](https://www.ilmatieteenlaitos.fi/latauspalvelun-pikaohje): useimmilla asemilla 10 min havaintoväli; näytteessäkin 10 min | Mittausväli ei ole saapumisviive eikä tuoreustakuu; Rankin palvelukohtaista viive-SLA:ta ei löydetty. |

**Seuraavat ovat suunnitteluehdotuksia, eivät käyttäjän hyväksymiä rajoja
tai tässä toteutettuja sääntöjä:**

- Piste-ennusteiden keskitetty haku **60 min välein**: pienempi kuin
  kummankin dokumentoidun 3/6 h rytmin jakso, mutta ei turhaa minuuttipollausta.
  Saman malliajon uudelleen saaminen ei tarkoita uutta ennustetta.
  Näytä hakuhäiriö viimeistään **2 h** viimeisestä onnistuneesta hausta
  (kaksi väliin jäänyttä hakua). Merkitse aineisto vanhentuneeksi myös, jos
  mallin viiteaika on yli **12 h** vanha (kaksi dokumentoidun hitaamman
  6 h jaksoa); tämä on varovainen ehdotus julkaisuristiriidan vuoksi.
  Tuleva kattavuus tarkistetaan aina kentittäin: voimassaoloajan päättyminen
  estää päättyneiden rivien näyttämisen tulevana ennusteena rajoista riippumatta.
- Rankin haku **10 min välein**, vanhentumisvaroitus mittausajan iän
  ylittäessä **30 min** (kolme odotettua mittausväliä).
  HTTP 200 tai uusi hakuaika ei nuorenna vanhaa havaintoa.
  Varoitusraja ei vielä ratkaise vanhan arvon lopullista piilotusrajaa;
  sen käyttäytyminen sovitaan havaintotoiminnon omassa suunnitelmassa.

Rajojen perustana ovat tämän aineiston rytmit ja puuttuva viivetakuu,
eivät Haminan vedenkorkeuden 30 min / 12 h haku-/vanhentumisrajat.
Sääennusteessa esiintyvä 12 h on erikseen perusteltu malliviiteajan
ikäehdotus. Häiriön esitys, tyhjät kentät ja raja-arvot vaativat myöhemmän
OpenSpec-päätöksen sekä häiriötestit; näillä ehdotuksilla ei hyväksytä
toimintoa tai luvata jatkuvaa saatavuutta.

## 8. Myöhempi Linux-VM-/Proxmox-testaus

Käyttäjän suunnitteluvaiheessa ilmoittama Proxmox-palvelin voi myöhemmin
tarjota Linux-VM:n integraatio-, verkkokatko-, uudelleenkäynnistys- ja
julkaisuketjun testeihin sekä palvelun oikeiden mobiiliselainten testejä varten.
Tässä ei selvitetä palvelimen kokoonpanoa, tehdä asennuksia, kerätä tunnuksia,
muuteta verkkoa tai avata palvelua. VM:n OS, resurssit, arkkitehtuuri,
osoitteet ja käyttöönotto ratkaistaan erillisessä muutoksessa.

| Näyttö | Mitä sillä voi myöhemmin osoittaa | Mitä se ei hyväksy |
| --- | --- | --- |
| Linux-VM:n integraatio- ja häiriötestit | Toiminta kyseisessä VM-ympäristössä | Pi:n muistia, lämpöä, I/O:ta tai suorituskykyä |
| Oikea puhelin VM-palvelua vasten | Kyseisen mobiiliselaimen toiminta oikealla laitteella | Pi:n laiteajoa; pelkkä selainemulointi ei ole tämä testi |
| arm64-ristiinkäännös | Kohdebinaarin muodostuminen | Binaarin käynnistymistä tai nopeutta oikealla Pi:llä |
| Oikea Raspberry Pi -ajo | Laitteen toimivuus ja sovitut laite-/suorituskykykriteerit | Muiden ympäristöjen automaattista hyväksyntää |

Mitään näistä neljästä testistä ei tehty tässä dokumentaatioselvityksessä.
Pi:n kunnostus odottaa käyttäjän päätöksen mukaisesti; **issue #1 säilyy
tekemättömän Pi-laitehyväksynnän viitteenä**, eikä sen hyväksymisehtoja
poisteta tai issuea suljeta VM-ajoihin vedoten.

## 9. Jatkosuositus ja hyväksymisehtojen läpikäynti

**Suositus yhdeksi pieneksi seuraavaksi toiminnoksi:** testimökkipisteen
lämpötila- ja tuuliennuste (lämpötila, keskituuli, puuska ja tulosuunta)
todelliseen saatavilla olevaan, enintään 48 h aikajänteeseen. Tavoite olisi
yksi kokonainen lähde–palvelin–näkymä-ketju nykyisen vedenkorkeustoiminnon
rinnalle, ei koko loppu-MVP tai yleinen kaikkien aineistojen kehys.
Tässä **ei käynnistetä** sitä koskevaa suunnittelua tai sovellustoteutusta.

**Varmennettu:** nimettyjen neljän ydinkentän lähdemerkitykset ja yksiköt,
kohdistuva todellinen mökkivastaus, 49 äärellistä tuntipistettä ilman aukkoja,
erilliset lähdeajat ja lähes 48 h tuleva kattavuus. Ulappapisteen kolme
tuulikenttää ja hyväksytyn Rankin kolme havaintokenttää on lisäksi saatu;
niiden näyttö ei kasvata seuraavan ehdotetun toiminnon rajausta.
Myös sääsymboli ja tuntisade ovat numeerisesti saatavilla, mutta niitä ei
ehdoteta samaan ensimmäiseen pieneen toimintoon.

**Puutteet:** suora näennäinen lämpötila ja sateen todennäköisyys jäävät
varmentamatta; kertymäankkuri estää hyväksymästä vuorokausisummaa.
Täsmällinen WFS-julkaisurytmi on ristiriitainen, erillistä hilaruudun
maa-/merimaskia ei saatu eikä puuttuvien arvojen koodausta tai mallin
tyynen koodia esiinny näytteissä. Saatavuuden jatkuvuutta ei testattu.
Nämä on kirjattava myöhempään suunnitelmaan, eivät hiljaisia
korvauksia tai hyväksyttyjä poikkeuksia.

**Käyttäjälle palautettavat päätökset:** hyväksytäänkö yllä oleva rajattu
seuraava toiminto erilliseen OpenSpec-suunnitteluun; sen haku-/vanhentumisrajat,
puuttuvan tiedon ja lyhyemmän kattavuuden esitys sekä käyttöönottopaikka.
Lopullinen mökkisijainti ja havaintoaseman edustavuus arvioidaan uudelleen
ennen käyttöönottoa. Täyden 48 h jatkuvaa takuuta ei ehdoteta.
Mitään ydintietoa ei hyväksytä korvattavaksi laskennallisesti.

Designin §6:n läpikäynti:

| Ehto | Näyttö / raja |
| --- | --- |
| Testipisteet ja valintojen alkuperä, aseman hyväksyntä | §1–2; Rankki hyväksytty ennen havaintohakua |
| Suurematriisi, merkitykset ja yksiköt tai nimenomainen puute | §3–6; toissijaiset puutteet ja sadesumman este näkyvissä |
| Todellinen kohdistettu vastaus jokaisesta aineistosta | Kolme muuttamatonta XML:ää ja manifestin URL:t; HTTP 200:n lisäksi rakenne, sijainti ja äärelliset arvot tarkistettu |
| Offline-toistettavuus ja vedenkorkeusnäytteiden säilyminen | §4:n komento; lopputarkistukset §10 |
| Varmennettu tieto, ehdotukset ja päätökset eroteltu | §7–9; vain yksi rajattu jatkoehdotus, ei uutta sovellustoteutusta |
| Ei MVP-/sijainti-/laitehyväksyntää | §1 ja §8; oikeat mobiili-, VM- ja Pi-testit tekemättä, issue #1 säilyy |

## 10. Toteutuneet tarkistukset ja puuttuva näyttö

Suoritettu 2026-11-17:

- `python3 docs/feasibility/evidence/weather-wind/verify.py` — läpi.
  Kaikkien kolmen XML:n SHA-256 vastaa manifestia ja kiinteitä odotuksia;
  rakenteet, pyydetty/palautettu kohde, suureet, rivimäärät, lähdeajat,
  askellus ja kenttäkohtainen äärellinen kattavuus tarkistettu.
  Myös manifestin kyselyiden kentät, aikaikkunat (48 h / 1 h),
  sijainti-/asematunnisteet ja raportin/manifestin paikalliset linkit tarkistettu.
- `openspec validate verify-weather-wind-feasibility --strict` — läpi.
  `skip_specs` on tarkoituksellinen: ei uusia sovelluskäyttäytymisen deltoja.
- `git diff --cached --check -- . ':(exclude)docs/feasibility/evidence/weather-wind/*.xml'`
  — läpi. Rajaamaton staged-tarkistus palauttaa 2: FMI:n alkuperäisissä
  XML-vastauksissa on rivinloppuvälilyöntejä ja välilyönti–tabulaattorisisennyksiä.
  Niitä ei siivota, koska näytteiden tavut ja SHA-256 säilytetään muuttamattomina.
  Muutosten rajaus: vain tämä raportti,
  uusi `evidence/weather-wind/` ja muutoksen tehtävämerkinnät.
  Aiemmat vedenkorkeusnäytteet sekä `cmd/`, `internal/` ja `web/`
  säilyvät muuttamattomina.
- `go test ./...` — kaikki viisi pakettia läpi.
- FMI:n dynaamiset asemakuvaukset tarkistettu Playwrightilla; selain suljettu.
  Tämä ei ole sovelluksen käyttöliittymätesti.

Sovelluksen käyttöliittymää tai riippuvuuksia ei muutettu, joten web-buildiä
tai selainhyväksyntätestejä ei ajettu tämän dokumentaatiomuutoksen vuoksi.
Ulkoista XSD-validointia, pitkittäistä saatavuusseurantaa, kuormitustestejä,
VM-ajoa, arm64-käännöstä eikä oikeaa mobiili-/Pi-hyväksyntää tehty.
Puuttuvan/tyynen arvon XML-esitystä ei tuotettu keinotekoisella FMI-haulla.
Säädokumentaation ristiriidat ja puutteet säilyvät §4–9:ssä;
offline-toisto ei todista verkkodokumenttien tulevaa sisältöä.

Selvityksen valmistuminen ei arkistoi muutosta, avaa seuraavaa
toteutusmuutosta eikä anna lupaa yhdistää `main`iin. Toimitus tapahtuu PR:n
ja katselmoinnin kautta; yhdistäminen odottaa erillistä lupaa ja CI-tarkistuksia.
