FROM golang:1.24-bullseye AS builder

WORKDIR /app
COPY . .

RUN go build -o gochujang ./cmd/api

FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/gochujang ./gochujang
COPY migrations ./migrations

CMD ["./gochujang"]
