# syntax=docker/dockerfile:1

ARG NODE_VERSION=24
ARG GO_VERSION=1.26
ARG PNPM_VERSION=11.7.0

FROM node:${NODE_VERSION}-bookworm AS web-builder
ARG PNPM_VERSION
WORKDIR /src/web
RUN corepack enable && corepack prepare "pnpm@${PNPM_VERSION}" --activate
COPY web/package.json web/pnpm-lock.yaml web/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile --trust-lockfile
COPY web/ ./
RUN pnpm build

FROM golang:${GO_VERSION}-bookworm AS go-builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=web-builder /src/web/dist ./web/dist
ARG NOVA_VERSION=docker
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w -X nova/internal/buildinfo.Version=${NOVA_VERSION}" -o /out/nova ./cmd/nova/

FROM debian:bookworm-slim AS runtime
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=go-builder /out/nova /app/nova
COPY --from=web-builder /src/web/dist /app/web
COPY skills /app/skills
COPY config.toml /app/config.toml

ENV NOVA_BACKEND_PORT=8080 \
    NOVA_DIR=/data/.nova \
    NOVA_SKILLS_DIR=/app/skills \
    NOVA_WEB_DIR=/app/web \
    NOVA_ALLOW_LAN_ACCESS=true

EXPOSE 8080
VOLUME ["/data"]

CMD ["/app/nova", "--port", "8080", "--no-open"]
