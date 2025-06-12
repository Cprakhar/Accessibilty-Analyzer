FROM node:20-slim
WORKDIR /app
ENV PUPPETEER_SKIP_DOWNLOAD=true
RUN apt-get update && apt-get install -y chromium && rm -rf /var/lib/apt/lists/*
COPY package.json ./
RUN npm install --omit=dev
COPY axe-runner.js ./
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium
ENTRYPOINT ["node", "axe-runner.js"]
