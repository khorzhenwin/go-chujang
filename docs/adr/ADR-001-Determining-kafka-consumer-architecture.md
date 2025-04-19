### 🧠 Goal Summary
Build a Kafka consumer that:
- ✅ Subscribes to the ticker-prices topic
- ✅ Maintains a sliding window of price data per symbol
- ✅ Periodically evaluates data for buy signal conditions

---
## 🏗️ High-Level Architecture
#### 🔁 1. Continuous Kafka Consumer
- Use franz-go to subscribe to ticker-prices
- Group messages by symbol
- Maintain a sliding window (e.g. last 10–30 minutes) per symbol in memory

### 🧮 2. Buy Signal Engine
- Every X seconds, evaluate each symbol’s sliding window:
- Check for price momentum (e.g. ↑ 3x in a row)
- Check for threshold (e.g. sudden dip, spike, etc.)
- Trigger action (log, webhook, notification) on signal

### 💾 3. Optional: State Durability
- Keep everything in memory (fastest, simplest)
- For fault tolerance: write to Redis, SQLite, or time series DB (optional)

---
### 🔄 Suggested Polling Frequency
- Kafka consumer: real-time stream (don’t delay)
- Signal evaluation: every 10s – 60s is reasonable
- Sliding window: typically last 5–30 minutes, depending on strategy

### ⏳ Sliding Window Strategies
Strategy	|   Memory Use  |	Relevance   |	Simplicity
Last N prices   |   🟢 Low   |   🟢 Good  |   🟢 Easy
Last T minutes  |   🟠 Medium    |   🟢 Great |   🟠 Moderate
EWMA (smoothed) |   🟢 Low   |   🔵 Insightful    |   🟠 Moderate

---
### 🧪 Buy Signal Logic (MVP Ideas)
- ✅ Price dropped 10% in last 5 ticks
- ✅ Price rose 3 times in a row
- ✅ Current price < X-day moving average
