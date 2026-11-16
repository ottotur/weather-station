# Tasks

Tämä tehtävälista koskee vain selvityksen toteuttamista. Sovelluskoodia, palveluasennusta tai käyttöjärjestelmämuutoksia ei tehdä. Delta-speksit on tarkoituksellisesti ohitettu; menetelmä ja hyväksymiskriteerit ovat `design.md`:ssä.

## 1. Lähtötiedot ja laitetuki

- [ ] 1.1 Pyydä käyttäjältä Pi:n malli/revisio, suoritinarkkitehtuuri, käyttöjärjestelmäversio ja käyttäjätilan bittisyys tarvittaessa lukevien komentojen avulla. Kirjaa tiedot ja alkuperä raporttiin `docs/feasibility/water-level.md`; varmista, että puuttuvat tiedot on merkitty avoimiksi eikä projektiin tallennu tunnuksia tai yksityisiä verkkotietoja.
- [ ] 1.2 Tarkista laitteen, käyttöjärjestelmän ja Go:n virallisista tukitiedoista soveltuva yhdistelmä sekä käännöskohde ARM-asetuksineen. Kirjaa ehdotus, lähde-URL:t ja tarkistuspäivä raporttiin; varmista, että johtopäätös on jäljitettävissä lähteisiin ja erottaa dokumentoidun tuen tekemättömästä laiteajotestistä. Jos lähtötiedot tai tuki puuttuvat, kirjaa este.

## 2. Asema ja dokumentoidut aineistot

- [ ] 2.1 Pyydä käyttäjää valitsemaan ensimmäinen vedenkorkeusasema; tarjoa tarvittaessa FMI:n asematietoihin perustuvia vaihtoehtoja itäiseltä Suomenlahdelta. Kirjaa aseman nimi, lähteiden tunnisteet, valinnan peruste ja käyttäjän hyväksyntä raporttiin; varmista, ettei ehdotettua asemaa merkitä hyväksytyksi ilman vastausta.
- [ ] 2.2 Etsi valitun aseman FMI-havainnon ja pääennusteen dokumentoidut kyselyt sekä käyttöehdot, hakurajat ja lähdemerkintäohjeet. Kirjaa täsmälliset kyselyparametrit, lähde-URL:t ja tarkistuspäivä raporttiin; varmista ennen hakuja, että kyselyt ja suunniteltu rajattu käyttö noudattavat dokumentoituja ehtoja. Kirjaa puuttuva dokumentaatio näkyväksi esteeksi tai selvitettäväksi rajoitteeksi.

## 3. Todelliset vastaukset ja niiden tulkinta

- [ ] 3.1 Suorita hyväksytylle asemalle rajatut havainto- ja ennustekyselyt ja tallenna tarpeelliset alkuperäisvastaukset hakemistoon `docs/feasibility/evidence/`. Lisää manifestiin kyselyt, UTC-hakuajat ja dokumentaatiolähteet; varmista vastausten onnistuminen ja molempien aineistojen asematunnisteet. Virhetilanteessa tarkista kysely ja dokumentaatio sekä kirjaa havaittu este ilman perusteetonta päätelmää aineiston puuttumisesta.
- [ ] 3.2 Täytä raportin tarkistusmatriisiin suureet, yksiköt, N2000-osoitus, aikaleimojen merkitys, viimeisin havainto, ennusteen alku ja loppu, askel ja aukot sekä puuttuvien arvojen esitystapa. Varmista tiedot tallennetuista vastauksista ja lähdedokumenteista; vertaa todellista kattavuutta 48 tunnin tavoitteeseen äläkä tulkitse hakuaikaa julkaisuajaksi.
- [ ] 3.3 Tarkista pääennusteen ja mahdollisten epävarmuusrajojen merkitys sekä dokumentoitu julkaisurytmi. Kirjaa saatavuus, epäselvyydet ja vaikutus myöhempään tuoreuden arviointiin raporttiin; varmista, ettei epävarmuusrajoja nimetä absoluuttisiksi ääriarvoiksi eikä dokumentoimatonta julkaisurytmiä tai vanhentumisrajaa esitetä faktana.
- [ ] 3.4 Lisää raporttiin toistamisohjeet ja tarkista niiden kyselyt manifestia ja tallennettuja vastauksia vasten. Varmista, että tiedostot ja linkit löytyvät, vastaukset ovat alkuperäisiä ja mahdollinen uusintahaku rajataan käyttöehtojen mukaiseksi; kirjaa, että ajan mukana vaihtuvat arvot eivät ole virhe toistettavuudessa.

## 4. Yhteinen loppuarvio

- [ ] 4.1 Arvioi kerätty näyttö `design.md`:n jatkosuunnitteluehtoja vasten ja kirjoita raporttiin tila sekä puuttuvat tiedot. Varmista, että ydintarpeiden puutteet ja alustavien teknologiavalintojen hyväksynnät palautetaan käyttäjän päätettäviksi, eikä avoimia kohtia kuitata hiljaisesti ratkaistuiksi.
- [ ] 4.2 Tarkista raportti, lähdenäytteet ja manifesti yhtenä kokonaisuutena: jokaisella keskeisellä johtopäätöksellä on lähde, tarkistusaika ja tila, jokainen olennainen ristiriita näkyy ja ulkopuolelle jääneet vaiheen 0 työt on lueteltu. Varmista myös, ettei muutos sisällä sovelluskoodia tai väitä laiteajoa, koko vaiheen 0 valmistumista tai MVP:n hyväksyntää tapahtuneeksi.

## Workflow follow-up

- Esitä valmis selvitys ja sen mahdolliset päätöstarpeet käyttäjälle.
- Kun tarvittavat lähtötiedot ja hyväksynnät ovat kasassa, suunnittele yhden aseman kokonainen vedenkorkeustoiminto erillisenä OpenSpec-muutoksena käyttäjän uudesta pyynnöstä.
