<script lang="ts">
  import { onMount } from 'svelte';

  type Point = { time: string; valueCm: number | null };
  type Attempt = { result: string; errorCode: string | null };
  type Response = {
    serverTime: string;
    observation: { status: string; data: { valueCm: number; observedAt: string } | null; attempt: Attempt };
    forecast: { status: string; data: { analysisTime: string; unitStatus: string } | null; points: Point[]; hasFutureValues: boolean; lastShownAt: string | null; attempt: Attempt };
    summary: { rise: { changeCm: number; at: string | null }; fall: { changeCm: number; at: string | null }; hasGaps: boolean } | null;
    summaryReason: string;
  };

  let data: Response | null = $state(null);
  let disconnected = $state(false);
  let loading = $state(true);
  let selectedTime: string | null = $state(null);
  let pending = false;
  const date = new Intl.DateTimeFormat('fi-FI', {
    timeZone: 'Europe/Helsinki', day: 'numeric', month: 'numeric', hour: '2-digit', minute: '2-digit', timeZoneName: 'shortOffset'
  });
  const formatTime = (s: string) => date.format(new Date(s));
  const number = (n: number) => n.toLocaleString('fi-FI', { minimumFractionDigits: 1, maximumFractionDigits: 1 });
  const status = (s: string) => ({ fresh: 'Tuore', stale: 'Vanhentunut', missing: 'Puuttuu' }[s] ?? s);
  const attemptMessage = (attempt: Attempt) => ({
    cache_durability_uncertain: 'Tallennuksen pysyvyys on epävarma; aiempi tieto on käytössä.',
    unit_metadata_error: 'Lähteen yksikkötieto on ristiriitainen tai tuntematon; uutta aineistoa ei käytetä.'
  }[attempt.errorCode ?? ''] ?? 'Viimeisin haku epäonnistui');
  const summaryMessage = (reason: string) => ({
    observation_missing: 'Havainto puuttuu.',
    observation_stale: 'Havainto on vanhentunut.',
    forecast_missing: 'Ennuste puuttuu.',
    forecast_stale: 'Ennuste on vanhentunut.',
    no_future_values: 'Ennusteessa ei ole tulevia ennustearvoja.'
  }[reason] ?? 'Vertailuun tarvittavat tiedot eivät ole saatavilla.');

  async function refresh() {
    if (pending) return;
    pending = true;
    const controller = new AbortController();
    const timeout = window.setTimeout(() => controller.abort(), 10_000);
    try {
      const response = await fetch('/api/water-level/hamina', { cache: 'no-store', signal: controller.signal });
      if (!response.ok) throw new Error('Palvelinvirhe');
      const next = await response.json() as Response;
      if (selectedTime && !next.forecast.points.some(point => point.time === selectedTime && point.valueCm !== null)) {
        selectedTime = null;
      }
      data = next;
      disconnected = false;
    } catch {
      disconnected = true;
    } finally {
      window.clearTimeout(timeout);
      loading = false;
      pending = false;
    }
  }

  onMount(() => {
    void refresh();
    const timer = window.setInterval(refresh, 60_000);
    const visibility = () => { if (document.visibilityState === 'visible') void refresh(); };
    document.addEventListener('visibilitychange', visibility);
    return () => { clearInterval(timer); document.removeEventListener('visibilitychange', visibility); };
  });

  let points: Point[] = $derived.by(() => data?.forecast.points ?? []);
  let selected = $derived(points.find(point => point.time === selectedTime && point.valueCm !== null) ?? null);
  let finite = $derived(points.filter((p): p is Point & { valueCm: number } => p.valueCm !== null));
  let min = $derived(finite.length ? Math.min(...finite.map(p => p.valueCm)) : 0);
  let max = $derived(finite.length ? Math.max(...finite.map(p => p.valueCm)) : 1);
  const x = (point: Point, all: Point[]) => {
    if (all.length < 2) return 300;
    const first = Date.parse(all[0].time);
    const last = Date.parse(all[all.length - 1].time);
    return last === first ? 300 : 20 + (Date.parse(point.time) - first) * 560 / (last - first);
  };
  const y = (v: number, lo: number, hi: number) => 180 - (v - lo) / (hi - lo || 1) * 150;
  let segments = $derived.by(() => {
    const result: string[] = [];
    let segment: string[] = [];
    points.forEach((p) => {
      if (p.valueCm === null) { if (segment.length) result.push(segment.join(' ')); segment = []; }
      else segment.push(`${x(p, points)},${y(p.valueCm, min, max)}`);
    });
    if (segment.length) result.push(segment.join(' '));
    return result;
  });
</script>

