const puppeteer = require('puppeteer');
const axeCore = require('axe-core');

async function runAxe({ url, html }) {
  const browser = await puppeteer.launch({ args: ['--no-sandbox'] });
  const page = await browser.newPage();
  if (url) {
    await page.goto(url, { waitUntil: 'networkidle2' });
  } else if (html) {
    await page.setContent(html, { waitUntil: 'networkidle2' });
  } else {
    throw new Error('Must provide url or html');
  }
  await page.addScriptTag({ content: axeCore.source });
  const results = await page.evaluate(async () => {
    return await window.axe.run();
  });
  await browser.close();
  return results;
}

async function main() {
  let input = '';
  process.stdin.setEncoding('utf8');
  for await (const chunk of process.stdin) {
    input += chunk;
  }
  let params;
  try {
    params = JSON.parse(input);
  } catch (e) {
    console.error('Invalid JSON input');
    process.exit(1);
  }
  try {
    const results = await runAxe(params);
    console.log(JSON.stringify(results));
  } catch (e) {
    console.error(e.message);
    process.exit(1);
  }
}

main();
