import { test, expect } from '@playwright/test';
import { spawn, execFileSync, type ChildProcess } from 'node:child_process';
import { createServer, type Server } from 'node:http';
import { createServer as createNetServer } from 'node:net';
import { mkdtempSync, readFileSync, rmSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';

test('oikea Go-palvelin: haku, katko, vanheneminen, uudelleenkäynnistys ja toipuminen', async ({ page }) => {
  test.setTimeout(120_000);
  const root = resolve('..');
  const dir = mkdtempSync(join(process.env.WEATHER_E2E_TMPDIR ?? tmpdir(), 'weather-e2e-'));
  const binary = join(dir, 'weather-station');
  const samples = {
    observation: readFileSync(join(root, 'docs/feasibility/evidence/hamina-observation.xml')),
    forecast: readFileSync(join(root, 'docs/feasibility/evidence/hamina-forecast.xml'))
  };
  let mode: 'ok' | 'invalid' | 'offline' = 'ok';
  let child: ChildProcess | undefined;
  let mock: Server | undefined;
  try {
    execFileSync('go', ['build', '-o', binary, './cmd/weather-station'], { cwd: root });
    mock = createServer((req, res) => {
      if (mode === 'offline') { res.writeHead(503); res.end(); return; }
      if (mode === 'invalid') { res.end('<broken>'); return; }
      const query = new URL(req.url!, 'http://localhost').searchParams.get('storedquery_id') ?? '';
      res.end(query.includes('observations') ? samples.observation : samples.forecast);
    });
    await new Promise<void>(resolve => mock!.listen(0, '127.0.0.1', resolve));
    const sourcePort = (mock.address() as { port: number }).port;
    const allocate = createNetServer();
    await new Promise<void>(resolve => allocate.listen(0, '127.0.0.1', resolve));
    const appPort = (allocate.address() as { port: number }).port;
    await new Promise<void>((resolve, reject) => allocate.close(err => err ? reject(err) : resolve()));
    const url = `http://127.0.0.1:${appPort}`;
    async function start(now: string) {
      child = spawn(binary, [], { cwd: root, env: { ...process.env,
        WEATHER_LISTEN: `127.0.0.1:${appPort}`, WEATHER_SOURCE_URL: `http://127.0.0.1:${sourcePort}/wfs`,
        WEATHER_CACHE_DIR: join(dir, 'cache'), WEATHER_STATIC_DIR: join(root, 'web/dist'), WEATHER_TEST_NOW: now
      }, stdio: 'ignore' });
      await expect.poll(async () => {
        try { const r = await fetch(`${url}/api/water-level/hamina`); return r.status; } catch { return 0; }
      }, { timeout: 10_000 }).toBe(200);
      await expect.poll(async () => {
        const r = await fetch(`${url}/api/water-level/hamina`);
        const data = await r.json();
        return data.observation.attempt.result !== 'pending' && data.forecast.attempt.result !== 'pending';
      }, { timeout: 10_000 }).toBe(true);
    }
    async function stop() {
      if (!child) return;
      const process = child;
      child = undefined;
      await new Promise<void>((resolve, reject) => {
        process.once('exit', () => resolve());
        if (!process.kill('SIGTERM')) reject(new Error('Go process did not stop'));
      });
    }

    await start('2026-11-16T18:08:00Z');
    await page.goto(url);
    await expect(page.getByText('32,1 cm', { exact: true })).toBeVisible();
    await expect(page.getByText('Tuore', { exact: true })).toHaveCount(2);
    await stop();
    mode = 'invalid';
    await start('2026-11-16T18:09:00Z');
    await page.goto(url);
    await expect(page.getByText('32,1 cm', { exact: true })).toBeVisible();
    await expect(page.getByText(/Viimeisin haku epäonnistui/)).toHaveCount(2);
    await stop();
    mode = 'offline';
    await start('2026-11-16T18:50:00Z');
    await page.goto(url);
    await expect(page.getByText('32,1 cm', { exact: true })).toBeVisible();
    await expect(page.getByText(/Vanhentunut.*Viimeisin haku epäonnistui/)).toBeVisible();
    await expect(page.getByText(/Yhteenveto ei ole saatavilla/)).toBeVisible();
    await stop();
    mode = 'ok';
    await start('2026-11-16T18:10:00Z');
    await page.goto(url);
    await expect(page.getByText('Tuore', { exact: true })).toHaveCount(2);
    await expect(page.getByText(/Viimeisin haku epäonnistui/)).toHaveCount(0);
    await stop();
  } finally {
    if (child) {
      const process = child;
      await new Promise<void>(resolve => {
        process.once('exit', () => resolve());
        if (!process.kill('SIGTERM')) resolve();
      });
    }
    if (mock) await new Promise<void>(resolve => mock!.close(() => resolve()));
    rmSync(dir, { recursive: true, force: true });
  }
});
