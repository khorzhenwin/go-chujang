# 📝 Project Summary: Penny Stock Buy Signal App

A fully containerized **real-time stock signal** system built in Go, using Kafka for decoupled event streaming, Alpha Vantage for live market data, PostgreSQL (AWS RDS) for watchlist storage, and Telegram for push notifications.

### 📡 Designed for production, tested in production. 📡
Gochujang-powered trading signals.  
**Simple. Fast. Real-time.**

---

## 📦 Architecture Overview

| Component              | Description                                  |
|-------------------------|----------------------------------------------|
| **Storage**             | In-memory sliding windows (for prices), PostgreSQL (AWS RDS) for ticker watchlist |
| **Market Data Fetcher** | Polls Alpha Vantage GLOBAL_QUOTE API every 5 minutes |
| **Kafka Cluster**       | Redpanda single-node cluster (self-hosted)   |
| **Signal Worker**       | Evaluates buy signals every 15 minutes       |
| **Notification Service**| Pushes signals to Telegram chat             |
| **Dockerized**          | Fully containerized using Docker + Compose  |

---

## ⚙️ Services

- **Ticker Price Fetcher**
    - Fetches latest ticker prices from Alpha Vantage.
    - Pushes results to Kafka topic `ticker-prices`.

- **Kafka Broker (Redpanda)**
    - Self-hosted Redpanda cluster exposed on port `9092`.

- **AWS RDS PostgreSQL**
    - Stores the list of tickers to watch.
    - Integrated with GORM ORM inside the Go application.

- **Signal Worker**
    - Maintains a sliding window of prices per symbol.
    - Detects:
        - 🚀 Uptrend (>1% increase, 4/5 points rising)
        - 🔻 Downtrend (Buy the Dip, >1% decrease, 4/5 points falling)
    - Sends Telegram notifications for confirmed signals.

- **Telegram Notifier**
    - Sends alerts instantly to your phone via bot messages.

---

## 📄 Environment Variables

```env
# RDS
DB_HOST=
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=
DB_SSL=

# Poller setting
FORCE_POLL=true

# Alpha Vantage API
ALPHA_VANTAGE_API_KEY=
ALPHA_VANTAGE_BASE_URL=https://www.alphavantage.co

# Kafka
KAFKA_USERNAME=gochujang-app
KAFKA_PASSWORD=
KAFKA_TICKER_PRICE_TOPIC=ticker-prices
KAFKA_BROKER=
KAFKA_CLIENT_ID=
KAFKA_CLIENT_SECRET=

# Telegram
TELEGRAM_BOT_TOKEN=
TELEGRAM_CHAT_ID=
```

---

## 🐳 Dockerized Deployment

**docker-compose.yml**

```yaml
version: '3.9'

services:
  swagger:
    image: golang:1.21
    working_dir: /app
    volumes:
      - .:/app
    entrypoint: [ "sh", "-c" ]
    command: |
      go install github.com/swaggo/swag/cmd/swag@latest && \
      swag init -g cmd/api/api.go --output docs

  app:
    working_dir: /app
    build:
      context: .
      dockerfile: Dockerfile
    env_file:
      - .env
    environment:
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - DB_SSL=${DB_SSL}
      - FORCE_POLL=${FORCE_POLL}
      - ALPHA_VANTAGE_API_KEY=${ALPHA_VANTAGE_API_KEY}
      - ALPHA_VANTAGE_BASE_URL=${ALPHA_VANTAGE_BASE_URL}
      - KAFKA_USERNAME=${KAFKA_USERNAME}
      - KAFKA_PASSWORD=${KAFKA_PASSWORD}
      - KAFKA_TICKER_PRICE_TOPIC=${KAFKA_TICKER_PRICE_TOPIC}
      - KAFKA_BROKER=${KAFKA_BROKER}
      - KAFKA_CLIENT_ID=${KAFKA_CLIENT_ID}
      - KAFKA_CLIENT_SECRET=${KAFKA_CLIENT_SECRET}
      - TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN}
      - TELEGRAM_CHAT_ID=${TELEGRAM_CHAT_ID}
    ports:
      - "8080:8080"
    command: [ "./gochujang" ]

```

### 🚀 To Start:

```bash
make up
```

---

## 📈 Features

- 5-minute price polling
- 15-minute signal evaluation
- 1-hour cooldown per symbol
- Uptrend and Downtrend detection
- Telegram notifications
- Watchlist stored in AWS RDS PostgreSQL
- Dockerized for easy deployment

---

## 🎯 Future Enhancements

- 📚 Historical price persistence (Redis/TSDB)
- 📈 Dashboard visualization of signals
- 🧠 Smarter quant rules (moving averages, RSI)
- 📤 Webhook integrations for broader alerting

---

![image](https://github.com/user-attachments/assets/95fe8883-ce9f-4893-9f5f-ad52f4f9f8a1)
