# Tasks

Tehtävät koskevat vain selvityksen toteuttamista erillisen apply-pyynnön jälkeen. Ei sovelluskoodia, palvelinasennusta tai verkkomuutoksia. Delta-speksit on ohitettu tarkoituksellisesti; hyväksymisehdot ovat `design.md`:n kohdassa 6.

## 1. Kohteet ja lähdekäytön edellytykset

- [x] 1.1 Aloita `docs/feasibility/weather-wind.md`: kirjaa käyttäjän valitsemat testimökkipiste 60.434318, 26.687825 ja testiulappapiste 60.414616, 26.697395 sekä valintojen alkuperä ja selvityksen rajaus. Varmista, ettei raportti esitä testipisteitä lopullisina sijainteina tai edellytä ulappapisteen uutta valintaa.
- [x] 1.2 Tarkista FMI:n ajantasaiset dokumentoidut sääennuste- ja tuulihavaintokyselyt, koordinaattikäytäntö, käyttöehdot, hakurajat ja lähdemerkintä. Kirjaa URL:t ja tarkistuspäivät raporttiin; varmista sallitut parametrit ennen aineistonäytehakuja.
- [x] 1.3 Selvitä enintään kaksi tuulihavaintoasemaehdokasta virallisten sijainti-/asematietojen avulla. Kirjaa sijainnit, asematunnisteet, etäisyydet ja edustavuuden perustelut raporttiin; tarkista, ettei lähintä asemaa nimetä automaattisesti edustavimmaksi.
- [x] 1.4 Esitä käyttäjälle suositeltu havaintoasema ja mahdollinen vaihtoehto. Kirjaa hyväksytty valinta ja hyväksynnän alkuperä raporttiin ennen aseman aineistonäytteitä; pysähdy aseman osalta, jos hyväksyntää ei saada. Testimökki- ja testiulappapiste eivät tarvitse uutta sijaintihyväksyntää.

## 2. Testimökkipisteen sää- ja tuuliennuste

- [x] 2.1 Valitse dokumentoitu lähde ja tee testimökkipisteelle rajattu enintään 48 h ennustehaku käyttöehtojen tarkistuksen jälkeen. Tallenna alkuperäinen vastaus ja URL, UTC-hakuaika, HTTP-tulos sekä SHA-256 hakemiston `docs/feasibility/evidence/weather-wind/` manifestiin. Varmista vastausrakenne ja pyydetyn sekä palautetun sijainnin vastaavuus.
- [x] 2.2 Täytä raporttiin lämpötilan, keskituulen, puuskien, suunnan, sääsymbolin ja tuntisateen matriisi sekä suoran näennäisen lämpötilan ja sateen todennäköisyyden saatavuus. Varmenna kentät, yksiköt, puuttuvat arvot ja dokumentoidut merkitykset näytteestä ja lähteistä; merkitse puutteet ilman omia korvaavia laskentoja.
- [x] 2.3 Tarkista näytteen ajat, askel, aukot ja todellinen äärellinen ennustekattavuus. Dokumentoi sateen kertymäjakso ja Suomen kalenterivuorokausisumman edellytykset; varmista, ettei hakuajasta tehdä julkaisuaikaa tai kertymästä intensiteettiä. Lisää näiden havaintojen offline-tarkistusohje raporttiin ja suorita se.

## 3. Hyväksytyn ulappapisteen ja aseman tuuli

- [x] 3.1 Tee käyttäjän valitsemalle testiulappapisteelle 60.414616, 26.697395 rajattu enintään 48 h tuuliennustehaku. Tallenna muuttamaton vastaus ja manifestitiedot; tarkista pyydetty ja palautettu mallipiste, keskituuli, puuskat, tulosuunta, yksiköt, lähdeajat, aukot ja kattavuus. Lisää tulokset raporttiin ja varmista ne offline-näytteestä; palauta mahdollinen sijaintiin liittyvä este käyttäjälle vaihtamatta pistettä hiljaisesti.
- [x] 3.2 Tee hyväksytylle havaintoasemalle rajattu havaintohaku aluksi enintään viimeiseltä tunnilta. Tallenna vastaus ja manifestitiedot; varmista aseman tunniste, viimeisin saatavilla oleva äärellinen arvo suureittain, mittausajat ja tuulisuureiden keskiarvoistusjaksot. Kirjaa puutteet raporttiin; älä korvaa puuttuvaa puuskaa laskennallisella arvolla.
- [x] 3.3 Kokoa aineistokohtaiset dokumentoidut päivitysrytmit ja ehdotus myöhemmistä haku-/vanhentumisrajoista perusteluineen. Varmista raportista, että ehdotukset erotetaan hyväksytyistä päätöksistä ja ettei vedenkorkeuden rajoja ole kopioitu ilman sääaineiston perusteita.

## 4. Jatkopäätös ja kokonaisuuden tarkistus

- [x] 4.1 Kirjaa raporttiin myöhemmän Linux-VM-/Proxmox-testauksen rooli ja rajat ilman asennustoimia tai alustakohtaisia käyttöohjeita. Varmista, että oikean mobiiliselaimen testi, VM-suorituskyky, arm64-käännös ja Pi:n laitehyväksyntä erotetaan ja issue #1 säilyy tekemättömän laitehyväksynnän viitteenä.
- [x] 4.2 Laadi lähdenäytön perusteella suositus yhdestä pienestä seuraavasta sovellustoiminnosta. Tarkista designin hyväksymisehdot ja kirjaa varmennetut tiedot, puutteet sekä käyttäjälle palautettavat rajauspäätökset; älä käynnistä sovellustoteutusta tai hyväksy ydintiedon korvausta.
- [x] 4.3 Suorita raportin offline-toistamisohje, tarkista kaikkien näytteiden rakenne ja SHA-256-tiivisteet sekä raportin paikalliset linkit. Varmista diffistä, ettei aiempia vedenkorkeusnäytteitä tai sovelluskoodia ole muutettu, ja validoi OpenSpec-muutos. Kirjaa toteutuneet tarkistukset ja puuttuva näyttö raporttiin.

## Workflow follow-up

- Pyydä selvityksen jälkeen erillinen päätös seuraavan sovellustoiminnon OpenSpec-suunnittelusta.
- Toimita muutos `main`iin GitHub-PR:n kautta; yhdistä vasta käyttäjän nimenomaisella luvalla ja vaadittujen tarkistusten onnistuttua.
- Arkistoi selvitys vasta sen toteutuksen ja katselmoinnin jälkeen. Suunnitteluasiakirjojen valmistuminen ei merkitse selvitystehtäviä tehdyiksi.
