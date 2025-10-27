package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"18/internal/service"
)

type Handler struct {
	calendarService *service.CalendarService
}

type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

type CreateEventRequest struct {
	UserID int       `json:"user_id"`
	Date   time.Time `json:"date"`
	Title  string    `json:"title"`
}

type UpdateEventRequest struct {
	ID     string    `json:"id"`
	UserID int       `json:"user_id"`
	Date   time.Time `json:"date"`
	Title  string    `json:"title"`
}

type DeleteEventRequest struct {
	ID string `json:"id"`
}

func NewHandler(calendarService *service.CalendarService) *Handler {
	return &Handler{calendarService: calendarService}
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	event, err := h.calendarService.CreateEvent(req.UserID, req.Date, req.Title)
	if err != nil {
		status := http.StatusServiceUnavailable
		if errors.Is(err, service.ErrInvalidDate) {
			status = http.StatusBadRequest
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Response{Result: event})
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	event, err := h.calendarService.UpdateEvent(req.ID, req.UserID, req.Date, req.Title)
	if err != nil {
		status := http.StatusServiceUnavailable
		if errors.Is(err, service.ErrNotFound) || errors.Is(err, service.ErrInvalidDate) {
			status = http.StatusBadRequest
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Response{Result: event})
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req DeleteEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.calendarService.DeleteEvent(req.ID)
	if err != nil {
		status := http.StatusServiceUnavailable
		if errors.Is(err, service.ErrNotFound) {
			status = http.StatusBadRequest
		}
		respondWithError(w, status, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Response{Result: "event deleted"})
}

func (h *Handler) GetDayEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, date, err := parseQueryParams(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	events := h.calendarService.GetDayEvents(userID, date)
	respondWithJSON(w, http.StatusOK, Response{Result: events})
}

func (h *Handler) GetWeekEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, date, err := parseQueryParams(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	events := h.calendarService.GetWeekEvents(userID, date)
	respondWithJSON(w, http.StatusOK, Response{Result: events})
}

func (h *Handler) GetMonthEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, date, err := parseQueryParams(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	events := h.calendarService.GetMonthEvents(userID, date)
	respondWithJSON(w, http.StatusOK, Response{Result: events})
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, Response{Error: message})
}

func parseQueryParams(r *http.Request) (int, time.Time, error) {
	userIDStr := r.URL.Query().Get("user_id")
	dateStr := r.URL.Query().Get("date")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, time.Time{}, service.ErrInvalidDate
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, time.Time{}, service.ErrInvalidDate
	}

	return userID, date, nil
}