<svelte:head><title>Haminan vedenkorkeus</title></svelte:head>
<main>
  <h1>Haminan vedenkorkeus</h1>
  <p>Hamina Pitäjänsaari · cm N2000</p>
  {#if loading}<p role="status">Ladataan tietoja…</p>{/if}
  {#if disconnected}<p role="alert">Yhteys paikalliseen palvelimeen katkesi. Näytössä on viimeksi saatu tieto.</p>{/if}
  {#if data}
    <section>
      <h2>Viimeisin havainto</h2>
      <p>{disconnected ? 'Viimeksi saatu tieto; tuoreutta ei voi tarkistaa' : status(data.observation.status)}{#if data.observation.attempt.result === 'error'} · {attemptMessage(data.observation.attempt)}{/if}</p>
      {#if data.observation.data}
        <p class="value">{number(data.observation.data.valueCm)} cm</p>
        <p>Mitattu {formatTime(data.observation.data.observedAt)}</p>
      {:else if data.observation.attempt.result === 'pending'}<p>Havaintoa haetaan…</p>
      {:else}<p>Havaintoa ei ole saatavilla.</p>{/if}
    </section>
    <section>
      <h2>Ennuste seuraaville 48 tunnille</h2>
      <p>{disconnected ? 'Viimeksi saatu tieto; tuoreutta ei voi tarkistaa' : status(data.forecast.status)}{#if data.forecast.attempt.result === 'error'} · {attemptMessage(data.forecast.attempt)}{/if}</p>
      {#if data.forecast.data}<p>Mallin viiteaika {formatTime(data.forecast.data.analysisTime)}</p>{/if}
      {#if data.forecast.lastShownAt}<p>Näytettävää ennustetta viimeistään {formatTime(data.forecast.lastShownAt)} asti.</p>{/if}
      {#if !data.forecast.data && data.forecast.attempt.result === 'pending'}<p>Ennustetta haetaan…</p>
      {:else if !data.forecast.hasFutureValues}<p>Voimassa olevaa ennustetta ei ole saatavilla.</p>{/if}
      {#if finite.length}
        <svg viewBox="0 0 600 200" role="img" aria-label="Ennustetun vedenkorkeuden kaavio">
          {#each segments as segment}<polyline points={segment} fill="none" stroke="currentColor" stroke-width="2" />{/each}
          {#each points as point}
            {#if point.valueCm !== null}
              <circle cx={x(point, points)} cy={y(point.valueCm, min, max)} r="7" fill="#125f86" tabindex="0"
                role="button" aria-label={`${formatTime(point.time)}, ${number(point.valueCm)} cm`}
                onclick={() => selectedTime = point.time} onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') selectedTime = point.time; }} />
            {/if}
          {/each}
        </svg>
        {#if selected && selected.valueCm !== null}<p>{formatTime(selected.time)}: {number(selected.valueCm)} cm</p>{/if}
      {/if}
      <details><summary>Tuntitaulukko</summary>
        <table><thead><tr><th>Aika Suomessa</th><th>cm N2000</th></tr></thead>
          <tbody>{#each points as point}<tr><td>{formatTime(point.time)}</td><td>{point.valueCm === null ? 'Puuttuu' : number(point.valueCm)}</td></tr>{/each}</tbody>
        </table>
      </details>
    </section>
    <section>
      <h2>Muutos viimeisimmästä havainnosta</h2>
      {#if !disconnected && data.summary}
        <p>Nousu: {data.summary.rise.changeCm ? `${number(data.summary.rise.changeCm)} cm` : 'ei nousua havaintotasosta'}{#if data.summary.rise.at} ({formatTime(data.summary.rise.at)}){/if}</p>
        <p>Lasku: {data.summary.fall.changeCm ? `${number(data.summary.fall.changeCm)} cm` : 'ei laskua havaintotasosta'}{#if data.summary.fall.at} ({formatTime(data.summary.fall.at)}){/if}</p>
        {#if data.summary.hasGaps}<p>Ennusteessa on aukkoja; yhteenveto koskee saatavilla olevia arvoja.</p>{/if}
      {:else}<p>Yhteenveto ei ole saatavilla. {disconnected ? 'Palvelinyhteys on poikki.' : summaryMessage(data.summaryReason)}</p>{/if}
    </section>
  {/if}
  <footer><p>Lähde: Ilmatieteen laitos, avoin data, <a href="https://creativecommons.org/licenses/by/4.0/">CC BY 4.0</a>.</p>
    <p><a href="https://opendata.fmi.fi/wfs?service=WFS&amp;version=2.0.0&amp;request=DescribeStoredQueries&amp;storedquery_id=fmi::observations::mareograph::instant::timevaluepair">Havaintokyselyn kuvaus</a> ·
      <a href="https://opendata.fmi.fi/wfs?service=WFS&amp;version=2.0.0&amp;request=DescribeStoredQueries&amp;storedquery_id=fmi::forecast::sealevel::point::timevaluepair">Ennustekyselyn kuvaus</a></p>
    <p>Havainnon millimetrit muunnetaan senttimetreiksi. Ennusteen cm-yksikkö on käyttäjän hyväksymä oletus, ei suoraan varmennettu metatieto. Muutokset on laskettu ennusteesta.</p>
  </footer>
</main>

<style>
  :global(body) { margin: 0; font-family: system-ui, sans-serif; color: #122b38; background: #f1f7fa; }
  main { max-width: 70rem; margin: auto; padding: 1rem; }
  section { background: white; padding: 1rem; margin: 1rem 0; border-radius: .5rem; }
  h1 { font-size: 1.7rem; }
  .value { font-size: 2rem; font-weight: bold; }
  svg { display: block; width: 100%; height: auto; color: #125f86; }
  circle { cursor: pointer; }
  circle:focus { stroke: #111; stroke-width: 3; }
  summary { cursor: pointer; padding: .75rem 0; }
  table { width: 100%; border-collapse: collapse; }
  th, td { padding: .5rem; border-bottom: 1px solid #aabcc5; text-align: left; }
</style>
