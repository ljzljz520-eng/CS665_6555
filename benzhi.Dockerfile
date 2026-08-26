# Official Go image with the complete toolchain required by the evaluator.
FROM golang:1.22

# Install Node.js 20 while retaining the complete Go toolchain.
RUN apt-get update && apt-get install -y curl \
    && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs \
    && corepack enable \
    && corepack prepare npm@10.8.2 --activate \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ["go.mod","go.sum","/app/"]
RUN cd /app && GOWORK=off GOTOOLCHAIN=local go mod download

COPY ["script-frontend/package.json","script-frontend/package-lock.json","/app/script-frontend/"]
RUN cd /app/script-frontend && npm install --no-audit --no-fund

# Keep the complete project, including any project-owned Dockerfile and BENZHI_README.md.
COPY . .

RUN cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go build ./...
RUN cd /app/script-frontend && npm_config_offline=true npm run build

CMD ["bash"]
