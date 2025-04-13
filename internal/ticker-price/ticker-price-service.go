package ticker_price

import (
	"encoding/json"
	"fmt"
	"github.com/khorzhenwin/go-chujang/internal/config"
	"github.com/khorzhenwin/go-chujang/internal/watchlist"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Service struct {
	watchlistService watchlist.Service
	config           config.VantageConfig
}

func NewService(watchlistService *watchlist.Service, config *config.VantageConfig) *Service {
	return &Service{watchlistService: *watchlistService, config: *config}
}

func (s *Service) FindBySymbol(symbol string) *TickerPrice {
	vantageApiUrl := s.config.GetGlobalQuoteUrl(symbol)
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

func fetchPrice(externalApiUrl string, symbol string) (*TickerPrice, error) {
	resp, err := http.Get(externalApiUrl)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	log.Printf("📦 Raw response for %s: %s\n", symbol, string(body))

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

	return &TickerPrice{
		Symbol:    symbol,
		Price:     price,
		Timestamp: timestamp + "T00:00:00Z", // Add time if needed
	}, nil
}

func pollPrices(tickerService *Service, symbols []string, results chan<- TickerPrice) {
	for _, symbol := range symbols {
		go func(s string) {
			vantageApiUrl := tickerService.config.GetGlobalQuoteUrl(symbol)
			resp, err := fetchPrice(vantageApiUrl, s)
			if err != nil {
				log.Printf("❌ Error fetching %s: %v", s, err)
				return
			}
			results <- *resp
		}(symbol)
	}
}

func PollAndPushToKafka(tickerService *Service, watchlistService *watchlist.Service) {
	// poll every 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	results := make(chan TickerPrice)

	log.Println("📈 Ticker-price fetcher started")

	// Start first run immediately
	tickerList, _ := getTickersFromWatchlist(watchlistService)
	go pollPrices(tickerService, tickerList, results)

	for {
		select {
		case res := <-results:
			bytes, _ := json.Marshal(res)
			log.Printf("✅ Price: %s", bytes)

			// TODO: Push to Kafka here
		case <-ticker.C:
			if IsTradingHours(time.Now()) || os.Getenv("FORCE_POLL") == "true" {
				log.Println("🔄 Polling watchlist...")
				go pollPrices(tickerService, tickerList, results)
			}
		}
	}
}
