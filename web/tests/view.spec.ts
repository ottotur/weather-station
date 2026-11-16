import { expect, test } from '@playwright/test';

const data = {
  serverTime: '2026-10-25T00:00:00Z',
  observation: {
    status: 'fresh',
    data: { valueCm: 32.1, observedAt: '2026-10-24T23:55:00Z' },
    attempt: { result: 'success', errorCode: null }
  },
  forecast: {
    status: 'fresh',
    data: { analysisTime: '2026-10-24T23:00:00Z', unitStatus: 'assumed-cm' },
    points: [
      { time: '2026-10-25T00:00:00Z', valueCm: 40 },
      { time: '2026-10-25T01:00:00Z', valueCm: null },
      { time: '2026-10-25T02:00:00Z', valueCm: 20 }
    ],
    hasFutureValues: true,
    lastShownAt: '2026-10-25T02:00:00Z',
    attempt: { result: 'success', errorCode: null }
  },
  summary: {
    rise: { changeCm: 7.9, at: '2026-10-25T00:00:00Z' },
    fall: { changeCm: 12.1, at: '2026-10-25T02:00:00Z' },
    hasGaps: true
  },
  summaryReason: ''
};

for (const width of [360, 768, 1024]) {
  test(`näkymä ja kaavio mukautuvat leveyteen ${width}`, async ({ page }) => {
    await page.setViewportSize({ width, height: 768 });
    await page.route('**/api/water-level/hamina', route => route.fulfill({ json: data }));
    await page.goto('/');
    await expect(page.getByText('32,1 cm', { exact: true })).toBeVisible();
    await expect(page.getByText(/7,9 cm/)).toBeVisible();
    await expect(page.getByText(/Ennusteessa on aukkoja/)).toBeVisible();
    await expect(page.locator('svg polyline')).toHaveCount(2);
    await page.getByText('Tuntitaulukko').click();
    await expect(page.locator('td').getByText('Puuttuu')).toBeVisible();
    await page.locator('svg circle').first().focus();
    await page.keyboard.press('Enter');
    await expect(page.getByText(/40,0 cm/).last()).toBeVisible();
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
    const chart = await page.locator('svg').boundingBox();
    expect(chart).not.toBeNull();
    expect(chart!.width).toBeLessThanOrEqual(width);
    await page.setViewportSize({ width, height: Math.round(width * 1.35) });
    await expect(page.getByText(/7,9 cm/)).toBeVisible();
    await expect(page.locator('table')).toBeVisible();
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
    await page.setViewportSize({ width: Math.max(width, 768), height: Math.round(width * 0.7) });
    await expect(page.locator('table')).toBeVisible();
    await expect(page.getByText(/7,9 cm/)).toBeVisible();
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
  });
}

test('paikallisen yhteyden katkos piilottaa yhteenvedon ja palautuu', async ({ page }) => {
  let offline = false;
  await page.clock.install();
  await page.route('**/api/water-level/hamina', route =>
    offline ? route.abort() : route.fulfill({ json: data }));
  await page.goto('/');
  await expect(page.getByText(/7,9 cm/)).toBeVisible();
  offline = true;
  await page.clock.fastForward('01:01');
  await expect(page.getByRole('alert')).toContainText('katkesi');
  await expect(page.getByText('Tuore', { exact: true })).toHaveCount(0);
  await expect(page.getByText(/7,9 cm/)).toHaveCount(0);
  offline = false;
  await page.clock.fastForward('01:01');
  await expect(page.getByText(/7,9 cm/)).toBeVisible();
});

test('puuttuva aineisto ja sivulle paluu päivitetään', async ({ page }) => {
  let count = 0;
  await page.route('**/api/water-level/hamina', route => {
    count++;
    const pending = {
      ...data,
      observation: { status: 'missing', data: null, attempt: { result: 'pending', errorCode: null } },
      forecast: { ...data.forecast, status: 'missing', data: null, points: [], hasFutureValues: false, lastShownAt: null, attempt: { result: 'pending', errorCode: null } },
      summary: null
    };
    return route.fulfill({ json: count === 1 ? pending : data });
  });
  await page.goto('/');
  await expect(page.getByText('Havaintoa haetaan…')).toBeVisible();
  await expect(page.getByText('Ennustetta haetaan…')).toBeVisible();
  await page.evaluate(() => document.dispatchEvent(new Event('visibilitychange')));
  await expect(page.getByText('32,1 cm', { exact: true })).toBeVisible();
  expect(count).toBe(2);
});

