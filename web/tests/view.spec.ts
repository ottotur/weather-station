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
    await page.setViewportSize({ width: 1024, height: 600 });
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
  await expect(page.getByText(/7,9 cm/)).toHaveCount(0);
  offline = false;
  await page.clock.fastForward('01:01');
  await expect(page.getByText(/7,9 cm/)).toBeVisible();
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
