# ---- frontend: runs natively on the build host ----
FROM --platform=$BUILDPLATFORM node:22-alpine AS web
WORKDIR /web
COPY web/package*.json ./
RUN npm ci
COPY web/ .
RUN npm run build

# ---- go: native too, cross-compiles to the target ----
FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=web /web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-s -w" -o /app ./cmd/server

# ---- runtime: the only stage that is actually arm64 ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
COPY --from=build /app /app
USER nobody
EXPOSE 8080
ENTRYPOINT ["/app"]
