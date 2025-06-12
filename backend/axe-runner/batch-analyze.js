// batch-analyze.js
// Usage: node batch-analyze.js <startUrl> <jwtToken> [maxPages] [concurrency]
// Crawls URLs and POSTs each to /api/analyze on your Go backend

const { spawn } = require('child_process');
const fetch = require('node-fetch');

async function runCrawler(startUrl, maxPages = 100, concurrency = 4) {
  return new Promise((resolveCrawl, reject) => {
    const args = ['crawler.js', startUrl, maxPages, concurrency];
    const proc = spawn('node', args, { cwd: __dirname });
    let output = '';
    proc.stdout.on('data', d => output += d);
    proc.stderr.on('data', d => process.stderr.write(d));
    proc.on('close', code => {
      if (code === 0) {
        try {
          const urls = JSON.parse(output);
          resolveCrawl(urls);
        } catch (e) {
          reject(e);
        }
      } else {
        reject(new Error('Crawler failed'));
      }
    });
  });
}

async function postAnalyze(url, jwtToken, backendUrl = 'http://localhost:8080/api/analyze') {
  try {
    const res = await fetch(backendUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': jwtToken
      },
      body: JSON.stringify({ url })
    });
    const data = await res.json();
    if (!res.ok) {
      console.error(`Analyze failed for ${url}:`, data.message || res.statusText);
    } else {
      console.log(`Analyze started for ${url}:`, data.data || data);
    }
  } catch (e) {
    console.error(`Error posting to /api/analyze for ${url}:`, e.message);
  }
}

async function main() {
  const startUrl = process.argv[2];
  const jwtToken = process.argv[3];
  const maxPages = process.argv[4] || 100;
  const concurrency = process.argv[5] || 4;
  const backendUrl = process.env.BACKEND_ANALYZE_URL || 'http://localhost:8080/api/analyze';
  if (!startUrl || !jwtToken) {
    console.error('Usage: node batch-analyze.js <startUrl> <jwtToken> [maxPages] [concurrency]');
    process.exit(1);
  }
  console.log(`Crawling: ${startUrl}`);
  const urls = await runCrawler(startUrl, maxPages, concurrency);
  console.log(`Found ${urls.length} URLs. Posting to /api/analyze...`);

  let idx = 0;
  async function analyzeBatch(batch) {
    await Promise.all(batch.map(url => postAnalyze(url, jwtToken, backendUrl)));
  }
  while (idx < urls.length) {
    const batch = urls.slice(idx, idx + Number(concurrency));
    await analyzeBatch(batch);
    idx += Number(concurrency);
  }
  console.log('Batch analyze complete.');
}

main();
