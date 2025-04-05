# 📝 Project Summary: Penny Stock Buy Signal App

## Description
A personal Go project designed to help identify buy signals for penny stocks in my Moomoo watchlist. This application fetches real-time and historical market data, analyzes each stock based on technical indicators (like moving averages and RSI), and alerts me when potential buy conditions are met.

### Goals
- Learn and practice Go in a real-world, finance-related context.
- Build a modular and extendable backend application.
- Gain insights into stock analysis and automate signal detection.
- Explore integrations with third-party APIs and schedulers.

### Core Features
- ⚙️ Watchlist Integration: Load penny stock tickers from a file (Moomoo export or manual).
- 📈 Market Data Fetching: Use free APIs (e.g., Alpha Vantage, Finnhub) to gather OHLCV data.
- 🔍 Signal Engine: Apply technical analysis strategies like moving average crossovers or RSI thresholds to detect buy signals.
- 📣 Notification System: Notify via CLI, email, or chat (Telegram/Slack) when a buy signal is triggered.
- 🕒 Scheduled Execution: Periodically run the analysis using a scheduler or cron job.
- 📊 (Optional) Historical signal storage for tracking performance.

### ☁️ Cloud Deployment Options

| Platform         | Pros                                                | Cons                                |
|------------------|-----------------------------------------------------|-------------------------------------|
| **AWS Lambda**   | - Pay per use, no server management                <br>- Easy to schedule with CloudWatch Events | - Cold start latency <br>- Limited execution time and memory |
| **AWS ECS Fargate** | - Run Docker containers without managing servers <br>- Supports cron-like scheduled tasks | - More configuration/setup required |
| **Render**       | - Simple UI <br>- Supports cron jobs out of the box <br>- Free tier available | - Limited customization options     |
| **Fly.io**       | - Easy Docker deployment with global edge runtime <br>- Free tier for small apps | - Still maturing, smaller ecosystem |
| **GitHub Actions** | - Easy to set up <br>- Free for public repos <br>- Cron job support with workflows | - Not suitable for long-running tasks |
| **Railway**      | - Simple full-stack deployment <br>- Environment variable management | - Slightly more opinionated workflow |

![image](https://github.com/user-attachments/assets/95fe8883-ce9f-4893-9f5f-ad52f4f9f8a1)
