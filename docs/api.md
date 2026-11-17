# Paikallinen vedenkorkeusrajapinta

`GET /api/water-level/hamina` vastaa JSON-tilakuvalla myös silloin, kun aineistoa
ei ole. Vastauksen `Cache-Control` on `no-store`. Selain ei hae FMI:ltä.

Esimerkki tyhjästä havainnosta:

```json
{
  "status": "missing",
  "data": null,
  "attempt": {"lastAttemptAt": null, "result": "pending", "errorCode": null}
}
```

Vastauksen juuritasolla ovat `schemaVersion`, `serverTime`, `station`, `unit`,
`datum`, `observation`, `forecast`, `summary`, `summaryReason` ja `source`.
Havainnon `data` sisältää `valueCm`, `observedAt`, `resultTime` ja `fetchedAt`.
Vanhoista välimuistitiedostoista, joihin tulosaikaa ei tallennettu, luetun
havainnon `resultTime` on `null`: sitä ei päätellä mittaus- tai hakuajasta.
Ennusteen
`data` sisältää `analysisTime`, `resultTime`, `fetchedAt`, `unitStatus` ja
viimeisen äärellisen arvon ajan `lastValidAt`. `analysisTime` on mallin
viiteaika, ei varmennettu julkaisuaika; `unitStatus: assumed-cm` on käyttäjän
hyväksymä yksikköoletus.

`forecast.points` sisältää enintään tulevan 48 tunnin pisteet. Puuttuva arvo
on `null`. `lastShownAt` on näytetyn jakson viimeinen äärellinen aika ja
`hasFutureValues` kertoo, onko tulevia arvoja. Tuoreus (`fresh`, `stale`,
`missing`) ja viimeisen yrityksen tulos ovat toisistaan riippumattomia.

Esimerkiksi tuore havainto voi säilyä käytössä hakuvirheestä huolimatta:

```json
{
  "status": "fresh",
  "data": {"valueCm": 32.1, "observedAt": "2026-11-16T18:00:00Z", "resultTime": "2026-11-16T18:00:00Z", "fetchedAt": "2026-11-16T18:01:00Z"},
  "attempt": {"lastAttemptAt": "2026-11-16T18:10:00Z", "result": "error", "errorCode": "source_error"}
}
```

`cache_error` tarkoittaa tallennusvirhettä ennen pysyvää julkaisua ja
`cache_durability_uncertain` korvauksen jälkeistä hakemistosynkronoinnin
virhettä; jälkimmäisessä levyllä voi olla uusi tiedosto vaikka muistissa on
vanha aineisto. `unit_metadata_error` tarkoittaa, että lähteen ilmoittama
yksikkö on ristiriitainen tai tuntematon; aiempi tarkistettu aineisto säilyy.
`summary` on `null`, jos vertailua ei voida tehdä; syy on
`summaryReason`-kentässä. Muuten `rise` ja `fall` sisältävät `changeCm` ja
`at` (nollamuutoksella `at: null`), ja `hasGaps` varoittaa sisäisistä aukoista.
