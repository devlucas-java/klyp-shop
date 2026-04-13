# ---------- BUILD STAGE ----------
FROM golang:1.26-alpine AS builder

# CGO deps
RUN apk add --no-cache \
    gcc \
    musl-dev \
    libc-dev \
    ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o app ./cmd/server


# ---------- FINAL STAGE ----------
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]