# Proposal

## Why

Haminan vedenkorkeustoiminnon kehityskoneketju on toteutettu, mutta seuraavien MVP-ominaisuuksien sää- ja tuulilähteet ovat vielä varmentamatta. Rajattu lähdeselvitys mahdollistaa seuraavan pienen kokonaisen toiminnon suunnittelun ilman lopullista mökkisijaintia tai käyttökuntoista Raspberry Pi:tä.

## What Changes

- Selvitetään FMI:n dokumentoiduista avoimista rajapinnoista testimökkipisteen **60.434318, 26.687825** sää- ja tuuliennusteiden saatavuus sekä erillisen ulappapisteen tuuliennuste ja soveltuvan lähiaseman tuulihavainnot.
- Testipiste on käyttäjän tätä selvitystä varten antama sijainti, ei lopullinen mökkipaikka. Ehdotetaan ulappapistettä ja tuulihavaintoasemaa perusteluineen ja pyydetään käyttäjän hyväksyntä ennen niiden varsinaisia aineistonäytehakuja.
- Kootaan lähteet, kyselyt, suureet, yksiköt, aikamerkitykset, todellinen kattavuus, puutteet ja käyttöehdot raporttiin `docs/feasibility/weather-wind.md`. Tallennetaan rajatut alkuperäiset julkiset vastaukset ja alkuperätiedot alihakemistoon `docs/feasibility/evidence/weather-wind/`.
- Suositellaan näytön perusteella yhtä pientä seuraavaa toteutusmuutosta, ei koko muun MVP:n toteutusta.
- Kirjataan Proxmoxissa ajettavan Linux-virtuaalikoneen mahdollinen rooli myöhemmässä integraatio- ja oikeiden mobiililaitteiden testauksessa. Pi:n laitehyväksyntä jää erilliseksi; käyttäjä aikoo kunnostaa laitteen vasta projektin valmistuttua.

## Capabilities

### New Capabilities

Ei uusia sovelluksen toiminnallisia vaatimuksia. Tämä on selvitys- ja dokumentaatiomuutos; `.openspec.yaml` sisältää `skip_specs: true`, kuten aiemmassa toteutettavuusselvityksessä.

### Modified Capabilities

Ei muutoksia vedenkorkeustoiminnon vaatimuksiin tai sovelluksen käyttäytymiseen. Selvityksen hyväksymisehdot määritellään `design.md`:ssä.

## Impact

- Vain suunnittelu, myöhemmän apply-vaiheen selvitysraportti ja lähdenäytteet. Ei muutoksia `cmd/`, `internal/` tai `web/` -toteutukseen, API:in tai ajonaikaisiin riippuvuuksiin.
- Ei aallokkoaineiston selvitystä, seitsemän päivän yleiskuvan toteutusta, uutta käyttöliittymää, asetustiedoston toteutusta tai muiden vedenkorkeusasemien lisäämistä.
- Ei Proxmoxin tai Pi:n asennusta, tunnusten keruuta, verkkomuutoksia tai palvelun avaamista. Proxmox ei korvaa issueen #1 kirjattua Pi:n suorituskyky- ja laitehyväksyntää.
- Ei FMI-livehakuja suunnitteluvaiheessa. Selvityksen toteutus käynnistyy erillisellä apply-pyynnöllä; lähdekäyttö rajataan käyttöehtojen tarkistuksen jälkeen.
