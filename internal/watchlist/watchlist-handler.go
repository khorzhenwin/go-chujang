package watchlist

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Repo Storage
}

func RegisterRoutes(r chi.Router, repo Storage) {
	h := &Handler{Repo: repo}

	r.Route("/watchlist", func(r chi.Router) {
		r.Get("/", h.GetAllHandler)
		r.Post("/", h.CreateHandler)
		r.Put("/{id}", h.UpdateHandler)
		r.Delete("/{id}", h.DeleteHandler)
	})
}

// GetAllHandler handles GET /watchlist
// @Summary      Get all watchlist items
// @Description  Returns the current watchlist
// @Tags         watchlist
// @Produce      json
// @Success      200  {array}  Ticker
// @Router       /api/v1/watchlist [get]
func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	tickers, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, "Failed to retrieve tickers", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(tickers)
	if err != nil {
		return
	}
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
func (h *Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var t Ticker
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Create(&t); err != nil {
		http.Error(w, "Failed to create ticker", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	err := json.NewEncoder(w).Encode(map[string]string{"message": "created"})
	if err != nil {
		return
	}
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
func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var t Ticker
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Update(uint(id), t); err != nil {
		http.Error(w, "Failed to update ticker", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"message": "updated"})
	if err != nil {
		return
	}
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
func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Delete(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Record not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete ticker", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
