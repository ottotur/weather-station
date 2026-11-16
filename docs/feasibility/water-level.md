# Vedenkorkeustoiminnon toteutettavuus: Hamina

Tarkistuspäivä: **2026-11-16**. Muutos: `verify-water-level-feasibility`.
Tämä on rajattu dokumentaatioselvitys, ei sovellustoteutus, asennus tai MVP:n hyväksyntä.

## Tila ja näytön tulkinta

**Jatkosuunnittelu: estynyt selvityksen ollessa kesken.**
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

## 3. Tallennetut vastaukset ja keskeneräinen tarkistusmatriisi

Alkuperäiset vastaukset ja niiden täsmälliset URL:t, UTC-hakuajat,
SHA-256-tiivisteet sekä dokumentaatiolähteet ovat
[manifestissa](evidence/manifest.md). Havainto- ja ennustehaku onnistuivat
HTTP- ja XML-tasolla. Asemat vastaavat toisiaan: havaintojen FMISID 134254
ja ennusteen geoid -10022823 yhdistyvät tallennetun asemametatiedon kautta.
Ennusteen vastausmuoto on XML/GML + WaterML 2.0, samoin havainnon aikasarja.

Alla olevat näytehavainnot tarkistettu **2026-11-16** tallennetuista vastauksista.
Matriisi on **keskeneräinen**, koska ennusteen yksikköä ei saatu varmennetuksi.

| Tarkistus | Havainto | Ennuste |
| --- | --- | --- |
| Suure | `WLEVN2K_PT1S_INSTANT`, laadunvalvottu vedenkorkeus | `SeaLevelN2000`; lisäksi `SeaLevel`, jota ei käytetä N2000-arvona |
| N2000-osoitus | Varmennettu: `observation-property.xml` sanoo N2000 | Parametrin nimi viittaa N2000:een, ja ennustekuvaus sanoo aineiston olevan saatavilla N2000:ssa. Parametrin oma metatieto on tyhjä. |
| Yksikkö | Varmennettu `mm` metatiedosta; myöhempi cm-esitys vaatii jaon kymmenellä | **Avoin/este:** ei yksikköä vastauksessa; `forecast-property.xml` on tyhjä. Suuruusluokasta tai verkkosivun esitysyksiköstä ei päätellä rajapinnan yksikköä. |
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

### Havaittu este ja seuraava päätös

**Ennusteen yksikköä ei vielä tunneta lähteisiin perustuen.** Vastauksen oma
metatietolinkki tarkistettiin sellaisenaan ja tallennettiin: HTTP 200, mutta
tyhjä `ObservableProperty`. Ennuste on olemassa; puute koskee sen tulkintaa.
Mittauksen ja ennusteen erotusta tai nousu- ja laskuyhteenvetoja ei saa
laskea ennen yksikön varmistamista. [FMI:n tukipyyntö](https://en.ilmatieteenlaitos.fi/open-data-support-request)
on mahdollinen jatkoselvitys; yhteydenottoa ei ole tehty tämän muutoksen nimissä.

Käyttäjälle palautettava kysymys: jatketaanko yksikön varmistamista FMI:n
täydentävästä virallisesta dokumentaatiosta tai käyttäjän välittämällä
FMI-vastauksella, vai viimeistelläänkö selvitys tämän esteen dokumentoivana
tuloksena? Yksikköä ei hyväksytä oletuksena kummassakaan tapauksessa.

## 4. Avoimet tarkistukset

- Käyttöjärjestelmä–Go-ehdotuksen käyttäjähyväksyntä ja myöhempi laiteajotesti.
- Ennusteen yksikkö, tarkistusmatriisin viimeistely, epävarmuusrajat,
  julkaisurytmin arvio ja toistamisohjeet.
- Näyttöön perustuva loppuarvio ja jäljellä olevat käyttäjän päätökset.
