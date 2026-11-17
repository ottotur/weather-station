# Vedenkorkeuden XML: aika ja yksikkö

Selvitys 17.11.2026. Tarkastelu perustuu tallennettuihin esimerkkeihin ja ensisijaisiin skeemoihin; live-FMI-dataa ei haettu. XSD:n sallima rakenne ei yksin todista FMI:n kaikkia tuotantokäytäntöjä tai mahdollisten profiili- ja Schematron-sääntöjen täyttymistä.

## Varmistetut muodot ja vertailu

- **Havaintoaika:** O&M:n `om:phenomenonTime` on `om:TimeObjectPropertyType`: sen sisällä voi olla GML-aikaobjekti tai ominaisuudessa voi olla assosiaation `xlink:href`. `om:resultTime` on vastaavasti `gml:TimeInstantPropertyType`. Ilmiöaika kuvaa aikaa, johon tulos pätee; tulosaika sitä, milloin tulos tuli saataville. [O&M observation.xsd](https://schemas.opengis.net/om/2.0/observation.xsd), [GML temporal.xsd](https://schemas.opengis.net/gml/3.2.1/temporal.xsd), [GML gmlBase.xsd](https://schemas.opengis.net/gml/3.2.1/gmlBase.xsd).
- **Aikajakson sisäinen rakenne:** `gml:TimePeriod` vaatii alulle joko `gml:beginPosition`- tai `gml:begin`-elementin ja lopulle joko `gml:endPosition`- tai `gml:end`-elementin. `begin` ja `end` ovat `TimeInstantPropertyType`-ominaisuuksia, joten ne voivat sisältää `TimeInstant`in tai viitata siihen `xlink:href`-attribuutilla. Koko aikajaksoon viitataan sitä vastoin **omistavassa ominaisuuselementissä**, esimerkiksi `<om:phenomenonTime xlink:href="#time-interval-1-1"/>`, ei korvaamalla `gml:TimePeriod`in pakollisia rajoja sen omalla `xlink:href`-attribuutilla. [GML temporal.xsd](https://schemas.opengis.net/gml/3.2.1/temporal.xsd).
- **Paikalliset esimerkit:** [havainto](evidence/hamina-observation.xml) sisältää suorat alku- ja loppuajankohdat ja erillisen tulosajan. [Ennusteen](evidence/hamina-forecast.xml) ensimmäinen jäsen sisältää jakson `gml:id="time-interval-1-1"` ja tulosajan `gml:id="time-1-1"`; toinen jäsen viittaa niihin omissa `om:phenomenonTime`- ja `om:resultTime`-elementeissään. Ennusteessa on lisäksi erillinen `analysisTime`: se ei ole sama kenttä kuin tulosaika.
- **Pisteiden aikamuoto:** WaterML:n `wml2:MeasurementTVP/wml2:time` on `gml:TimePositionType`, XSD:ssä valinnainen. GML:n unioni sallii esimerkiksi `dateTime`, `date`, vuoden, kellonajan, desimaaliluvun ja URI:n; myös `indeterminatePosition` voi ilmaista epätarkan ajan. Molempien tallenteiden pisteet käyttävät täsmällistä UTC-merkintää `...Z`. Kaikkia skeeman sallimia muotoja ei siis voi olettaa samaan yksinkertaiseen aikaleimamuotoon. [WaterML timeseries.xsd](https://schemas.opengis.net/waterml/2.0/timeseries.xsd), [GML temporal.xsd](https://schemas.opengis.net/gml/3.2.1/temporal.xsd).
- **WaterML-yksikkö:** `wml2:value` voi sisältää valinnaisen `uom`-attribuutin. Vaihtoehtoisia metatietopaikkoja ovat sarjan `wml2:defaultPointMetadata/wml2:DefaultTVPMeasurementMetadata/wml2:uom` sekä pisteen `wml2:metadata/wml2:TVPMeasurementMetadata/wml2:uom`; jälkimmäiset ovat `swe:UnitReference`-tyyppiä, jossa on `code` ja XLink-attribuutit. Yksikkökoodin on tarkoitus olla UCUM-koodi. XSD sallii näiden yksikkötietojen puuttumisen; tämä ei yksin osoita, että yksikötön arvo olisi semanttisesti riittävä. [WaterML timeseries.xsd](https://schemas.opengis.net/waterml/2.0/timeseries.xsd), [SWE basic_types.xsd](https://schemas.opengis.net/sweCommon/2.0/basic_types.xsd).
- **`omop:uom` ei ole WaterML-pisteen yksikköelementti:** [tallennetussa havaintosuureen metatiedossa](evidence/observation-property.xml) juuren oletusnimiavaruus on `http://inspire.ec.europa.eu/schemas/omop/2.9`, joten `<uom uom="mm"/>` on kyseisen nimiavaruuden `ObservableProperty`-juuren suora lapsi; sen `uom`-attribuutti on nimiavaruudeton. Havaintovastauksen `om:observedProperty` viittaa metatietoon erillisellä URL-osoitteella. Havainto- ja ennuste-XML:issä ei ole vastaavaa yksikköilmoitusta. INSPIREn [3.0-skeema](https://inspire.ec.europa.eu/schemas/omop/3.0/ObservableProperties.xsd) sijoittaa valinnaisen `uom`-elementin niin ikään `ObservableProperty`n lapseksi, mutta on eri nimiavaruuden versio eikä todista 2.9-skeeman validiutta.

## Aukot ja epäonnistumiset

- Tallenteen ilmoittama [INSPIRE omop 2.9 -luonnosskeeman osoite](https://inspire.ec.europa.eu/draft-schemas/omop/2.9/ObservableProperties.xsd) palautti INSPIRE-rekisterin «Page not found» -sivun, ei XSD:tä. Siksi 2.9:n sallitut yksikkömuodot ja tallenteen skeemavaliditeetti jäivät varmistamatta; 3.0:aa ei voi käyttää sen korvikkeena.
- Havaintometatieto ilmoittaa `mm`, mutta tallenteista ei yksin seuraa ennustesuureiden yksikkö, pystydatum tai yksikköjen mahdollinen etusijajärjestys ristiriitatilanteessa. Niille tarvitaan erillinen FMI:n ensisijainen sopimus/metatieto. Myöskään `omop:uom`-metatiedon yleistä periytymistä WaterML-arvoihin ei tässä vahvistettu.
- GML/XSD sallii enemmän aikamuotoja ja viittauksia kuin esimerkit näyttävät. Ei ole vahvistettu, mitä niistä FMI:n juuri nämä WFS-kyselyt takaavat tai tuottavat; viittausten resolvointi, puuttuva kohde ja etäviitteiden turvallinen käsittely ovat toteutuspäätöksiä. Täyttä skeema- tai Schematron-validointia ei tehty.

## Rajattu tulkinta ja testimatriisi

Sovellus ei pyri yleiseksi XSD-validoijaksi: se edellyttää tämän toiminnon
tarkkoja UTC-aikoja ja rajaa hyväksytyt yksiköt dokumentoidun tulkinnan
mukaisesti. Tallennetut vastausnäytteet säilyvät muuttamattomina;
poikkeustapaukset tehdään testeissä niistä erillisinä synteettisinä muunnoksina.

| Tapaus | Testattava tulos | Testi `internal/fmi/parse_test.go` |
| --- | --- | --- |
| Yksikkö puuttuu | Havainto mm; ennuste cm **käyttäjän hyväksymänä oletuksena** | `TestUnitMetadata`, alkuperäiset näytteet |
| Sarjan tai pisteen ilmoitettu odotettu yksikkö | Hyväksytään, mukaan lukien WaterML:n oletus- ja pistemetadata sekä arvon `uom` | `TestUnitMetadata`, `TestWaterMLUnitMetadataLocations`, `TestPointValueUnitAttributes` |
| Ristiriitainen, tyhjä tai tuntematon yksikkö | Hylätään; ei automaattista skaalausta tai viitteen hakua | `TestUnitMetadata`, `TestUnitMetadataNamespaceMatrix`, `TestWaterMLUnitMetadataLocations`, `TestPointValueUnitAttributes` |
| Tuntemattoman nimiavaruuden yksikköelementti valitussa sarjassa | Hylätään varovaisuusperiaatteella; toiselle suureelle kuuluvaa sarjaa ei käytetä N2000:n korvikkeena | `TestUnitMetadataNamespaceMatrix`, `TestOtherForecastSeriesUnitDoesNotAffectN2000` |
| Suora UTC-aika tai yksiselitteinen paikallinen aikaobjektiviite | Hyväksytään, myös jakson alku-/loppuominaisuuksissa | `TestObservationSample`, `TestForecastSample`, `TestTimePeriodBoundaryReferences` |
| Ratkaisematon tai ulkoinen aikaviite; virheellinen viite GML:n aikaobjektissa | Hylätään ilman verkkohakua | `TestUnresolvedPhenomenonTimeRejected`, `TestTimePeriodNestedReferenceMatrix`, `TestTimePeriodBoundaryReferences` |
| Epätarkka aikaleima tai tuntematon aikaviitejärjestelmä | Hylätään: ei esitetä ajankohtaa tarkkana. Oletus tai nimenomainen `frame="#ISO-8601"` hyväksytään. | `TestIndeterminateTimesRejected`, `TestExplicitISOFrameAccepted` |

`omop:uom`-elementin testitapaus on puolustava synteettinen tapaus:
käytettävissä oleva havaintometatietonäyte ei osoita sen esiintyvän näiden
WFS-vastausten valitussa aikasarjassa. INSPIRE 2.9 -skeemaa ei saatu, eikä
sen sallimia käyttökonteksteja tässä väitetä varmennetuiksi.

Tunnetut yksikköviitteet rajataan täsmällisesti QUDT:n määritelmiin
[millimetri](https://qudt.org/vocab/unit/MilliM) ja
[senttimetri](https://qudt.org/vocab/unit/CentiM). Muu ilmoitettu viite
hylätään, vaikka se voisi viitata samaan yksikköön: parseri ei hae viitteitä
verkosta eikä päättele niiden merkitystä URL:n viimeisen osan perusteella.
Sarjojen oletus- ja yksittäisen pisteen yksiköiden mahdollisen etusijan
sijasta sovellus hylkää ristiriidan varovaisuusperiaatteella.