test.describe('kosketus ja Suomen aika', () => {
  test.use({ hasTouch: true, timezoneId: 'America/New_York', viewport: { width: 360, height: 780 } });
  test('kosketusvalinta ja toistuva tunti erottuvat ilman ulkoisia hakuja', async ({ page }) => {
    const origins = new Set<string>();
    page.on('request', request => origins.add(new URL(request.url()).origin));
    const repeated = {
      ...data,
      forecast: { ...data.forecast, points: [
        { time: '2026-10-25T00:00:00Z', valueCm: 40 },
        { time: '2026-10-25T01:00:00Z', valueCm: 41 }
      ] }
    };
    await page.route('**/api/water-level/hamina', route => route.fulfill({ json: repeated }));
    await page.goto('/');
    await page.locator('svg circle').last().tap();
    await expect(page.getByText(/41,0 cm/).last()).toBeVisible();
    await page.getByText('Tuntitaulukko').click();
    await expect(page.getByText(/UTC\+3/).first()).toBeVisible();
    await expect(page.getByText(/UTC\+2/).first()).toBeVisible();
    await expect(page.getByRole('link', { name: 'CC BY 4.0' })).toBeVisible();
    await expect(page.getByText(/käyttäjän hyväksymä oletus/)).toBeVisible();
    await expect(page.getByRole('link', { name: 'Havaintokyselyn kuvaus' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Ennustekyselyn kuvaus' })).toBeVisible();
    await expect(page.getByText(/epävarmuusraj/)).toHaveCount(0);
    expect([...origins]).toEqual(['http://127.0.0.1:4173']);
  });
});

test('vanhentunut aineisto ja tallennuksen epävarmuus näkyvät', async ({ page }) => {
  const stale = {
    ...data, summary: null,
    observation: { ...data.observation, status: 'stale', attempt: { result: 'error', errorCode: 'cache_durability_uncertain' } }
  };
  await page.route('**/api/water-level/hamina', route => route.fulfill({ json: stale }));
  await page.goto('/');
  await expect(page.getByText(/Vanhentunut/)).toBeVisible();
  await expect(page.getByText(/Tallennuksen pysyvyys on epävarma/)).toBeVisible();
  await expect(page.getByText(/7,9 cm/)).toHaveCount(0);
});

test('ennusteen vanheneminen piilottaa yhteenvedon seuraavassa päivityksessä', async ({ page }) => {
  let stale = false;
  await page.clock.install();
  await page.route('**/api/water-level/hamina', route => route.fulfill({
    json: stale ? { ...data, forecast: { ...data.forecast, status: 'stale' }, summary: null, summaryReason: 'forecast_stale' } : data
  }));
  await page.goto('/');
  await expect(page.getByText(/7,9 cm/)).toBeVisible();
  stale = true;
  await page.clock.fastForward('01:01');
  await expect(page.getByText(/7,9 cm/)).toHaveCount(0);
  await expect(page.getByText('Vanhentunut', { exact: true })).toBeVisible();
});

test('vastaamatta jäävä pyyntö merkitsee yhteyskatkon ja sallii toipumisen', async ({ page }) => {
  let count = 0;
  let release: (() => void) | undefined;
  await page.clock.install();
  await page.route('**/api/water-level/hamina', async route => {
    count++;
    if (count === 2) {
      await new Promise<void>(resolve => { release = resolve; });
      await route.abort().catch(() => {});
      return;
    }
    await route.fulfill({ json: data });
  });
  await page.goto('/');
  await expect(page.getByText(/7,9 cm/)).toBeVisible();
  await page.clock.fastForward('01:00');
  await expect.poll(() => count).toBe(2);
  await page.clock.fastForward('00:31');
  await expect(page.getByRole('alert')).toBeVisible({ timeout: 1000 });
  await expect(page.getByText(/7,9 cm/)).toHaveCount(0);
  release?.();
  await page.clock.fastForward('00:30');
  await expect(page.getByText(/7,9 cm/)).toBeVisible();
});

test('vanha valinta ei jää näkyviin ennusteen vaihtuessa', async ({ page }) => {
  let changed = false;
  await page.route('**/api/water-level/hamina', route => route.fulfill({
    json: changed ? { ...data, forecast: { ...data.forecast, points: [{ time: '2026-10-25T02:00:00Z', valueCm: 22 }] } } : data
  }));
  await page.goto('/');
  await page.locator('svg circle').first().click();
  await expect(page.getByText(/40,0 cm/).last()).toBeVisible();
  changed = true;
  await page.evaluate(() => document.dispatchEvent(new Event('visibilitychange')));
  await expect(page.getByText(/40,0 cm/)).toHaveCount(0);
});

test('päättyneen ennusteen yhteenveto kertoo todellisen syyn', async ({ page }) => {
  await page.route('**/api/water-level/hamina', route => route.fulfill({
    json: { ...data, summary: null, summaryReason: 'no_future_values', forecast: { ...data.forecast, points: [], hasFutureValues: false, lastShownAt: null } }
  }));
  await page.goto('/');
  await expect(page.getByText(/ei ole tulevia ennustearvoja/i)).toBeVisible();
});
