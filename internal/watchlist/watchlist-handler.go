package watchlist

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RegisterRoutes(r chi.Router) {
	r.Route("/watchlist", func(r chi.Router) {
		r.Get("/", GetAllHandler)
		r.Post("/", CreateHandler)
		r.Put("/{id}", UpdateHandler)
		r.Delete("/{id}", DeleteHandler)
	})
}

// GetAllHandler handles GET /watchlist
// @Summary      Get all watchlist items
// @Description  Returns the current watchlist
// @Tags         watchlist
// @Produce      json
// @Success      200  {array}  Ticker
// @Router       /api/v1/watchlist [get]
func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	// return dummy response for now
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]string{"AAPL", "GOOG", "TSLA"})
}

// CreateHandler handles POST /watchlist
// @Summary      Create a watchlist entry
// @Description  Adds a new stock to the watchlist
// @Tags         watchlist
// @Accept       json
// @Produce      json
// @Param        ticker  body      Ticker  true  "Ticker to add"
// @Success      201     {string}  string            "created"
// @Router       /api/v1/watchlist [post]
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var t Ticker
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	// store t in DB or memory (future)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "created"})
}

// UpdateHandler handles PUT /watchlist/{id}
// @Summary      Update a watchlist entry
// @Description  Update the symbol or notes for a given watchlist item
// @Tags         watchlist
// @Accept       json
// @Produce      json
// @Param        id      path      string   true  "Ticker ID"
// @Param        ticker  body      Ticker   true  "Updated ticker"
// @Success      200     {string}  string   "updated"
// @Failure      400     {string}  string   "bad request"
// @Router       /api/v1/watchlist/{id} [put]
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// update logic here
	w.WriteHeader(http.StatusOK)
}

// DeleteHandler handles DELETE /watchlist/{id}
// @Summary      Delete a watchlist entry
// @Description  Remove a ticker from your watchlist by ID
// @Tags         watchlist
// @Produce      json
// @Param        id  path  string  true  "Ticker ID"
// @Success      204  {string}  string  "no content"
// @Failure      400  {string}  string  "bad request"
// @Router       /api/v1/watchlist/{id} [delete]
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// delete logic here
	w.WriteHeader(http.StatusNoContent)
}
