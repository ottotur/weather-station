# Design

## Context

Motivaatio: `proposal.md`, Why. Projektissa on dokumentaatio ja OpenSpec-runko, ei sovellusta, riippuvuuksia, testejä tai rakennusasetuksia. Pääspeksifikaatioita ei ole. Lähdenäyttö: `docs/feasibility/water-level.md` ja sen kuusi alkuperäistä XML-näytettä manifestissa.

Käyttäjä hyväksyi tässä keskustelussa 30 min / 12 h vanhentumisrajat, häiriöesityksen, 10 min / 30 min hakurytmin ja kehityskone-/Pi-hyväksynnän erottamisen. Kehitysympäristö on Windows 11 + Ubuntu WSL; tämän työympäristön `uname -m` on `x86_64`. WSL-versiota tai Windowsin selaimia ei ole vielä testattu.

Hyväksytty Pi-kohde on Raspberry Pi 3 Model B+ / Raspberry Pi OS Lite 64-bit / Debian 13 Trixie, Go `linux/arm64`. Go:n ajantasainen tuettu korjausversio tarkistetaan toteutuksen alussa, eikä selvityksen versionumeroa pidetä pysyvänä valintana. Pi ei ole käyttökunnossa eikä laiteajoa ole tehty. Projektikuvauksen vanhat Pi 2 -maininnat ja alustava OS-/Go-tila eivät kumoa käyttäjän myöhempiä hyväksyntöjä; projektikuvausta ei muuteta tässä suunnittelumuutoksessa.

## Goals / Non-Goals

**Goals:**

- Erotetaan lähdemuodon tulkinta, tuoreussäännöt, pysyvyys ja esitys, jotta koko ketju voidaan testata ilman internetiä tai Pi:tä.
- Sama rakennettu sovellus palvelee käyttöliittymän ja API:n samasta alkuperästä sekä kehityskoneella että Pi:llä.
- Rajataan riippuvuuksia ja muistinkulutusta yhdelle asemalle ja pienelle ennustejaksolle.

**Non-Goals:**

