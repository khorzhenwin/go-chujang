FROM golang:1.24-bullseye as builder

WORKDIR /app
COPY . .

RUN go build -o gochujang ./cmd/api

FROM debian:bullseye-slim

WORKDIR /root/
COPY --from=builder /app/gochujang .

CMD ["./gochujang"]
