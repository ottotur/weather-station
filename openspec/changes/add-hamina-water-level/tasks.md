# Tasks

Tämä on tulevan toteutuksen tehtävälista, ei lupa aloittaa toteutusta. Taso A tarkoittaa kehityskoneella toimivaa kokonaisuutta; taso B oikealla Pi:llä ja mobiililaitteilla varmennettua käyttöä. Tekemättömät laitetestit jätetään avoimiksi.

## 1. Sovellusrunko ja paikallinen kehitys

- [ ] 1.1 Tarkista tuetut Go-, Node- ja käyttöliittymätyökalujen versiot sekä Ubuntu WSL:n versio ja arkkitehtuuri. Kirjaa valinnat kehitysohjeeseen; varmista versiot työkalujen omilla komennoilla ja päivitä selvityksen Go-version sijasta ajantasainen tuettu korjausversio.
- [ ] 1.2 Luo Go-palvelimen ja Svelte + TypeScript + Vite -käyttöliittymän runko designin mukaisiin hakemistoihin. Varmista, että Go-testi, käyttöliittymän tyyppitarkistus ja staattinen rakennus onnistuvat, eikä generoituja riippuvuushakemistoja lisätä Gitiin.
- [ ] 1.3 Lisää asetukset kuunteluosoitteelle, välimuistihakemistolle ja staattisille tiedostoille sekä testikäyttöön kellolle/lähdeosoitteelle. Dokumentoi WSL:n Linux-tiedostojärjestelmässä ajo ja localhost-oletus; varmista, että virheelliset asetukset tuottavat selkeän käynnistysvirheen.

## 2. FMI-aineistojen tulkinta ja normalisointi

- [ ] 2.1 Toteuta havaintojen parseri ja mm → cm -normalisointi. Varmista alkuperäisellä näytteellä asema, suure, viimeisin äärellinen arvo 32,1 cm ja mittausaika; lisää testit väärälle asemalle/suureelle, tyhjille arvoille ja aikavirheille.
- [ ] 2.2 Toteuta `SeaLevelN2000`-ennusteen parseri, sisäisten aikaviitteiden ratkaisu sekä viite-/tulos-/voimassaoloaikojen erottelu. Varmista näytteellä 114 äärellistä pistettä ja seitsemän puuttuvaa loppupistettä sekä testaa sisäiset aukot, puuttuva viiteaika ja `SeaLevel`-korvauksen estyminen.
- [ ] 2.3 Lisää yksikkötulkinnan alkuperä ja tarkistussäännöt aineistomalliin. Testaa cm-oletuksen säilyminen ja virheellisen XML:n, WFS-poikkeuksen, kaikkien arvojen puuttumisen, väärän lukutekstin ja ristiriitaisen yksikkötiedon hylkäys. Dokumentoi tulkinta muuttamatta alkuperäisiä XML-näytteitä.

## 3. Tuoreus ja muutoksen laskenta

- [ ] 3.1 Toteuta kellosta riippuvat 30 min / 12 h -tilat sekä kattavuuden erottaminen tuoreudesta. Varmista testeillä rajan alitus, tasaraja, ylitys, sama malliajo uusintahaussa ja päättynyt ennustejakso; dokumentoi `analysisTime` mallin viiteajaksi, ei julkaisuajaksi.
- [ ] 3.2 Toteuta näytetyn tulevan 48 h jakson nousu-/laskuyhteenveto viimeisimmästä havainnosta. Testaa molemmat muutossuunnat, vain nousu/lasku, nollamuutos, tasatulosten aikaisin hetki, aukkojen varoitus sekä laskelmien piilottaminen kummankin aineiston vanhentuessa erikseen.

## 4. Pysyvä välimuisti

- [ ] 4.1 Toteuta erilliset versioidut havainto- ja ennustetiedostot, tarkistettu lataus käynnistyksessä ja atominen korvaus ennen muistijulkaisua. Varmista tilapäishakemistotesteillä onnistunut uudelleenkäynnistys ja lähde-/yksikköaikojen säilyminen.
- [ ] 4.2 Testaa kirjoitusvirhe, katkennut kirjoitus, jäljelle jäänyt väliaikaistiedosto, vioittunut tiedosto ja tuntematon tiedostoversio. Varmista aiemman aineiston säilyminen ja toisen aineiston riippumaton käyttö; dokumentoi tiedostomuoto ja virheestä palautuminen.

## 5. Haku ja ajastus

- [ ] 5.1 Toteuta rajatut dokumentoidut FMI-kyselyt, aikakatkaisu ja vastauskokoraja. Varmista paikallisella HTTP-valepalvelimella kyselyparametrit, HTTP-/WFS-virheet ja ylikokoisen/hitaan vastauksen hallittu hylkäys; dokumentoi lähteet, käyttörajat ja CC BY 4.0.
- [ ] 5.2 Yhdistä käynnistyshaku ja riippumattomat 10 min / 30 min ajastukset tarkistukseen ja pysyvään julkaisuun. Varmista testikellolla virheen jälkeinen normaali hakuväli, päällekkäisyyksien esto, vanhan snapshotin säilyminen sekä ettei asiakasmäärä muuta FMI-hakumäärää.

