FROM node:22.16-slim

# Install Chromium dependencies
RUN apt update && \
    apt install chromium -y && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Install dependencies, but skip Puppeteer's Chromium download
ENV PUPPETEER_SKIP_DOWNLOAD=true
COPY package.json package-lock.json* ./
RUN npm install

COPY . .

# Set Puppeteer to use system Chromium
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium

CMD ["node", "axe-runner.js"]
