package ticker_price

import (
	"encoding/json"
	"fmt"
	"github.com/khorzhenwin/go-chujang/internal/config"
	"github.com/khorzhenwin/go-chujang/internal/kafka"
	"github.com/khorzhenwin/go-chujang/internal/models"
	"github.com/khorzhenwin/go-chujang/internal/watchlist"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Service struct {
	watchlistService watchlist.Service
	vantageConfig    config.VantageConfig
	kafkaConfig      config.KafkaConfig
}

func NewService(watchlistService *watchlist.Service, vantageConfig *config.VantageConfig, kafkaConfig *config.KafkaConfig) *Service {
	return &Service{watchlistService: *watchlistService, vantageConfig: *vantageConfig, kafkaConfig: *kafkaConfig}
}

func (s *Service) FindBySymbol(symbol string) *models.TickerPrice {
	vantageApiUrl := s.vantageConfig.GetGlobalQuoteUrl(symbol)
	tickerPrice, _ := fetchPrice(vantageApiUrl, symbol)
	return tickerPrice
}

func getTickersFromWatchlist(watchlistService *watchlist.Service) ([]string, error) {
	tickers, err := watchlistService.FindAll()
	if err != nil {
		return nil, err
	}

	var symbols []string
	for _, t := range tickers {
		symbols = append(symbols, t.Symbol)
	}

	return symbols, nil
}

func fetchPrice(externalApiUrl string, symbol string) (*models.TickerPrice, error) {
	resp, err := http.Get(externalApiUrl)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("request failed: %w", err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("❌ failed to read response: %w", err)
	}

	var result map[string]map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("❌ failed to decode JSON: %w", err)
	}

	data := result["Global Quote"]
	price := data["05. price"]
	timestamp := data["07. latest trading day"]

	return &models.TickerPrice{
		Symbol:    symbol,
		Price:     price,
		Timestamp: timestamp + "T00:00:00Z", // Add time if needed
	}, nil
}

func pollPrices(tickerService *Service, symbols []string, results chan<- models.TickerPrice) {
	for _, symbol := range symbols {
		go func(s string) {
			vantageApiUrl := tickerService.vantageConfig.GetGlobalQuoteUrl(symbol)
			resp, err := fetchPrice(vantageApiUrl, s)
			log.Printf("Raw response : " + fmt.Sprintf("%+v", resp))
			if err != nil {
				log.Printf("❌ Error fetching %s: %v", s, err)
				return
			}
			results <- *resp
		}(symbol)
	}
}

func PollAndPushToKafka(tickerService *Service, watchlistService *watchlist.Service, kafkaConfig *config.KafkaConfig) {
	// poll every 5 minutes
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	results := make(chan models.TickerPrice)

	log.Println("📈 Ticker-price fetcher started")

	// Start first run immediately
	tickerList, _ := getTickersFromWatchlist(watchlistService)
	go pollPrices(tickerService, tickerList, results)

	for {
		select {
		case res := <-results:
			bytes, _ := json.Marshal(res)
			log.Printf("✅ Price: %s", bytes)

			go kafka.PushToKafkaTopic(kafkaConfig.TickerPriceTopic, res, res.Symbol)

		case <-ticker.C:
			if IsTradingHours(time.Now()) || os.Getenv("FORCE_POLL") == "true" {
				log.Println("🔄 Polling watchlist...")
				go pollPrices(tickerService, tickerList, results)
			}
		}
	}
}

// StartSignalWorker Refer to ADR-001
func StartSignalWorker(input <-chan models.TickerPrice) {
	type PriceEntry struct {
		Timestamp time.Time
		Price     float64
	}

	priceWindows := make(map[string][]PriceEntry)
	var mu sync.Mutex

	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case msg := <-input:
				parsedTime, err := time.Parse(time.RFC3339, msg.Timestamp)
				if err != nil {
					log.Printf("⚠️ Invalid timestamp: %v", msg.Timestamp)
					continue
				}
				price, err := strconv.ParseFloat(msg.Price, 64)
				if err != nil {
					log.Printf("⚠️ Invalid price: %v", msg.Price)
					continue
				}
				entry := PriceEntry{Timestamp: parsedTime, Price: price}

				mu.Lock()
				priceWindows[msg.Symbol] = append(priceWindows[msg.Symbol], entry)
				if len(priceWindows[msg.Symbol]) > 10 {
					priceWindows[msg.Symbol] = priceWindows[msg.Symbol][len(priceWindows[msg.Symbol])-10:]
				}
				mu.Unlock()

			case <-ticker.C:
				mu.Lock()
				for symbol, window := range priceWindows {
					if len(window) < 3 {
						continue
					}
					latest := window[len(window)-1]
					prev := window[len(window)-2]
					oldest := window[0]

					if latest.Price > prev.Price && prev.Price > oldest.Price {
						log.Printf("🚀 BUY SIGNAL for %s - price increasing trend detected", symbol)
					}
				}
				mu.Unlock()
			}
		}
	}()
}
