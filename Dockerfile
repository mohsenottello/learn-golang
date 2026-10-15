# 1. React
FROM node:22-alpine AS web
WORKDIR /web
COPY web/package*.json ./
RUN npm ci
COPY web ./
RUN npm run build

# 2. Go
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.su[m] ./
RUN go mod download
COPY . .
COPY --from=web /web/dist ./web/dist
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app ./cmd/server

# 3. Runtime
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /app /app
EXPOSE 8080
USER nobody
ENTRYPOINT ["/app"]