- Ei yleiskäyttöistä säätietokehystä tai varautumista kaikkiin tuleviin aineistoihin.
- Ei varsinaista Pi-asennusta, SD-kortin alustusta tai käyttöjärjestelmän ylläpitoa tämän suunnittelupyynnön perusteella.
- Ei täyttä tuotannon julkaisu-/palautus-/systemd-ratkaisua. Pi:n manuaalisen ajon ohje kuuluu tähän muutokseen; varsinainen laitehyväksyntä siirretään [issueen #1](https://github.com/ottotur/weather-station/issues/1) ja ylläpidon viimeistely myöhempään muutokseen.

## Decisions

### 1. Yksi Go-prosessi ja erillinen staattinen käyttöliittymä

Ehdotettu rakenne toteutukseen: `cmd/weather-station`, `internal/fmi`, `internal/waterlevel`, `internal/cache`, `internal/httpapi` ja `web/`. FMI-osa omistaa XML/WFS-tulkinnan ja HTTP-haun; vedenkorkeusosa normalisoinnin, tuoreuden ja yhteenvedot; välimuisti pysyvyyden; HTTP-osa näkymäsopimuksen. Kellonaika, lähdeosoite ja välimuistihakemisto voidaan antaa testeissä erikseen.

Svelte + TypeScript + Vite rakennetaan staattiseksi paketiksi, jonka Go palvelee paikallisesta hakemistosta. Yksi julkaisu koostuu Go-ohjelmasta ja täsmälleen saman version käyttöliittymäpaketista. Tämä on yksinkertaisempi kuin erillinen Node-palvelin Pi:llä. Go:n standardikirjastoa suositaan HTTP:hen, XML:ään, JSON:iin ja testeihin. SQLitea tai Dockeria ei tarvita.

### 2. Rajattu haku ja tiukka lähdetulkinta

Havaintokysely: `fmi::observations::mareograph::instant::timevaluepair`, `fmisid=134254`, `parameters=WLEVN2K_PT1S_INSTANT`, `timezone=UTC`, `timestep=1`, liukuva viimeinen tunti. Säilytetään vain viimeisin äärellinen kelvollinen havainto; historia ei kerry.

Ennustekysely: `fmi::forecast::sealevel::point::timevaluepair`, sama FMISID, UTC ja `timestep=60`, nykyhetkestä enintään 48 h. `parameters`-rajausta ei lisätä, koska sitä ei dokumentoitu. Vastauksesta valitaan vain `SeaLevelN2000`; `SeaLevel` ei ole epävarmuusraja tai korvaava N2000-sarja. Lyhyempi kattavuus hyväksytään näkyvänä puutteena, ei haeta toista lähdettä.

Parseri käyttää nimiavaruuksia, tunnistaa WFS-virheet myös HTTP 200:ssa ja ratkaisee vain vastauksen sisäiset `xlink:href`-aikaviitteet. Ulkoisia XML-viitteitä ei haeta. Tarkistetaan tunnisteet, odotettu suure, UTC-ajat, yksikäsitteiset järjestettävät aikapisteet ja ennusteen `analysisTime`. Puuttuva tai virheellinen vaadittu viiteaika hylkää uuden ennusteen; hakuaika ei ole varaviiteaika. Tulevaisuudessa oleva havainto tai mallin viiteaika hylätään selkeänä aikavirheenä; ennustepisteiden tulevat ajat ovat normaaleja. Täysin äärellisiä pisteitä vailla oleva sarja ei korvaa vanhaa. Äärellistä mutta jo vanhaa aineistoa voidaan käyttää vanhentuneeksi merkittynä.

`NaN`/puuttuvat arvot normalisoidaan `null`-aukoiksi. Muu virheellinen lukuteksti hylkää vastauksen. Tavoitteeksi asetetaan 30 s HTTP-aikakatkaisu ja 2 MiB vastausraja; ne ovat toteutuksen lähtöarvoja, eivät FMI-faktoja. Samalle aineistolle ei sallita päällekkäisiä hakuja. Ei automaattisia välittömiä HTTP-uusintoja. Palvelimen nimenomaisia käyttörajoitusohjeita noudatetaan, eikä 429-virhettä kierretä tiheämmällä haulla.

**Yksikköriski:** havainto jaetaan kymmenellä; ennustearvo jätetään sellaisekseen cm-oletuksella. Oletuksen alkuperä tallennetaan aineistoon (`unitStatus: assumed-cm`), API:in ja käyttöliittymän selitteeseen. Ristiriitainen uusi FMI-metatieto on pysäyttävä tulkintavirhe ja vaatii käyttäjälle palautettavan korjauspäätöksen, ei hiljaista skaalausta.

### 3. Ajat ja tilat eivät sekoitu toisiinsa

Tuoreus lasketaan palvelimella jokaiselle API-vastaukselle, ei vain haun yhteydessä. Havainto on vanhentunut, kun `now - observedAt > 30 min`; ennuste, kun `now - analysisTime > 12 h`. Täsmälleen rajalla aineisto on vielä tuore. Tämä täsmentää hyväksyttyä rajan ylittämisen käytäntöä.

Kummallakin aineistolla on toisistaan riippumattomat tiedot:

- aineiston saatavuus ja tuoreus: `missing | fresh | stale`;
- `lastAttemptAt`, yrityksen tulos `pending | success | error` ja rajattu virhekoodi;
- viimeisen julkaistun aineiston `fetchedAt`, havainto-/viiteaika, erillinen `resultTime` ja ennusteen viimeinen äärellinen voimassaolohetki;
- ennusteen `hasFutureValues` ja näkyvä kattavuus. Vanha malliajo voi edelleen kattaa tulevaisuutta; tuoreuden ja kattavuuden tiloja ei yhdistetä yhdeksi epämääräiseksi tilaksi.

Mallin viiteaika ei ole varmennettu julkaisuaika. UI nimeää sen mallin viiteajaksi. Tuoretta välimuistia ja hakuvirhettä voidaan näyttää samanaikaisesti. API ei nojaa selaimen kelloon eikä selain laske tuoreutta laitteen aikavyöhykkeestä.

### 4. Julkaisu vasta pysyvän tallennuksen jälkeen

Kummallekin aineistolle oma versioitu JSON-tiedosto konfiguroidussa välimuistihakemistossa. Tallennetaan normalisoitu aineisto, lähdeajat, suure, yksikkötulkinta ja viimeinen onnistunut hakuaika. Ei jatkuvaa raaka-XML-arkistoa eikä epäonnistumisella päivitetä onnistuneen aineiston aikaleimoja.

Tallennus kirjoittaa väliaikaistiedostoon samassa hakemistossa, synkronoi ja sulkee sen, korvaa kohteen atomisesti ja synkronoi hakemiston Linuxissa. Vasta onnistuneen hakemistosynkronoinnin jälkeen päivitetään muistissa julkaistu snapshot. Ennen korvausta tapahtuvassa virheessä aiempi kohdetiedosto ja muistissa julkaistu aineisto säilyvät. Jos korvaus onnistuu mutta hakemiston synkronointi epäonnistuu, muistissa säilyy aiempi aineisto ja häiriö ilmoitetaan erillisellä tallennuksen epävarman tilan syykoodilla: uusi tiedosto voi jo olla levyllä ja tulla luetuksi uudelleenkäynnistyksessä. Ei luvata palauttaa aiempaa tiedostoa korvauksen jälkeen. Toinen aineisto toimii riippumatta häiriöstä. Tämä vähentää kesken jääneen kirjoituksen riskiä mutta ei takaa SD-kortin fyysistä ehjyyttä.

Käynnistyksessä luetaan ja tarkistetaan tiedostot erikseen, sivuutetaan keskeneräiset väliaikaistiedostot ja aloitetaan molempien lähteiden haku. Vioittunut tai yhteensopimaton tiedosto merkitsee vain kyseisen aineiston puuttumista; prosessi ei kaadu. Viimeisen epäonnistuneen yrityksen tilaa ei tarvitse säilyttää käynnistysten yli. Tietokanta ja yhteinen molemmat aineistot kerralla korvaava tiedosto hylätään tarpeettomana kytkentänä.

### 5. Rajapinta ja laskenta

`GET /api/water-level/hamina` palauttaa HTTP 200 myös osittaisessa/tyhjässä datatilassa: kyse on onnistuneesta tilakyselystä. Odottamaton palvelinvirhe käyttää 5xx-vastausta. `Cache-Control: no-store` ehkäisee selaimen vanhan tilavastauksen uudelleenkäyttöä.

Sopimuksen kenttäryhmät: `schemaVersion`, `serverTime`, `station`, `unit`, `datum`, `observation`, `forecast`, `summary`, `source`. Aineistojen kentät kattavat kohdassa 3 mainitut tilat ja ajat. Puuttuva luku on `null`, ei nolla tai JSON:n vastainen `NaN`. Virheet ovat käyttäjälle sopivia koodeja; sisäisiä polkuja tai raakavirheitä ei lähetetä selaimeen.

Palvelin valitsee tulevat äärelliset pääennustepisteet väliltä `[serverTime, serverTime + 48 h]`. Ei interpolaatiota nykyhetkeen eikä extrapolaatiota. Suurin nousu = `max(0, max(forecast) - observation)`; suurin lasku = `max(0, observation - min(forecast))`. Nollamuutos nimetään "ei nousua/laskua havaintotasosta", ei tekaistua äärihetkeä. Tasatuloksissa valitaan aikaisin hetki. Lasketaan ennen näyttöpyöristystä; esitetään cm yhdellä desimaalilla yhtenäisesti.

Jos sisäisiä aukkoja on, yhteenveto koskee vain saatavilla olevia arvoja ja sen vieressä näkyy kattavuusvaroitus. Jos havainto tai ennuste ei ole tuore, tai tulevia äärellisiä pisteitä ei ole, `summary` on puuttuva ja mukana on syykoodi. API laskee yhteenvedon yhdestä johdonmukaisesta snapshotista. Selain ei toteuta rinnakkaista laskentalogiikkaa.

### 6. Kevyt mobiiliesitys

Yhden aseman käyrä toteutetaan aluksi kevyenä SVG:nä ilman kaaviokirjastoa. Kosketus/kohdistus valitsee pisteen; avattava semanttinen tuntitaulukko tarjoaa samat arvot ilman kaavion käyttöä. Näin vältetään suuren kirjaston kustannus ensimmäiselle yksinkertaiselle aikasarjalle. Kaaviota ei siloteta eikä piirretä aukon yli. Havainto erotetaan ennusteesta nimellä, ajalla ja esitystavalla.

Svelte hakee API:n noin minuutin välein ja `visibilitychange`-palautumisessa ilman rinnakkaisia pyyntöjä. Yhteysvirhe piilottaa muutossummat, merkitsee viimeisen selaimessa olevan aineiston yhteyskatkolla ja yrittää uudelleen seuraavalla normaalilla päivityksellä. Ei service workeria, localStorage-offlinevälimuistia tai CDN-riippuvuuksia.

Ajat muotoillaan `Intl.DateTimeFormat`-toiminnolla kiinteästi `Europe/Helsinki`-aikaan. Kesä-/talviajan toistuvat kellonajat erotetaan UTC-poikkeamalla. Tekstit, kontrasti, näppäimistökäyttö, kosketuskohteet ja 360 px leveys tarkistetaan. FMI-lähdemerkintä, CC BY 4.0 -linkki, lähdekyselylinkit ja muunnos-/laskentaselite tarjotaan paikallisesti renderöityinä; linkkien avaaminen on eri asia kuin sovelluksen ulkoinen verkkoriippuvuus.

### 7. Kaksi erillistä testaus- ja hyväksyntätasoa

**A — kehityskoneella toimiva kokonaisuus:** Go ja Node-työkalut ajetaan Ubuntu WSL:ssä, projekti WSL:n Linux-tiedostojärjestelmässä. Tuetut työkalujen versiot lukitaan toteutuksen alussa. Palvelin sidotaan paikalliseen osoitteeseen; rakennettu käyttöliittymä avataan Windowsin selaimessa `http://localhost:<portti>`. Ohje erottaa Vite-kehitysajon rakennettujen staattisten tiedostojen koko ketjun testistä. WSL:n localhost-välitys tarkistetaan, ei oleteta toimivaksi eikä avata palomuuria tai lähiverkkoa automaattisesti.

- Go-testit: paikallinen HTTP-valepalvelin, injektoitu kello, tilapäinen välimuisti. Alkuperäisiä XML-näytteitä ei muuteta; poikkeustestit käyttävät erillisiä synteettisiä fixtureitä. Testikello asetetaan vuoden 2026 näytteisiin sopivaksi, jotta ne eivät kelpaa vahingossa nykyhetken live-dataksi.
- Käyttöliittymä: komponentti-/tilatestit ja Playwright-selaintestit. Koko ketjun testi käyttää oikeaa Go-palvelinta ja paikallista WFS-valepalvelinta, ei vain valmista API-mockia.
- Tapaukset: mm/cm, sisäinen aikaviite, väärä asema/suure, aikavirheet, `NaN`/aukot, tuoreusrajan molemmat puolet ja tasaraja, sama malliajo uusintahakuna, hakuvirhe tuoreella datalla, vanha havainto/ennuste erikseen, ennusteen päättyminen, tyhjä käynnistys, vioittunut välimuisti, kirjoitusvirhe, uudelleenkäynnistys, kummankin palvelinyhteyden katko ja toipuminen.
- Rakennus, Go-testit/race-tarkistus soveltuvalla hostilla, `go vet`, TypeScript/Svelte-tarkistus ja käyttöliittymätestit; `linux/arm64`-ristiinkäännös asetuksilla `GOARM64=v8.0 CGO_ENABLED=0`. Linux/x86_64-kehitysajo ja arm64-käännös kirjataan erikseen.
- Erillinen manuaalinen Windows-selaintesti, kapea näkymä, tuntitaulukko, kesä-/talviaika sekä selaimen verkkopyyntöjen tarkistus. Automaattinen WebKit ei ole oikean iPhonen Safari-hyväksyntä.
- Valinnainen rajattu live-smoke: dokumentoidut kyselyt ja ehdot tarkistetaan, yksi haku kumpaankin aineistoon, normaalin ajastuksen rinnakkaisajo estetään. Tarkistetaan todellinen data, mutta live-FMI ei ole automaattisten testien riippuvuus.

**B — Pi:llä varmennettu käyttöönotto:** kun laite on valmis ja käyttäjä sallii ajon, sama rakennettu versio ajetaan hyväksytyllä Trixie/arm64-kohteella. Tarkistetaan oikea arkkitehtuuri, käynnistys, haku, pysyvä välimuisti, internetkatko ja prosessin uudelleenkäynnistys. Testataan oikealla mobiililaitteella Safari sekä Android Chrome/Firefox käytettävissä olevilla laitteilla; testaamattomat yhdistelmät kirjataan avoimiksi.

Mitataan käyttökelpoisen yhteenvedon 2 s tavoite valmiilla välimuistilla sekä 10 yhtäaikaisen asiakkaan käyttö, ja kirjataan laite, selain, verkko ja mittaustapa. WSL-tulokset eivät todista Pi:n suorituskykyä. Käyttäjän hyväksymänä rajauksena taso B ja sitä vastaavat tekemättömät tehtävät siirretään avoimeen [issueen #1](https://github.com/ottotur/weather-station/issues/1); taso A voidaan hyväksyä tässä muutoksessa. Siirto ei tarkoita laitetestien suorittamista tai laitehyväksyntää.

## Risks / Trade-offs

- [Ennusteen cm on hyväksytty oletus, ei suora metatietovarmennus] → Tulkinta näkyy kaikissa kerroksissa; ristiriita pysäyttää uuden tulkinnan julkaisun.
- [Vanha dokumentoitu julkaisurytmi ei vastaa näytteen malliajoa] → Vanhentumisrajat ovat käyttäjän hyväksymiä sovellusrajoja, eivät FMI-palvelulupauksia; sama malliajo ei nuorene haettaessa.
- [Vain tunnin havaintohaku] → Pitkässä lähdeviiveessä uutta havaintoa ei ehkä saada, mutta vanha tarkistettu havainto säilyy vanhentuneena. Kattavuutta ei laajenneta huomaamatta historiakeruuksi.
- [WSL ja Pi eroavat] → Natiivi Linux-kehitysajo, arm64-käännös ja oikea Pi-ajotesti pidetään erillisinä näyttöinä.
- [SVG:n oma kosketuskäsittely] → Tuntitaulukko tarjoaa saman tiedon; todellinen mobiilitestaus ei korvaudu desktop-emulaatiolla.
- [Sähkökatko tai SD-vika] → Atominen korvaus, uudelleenkäynnistystestit ja näkyvä vioittumistila; ei luvata fyysisen tallennusmedian vikasietoisuutta.

## Migration Plan

Ei olemassa olevaa sovellusta tai käyttäjätietomigraatiota. Toteutus aloitetaan vain uudesta apply-pyynnöstä. Tässä muutoksessa toteutetaan ja testataan taso A; taso B suoritetaan myöhemmin käyttäjän sallimalla laitteella [issuen #1](https://github.com/ottotur/weather-station/issues/1) mukaisesti. Sovelluksen oletuskuuntelu on localhost; Pi:n lähiverkkokuuntelu asetetaan erikseen, ei internetiin avaamista.

Ensimmäistä manuaalista Pi-ajoa varten ohjeistetaan erillinen version hakemisto ja välimuistihakemisto. Testiversion voi pysäyttää koskematta asetuksiin tai aiempaan julkaisuun. Tämä ei korvaa myöhempää tuotannon julkaisu- ja palautusmenettelyä.

## Open Questions

- WSL:n tarkka versio ja paikallisesti asennettujen työkalujen versiot tarkistetaan toteutuksen alussa; lähestymistapa ei edellytä tiettyä WSL-verkkotilaa.
- Pi:n ja oikeiden mobiililaitteiden saatavuus sekä mallikohtaiset mittaustulokset kirjataan vasta ajotestissä. Tekemättömälle testille ei anneta hyväksyttyä tulosta.
