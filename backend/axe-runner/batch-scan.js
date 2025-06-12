// batch-scan.js
// Usage: node batch-scan.js <startUrl> [maxPages] [concurrency]
// Runs the crawler, then runs axe-runner.js for each URL found

const { spawn } = require('child_process');
const { promisify } = require('util');
const { resolve } = require('path');

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

async function runAxeRunner(url) {
  return new Promise((resolveRun, reject) => {
    const args = ['axe-runner.js'];
    const proc = spawn('node', args, { cwd: __dirname });
    proc.stdin.write(JSON.stringify({ url }) + '\n');
    proc.stdin.end();
    let output = '';
    proc.stdout.on('data', d => output += d);
    proc.stderr.on('data', d => process.stderr.write(d));
    proc.on('close', code => {
      if (code === 0) {
        resolveRun(output);
      } else {
        reject(new Error(`axe-runner failed for ${url}`));
      }
    });
  });
}

async function main() {
  const startUrl = process.argv[2];
  const maxPages = process.argv[3] || 100;
  const concurrency = process.argv[4] || 4;
  if (!startUrl) {
    console.error('Usage: node batch-scan.js <startUrl> [maxPages] [concurrency]');
    process.exit(1);
  }
  console.log(`Crawling: ${startUrl}`);
  const urls = await runCrawler(startUrl, maxPages, concurrency);
  console.log(`Found ${urls.length} URLs. Starting accessibility scans...`);

  // Run axe-runner for each URL concurrently (limit concurrency)
  let idx = 0;
  async function scanBatch(batch) {
    await Promise.all(batch.map(async url => {
      console.log(`\nScanning: ${url}`);
      try {
        const result = await runAxeRunner(url);
        console.log(result);
      } catch (e) {
        console.error(`Error scanning ${url}:`, e.message);
      }
    }));
  }

  while (idx < urls.length) {
    const batch = urls.slice(idx, idx + Number(concurrency));
    await scanBatch(batch);
    idx += Number(concurrency);
  }
  console.log('Batch scan complete.');
}

main();
