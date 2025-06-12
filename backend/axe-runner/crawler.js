// Simple domain crawler for internal links
const puppeteer = require('puppeteer');
const { URL } = require('url');

const visited = new Set();
const toVisit = [];

async function crawl(startUrl, maxPages = 100, concurrency = 4) {
  const browser = await puppeteer.launch({ headless: "new", args: ['--no-sandbox', '--disable-setuid-sandbox'] });
  const root = new URL(startUrl);
  toVisit.push(startUrl);

  async function processUrl(url) {
    if (visited.has(url)) return;
    let page;
    try {
      page = await browser.newPage();
      await page.goto(url, { waitUntil: 'domcontentloaded', timeout: 20000 });
      visited.add(url);
      const links = await page.$$eval('a[href]', as => as.map(a => a.href));
      for (const link of links) {
        try {
          const u = new URL(link, url);
          if (u.hostname === root.hostname && !visited.has(u.href) && !toVisit.includes(u.href)) {
            toVisit.push(u.href);
          }
        } catch {}
      }
    } catch (e) {
      // Ignore navigation errors
    } finally {
      if (page) await page.close();
    }
  }

  while (toVisit.length > 0 && visited.size < maxPages) {
    const batch = [];
    while (batch.length < concurrency && toVisit.length > 0 && visited.size + batch.length < maxPages) {
      const url = toVisit.shift();
      if (!visited.has(url)) batch.push(processUrl(url));
    }
    await Promise.all(batch);
  }
  await browser.close();
  return Array.from(visited);
}

if (require.main === module) {
  const startUrl = process.argv[2];
  const maxPages = parseInt(process.argv[3] || '100', 10);
  const concurrency = parseInt(process.argv[4] || '4', 10);
  if (!startUrl) {
    console.error('Usage: node crawler.js <startUrl> [maxPages] [concurrency]');
    process.exit(1);
  }
  crawl(startUrl, maxPages, concurrency).then(urls => {
    console.log(JSON.stringify(urls, null, 2));
  });
}