## 6. Paikallinen rajapinta

- [ ] 6.1 Toteuta `GET /api/water-level/hamina`, palvelinkohtainen tuoreuslaskenta ja yhtenäinen snapshot. Varmista HTTP-testeillä täysi/osittainen/tyhjä vastaus, `null`-arvot, hakuyrityksen erottaminen tuoreudesta, syykoodit, lähdemerkintä ja `Cache-Control: no-store`; dokumentoi JSON-sopimus esimerkein.
- [ ] 6.2 Lisää rakennettujen staattisten tiedostojen tarjoilu samasta alkuperästä ja varmista integraatiotestillä käyttöliittymän latautuminen ilman Vite-palvelinta. Testaa kymmenen yhtäaikaista API-asiakasta ilman lisähakuja tai ristiriitaisia vastauksia.

## 7. Mobiilinäkymä

- [ ] 7.1 Toteuta havaintoyhteenveto, tuleva ennustejakso, aukot näyttävä SVG-kaavio ja avattava tuntitaulukko. Varmista komponentti-/selaintesteillä arvot, kosketus-/näppäimistövalinta, 360 px leveys ja sivun vaakavierityksen puuttuminen.
- [ ] 7.2 Toteuta nousu-/laskuesitys, lataus-, puute-, vanhentumis- ja hakuvirhetilat sekä Pi-yhteyskatko ja palautuminen. Testaa minuutin päivitys ja sivulle palaaminen, yhteenvetojen piilottuminen ja ettei vanha tieto näy tuoreena; dokumentoi tilojen merkitys käyttäjäohjeessa.
- [ ] 7.3 Lisää Suomen aika, toistuvien kellonaikojen erottaminen, lähde-/lisenssimerkinnät ja yksikköoletuksen selite. Testaa selaimen eri aikavyöhyke ja kesä-/talviajan vaihto sekä varmista, ettei näkymä tee FMI-/CDN-pyyntöjä tai esitä keksittyjä epävarmuusrajoja.

## 8. Kehityskoneen koko ketjun hyväksyntä — taso A

- [ ] 8.1 Lisää toistettava koko ketjun testiajo oikealla Go-palvelimella, paikallisella WFS-valepalvelimella ja rakennetulla käyttöliittymällä. Varmista hausta näkymään normaalitila, lähdekatko, vanheneminen, virheellinen vastaus, uudelleenkäynnistys ja toipuminen ilman live-FMI-riippuvuutta; dokumentoi täsmälliset ajokomennot.
- [ ] 8.2 Suorita Go-testit, soveltuva race-tarkistus, `go vet`, Svelte/TypeScript-tarkistus, käyttöliittymätestit ja rakennukset. Ristiinkäännä `GOOS=linux GOARCH=arm64 GOARM64=v8.0 CGO_ENABLED=0`; kirjaa tulokset erikseen host-ajosta äläkä merkitse Pi-ajotestiä tehdyksi.
- [ ] 8.3 Varmista rakennettu sovellus Windows 11:n selaimessa WSL:n localhost-osoitteesta. Kirjaa selain, WSL-versio, kosketusemulaation rajat ja näkyvät häiriötilat hyväksyntäraporttiin; dokumentoi mahdollinen localhost-välityksen este ilman automaattisia palomuuri-/verkkomuutoksia.

## 9. Oikean laitteen hyväksyntä — taso B

- [ ] 9.1 Laadi manuaalisen Pi-ajon ohje hyväksytylle Raspberry Pi OS Lite 64-bit / Debian 13 Trixie -kohteelle: version tiedostot, erillinen välimuisti, kuunteluosoite, käynnistys ja pysäytys. Varmista ohjeen tiedostopolut ja komennot; älä alusta tai asenna laitetta ilman erillistä lupaa.
- [ ] 9.2 Kun Pi on valmis ja ajo sallittu, suorita sama versio oikealla `linux/arm64`-laitteella. Kirjaa arkkitehtuuri, käynnistys, rajatut live-haut tarkistettujen käyttöehtojen mukaisesti, internetkatko, välimuistin säilyminen ja uudelleenkäynnistys. Jätä tehtävä avoimeksi, jos laitetta ei voi testata.
- [ ] 9.3 Testaa oikeiden mobiililaitteiden Safari-, Android Chrome- ja Firefox-käyttö soveltuvin laittein. Mittaa käyttökelpoisen yhteenvedon 2 s tavoite valmiilla välimuistilla ja 10 samanaikaisen asiakkaan käyttö; kirjaa laitteet, menetelmä, tulokset ja testaamattomat yhdistelmät. Älä korvaa puuttuvaa näyttöä WSL- tai emulaatiotuloksilla.

## Workflow follow-up

- Esitä tasojen A ja B tila erikseen; älä väitä koko MVP:tä tai Pi-käyttöönottoa hyväksytyksi.
- Pyydä erillinen hyväksyntä ennen yhdistämistä `main`iin. Arkistointicommit `b50b789` on eri haarassa eikä sisälly tähän muutokseen.
- Arkistoi tämä muutos vasta toteutuksen ja vaadittujen tarkistusten valmistuttua tai käyttäjän hyväksyttyä avoimet poikkeamat nimenomaisesti.
