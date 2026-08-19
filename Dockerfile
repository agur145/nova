# syntax=docker/dockerfile:1

ARG NODE_VERSION=24
ARG GO_VERSION=1.26.6
ARG PNPM_VERSION=11.7.0

FROM node:${NODE_VERSION}-bookworm AS web-builder
ARG PNPM_VERSION
WORKDIR /src/web
RUN corepack enable && corepack prepare "pnpm@${PNPM_VERSION}" --activate
COPY web/package.json web/pnpm-lock.yaml web/pnpm-workspace.yaml ./
COPY web/patches ./patches
RUN pnpm install --frozen-lockfile --trust-lockfile
COPY web/ ./
RUN pnpm build

FROM golang:${GO_VERSION}-bookworm AS go-builder
WORKDIR /src
COPY go.mod go.sum ./
COPY agent/go.mod agent/go.sum ./agent/
RUN go mod download
COPY . .
COPY --from=web-builder /src/web/dist ./web/dist
ARG DENOVA_VERSION=docker
ENV GOTOOLCHAIN=local
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w -X denova/internal/buildinfo.Version=${DENOVA_VERSION}" -o /out/denova ./cmd/denova/

FROM debian:bookworm-slim AS runtime
RUN apt-get update \
    && apt-get install -y --no-install-recommends bash ca-certificates curl git ripgrep tzdata \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --system --gid 10001 denova \
    && useradd --system --uid 10001 --gid denova --home-dir /data --no-create-home denova

WORKDIR /data
COPY --from=go-builder /out/denova /app/denova
COPY --from=web-builder /src/web/dist /app/web
COPY skills /app/skills
COPY config.toml /app/config.toml
RUN mkdir -p /data && chown -R denova:denova /data

ENV HOME=/data \
    DENOVA_BACKEND_PORT=8080 \
    DENOVA_DIR=/data/.denova \
    DENOVA_SKILLS_DIR=/app/skills \
    DENOVA_WEB_DIR=/app/web \
    DENOVA_ALLOW_LAN_ACCESS=true

USER denova

EXPOSE 8080
VOLUME ["/data"]

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 CMD curl --fail --silent http://127.0.0.1:8080/ >/dev/null || exit 1

CMD ["/app/denova", "--port", "8080", "--no-open"]
