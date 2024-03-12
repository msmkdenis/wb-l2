package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"dev11/internal/calendarerr"
	"dev11/internal/delivery/http/dto"
	"dev11/internal/delivery/http/mapper"
	"dev11/internal/delivery/http/middleware"
	"dev11/internal/model"
	"dev11/pkg"
)

type CalendarService interface {
	Add(ctx context.Context, e *model.Event) (*model.Event, error)
	Update(ctx context.Context, e *model.Event) (*model.Event, error)
	Delete(ctx context.Context, id, creatorID uint64) error
	GetEventsForDay(ctx context.Context, creatorId uint64, date time.Time) ([]model.Event, error)
	GetEventsForMonth(ctx context.Context, creatorId uint64, date time.Time) ([]model.Event, error)
	GetEventsForYear(ctx context.Context, creatorId uint64, date time.Time) ([]model.Event, error)
}

type Calendar struct {
	service       CalendarService
	mux           *http.ServeMux
	requestLogger *middleware.RequestLogger
}

func NewCalendar(s CalendarService, mux *http.ServeMux, requestLogger *middleware.RequestLogger) *Calendar {
	c := &Calendar{
		service: s,
		mux:     mux,
	}

	c.mux.Handle("/create_event", requestLogger.Logger(http.HandlerFunc(c.CreateEvent)))
	c.mux.Handle("/update_event", requestLogger.Logger(http.HandlerFunc(c.UpdateEvent)))
	c.mux.Handle("/delete_event", requestLogger.Logger(http.HandlerFunc(c.DeleteEvent)))
	c.mux.Handle("/get_events_for_day", requestLogger.Logger(http.HandlerFunc(c.GetEventsForDay)))
	c.mux.Handle("/get_events_for_month", requestLogger.Logger(http.HandlerFunc(c.GetEventsForMonth)))
	c.mux.Handle("/get_events_for_year", requestLogger.Logger(http.HandlerFunc(c.GetEventsForYear)))

	return c
}

func (h *Calendar) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Info("method not allowed", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("method not allowed"), w, http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		slog.Info("invalid content type", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid content type"), w, http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		slog.Info("unable to read request", "caller", pkg.Caller())
		h.writeErrorResponse(err, w, http.StatusInternalServerError)
		return
	}

	var request dto.EventCreateRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		slog.Info("unable to unmarshal request", "caller", pkg.Caller())
		h.writeErrorResponse(err, w, http.StatusBadRequest)
		return
	}

	eventId, err := strconv.ParseUint(request.ID, 10, 64)
	if err != nil || eventId == 0 {
		slog.Info("invalid event id", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid event id"), w, http.StatusBadRequest)
		return
	}

	creatorId, err := strconv.ParseUint(request.CreatorID, 10, 64)
	if err != nil || creatorId == 0 {
		slog.Info("invalid creator id", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid creator id"), w, http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), eventId, creatorId)
	if err != nil {
		if errors.Is(err, calendarerr.ErrEventNotFound) || errors.Is(err, calendarerr.ErrEventOperationNotAllowed) {
			slog.Info("event not found", "caller", pkg.Caller())
			h.writeErrorResponse(err, w, http.StatusServiceUnavailable)
			return
		}
		slog.Info("unable to add event", "caller", pkg.Caller())
		h.writeErrorResponse(err, w, http.StatusInternalServerError)
		return
	}

	h.writeResultResponse("event deleted", w)
}

func (h *Calendar) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		slog.Info("method not allowed", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("method not allowed"), w, http.StatusBadRequest)
		return
	}

	dateRequest := r.URL.Query().Get("date")
	if dateRequest == "" {
		slog.Info("invalid date", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid date"), w, http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateRequest)
	if err != nil {
		slog.Info("invalid date", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid date"), w, http.StatusBadRequest)
		return
	}

	creatorIdRequest := r.URL.Query().Get("user_id")
	if creatorIdRequest == "" {
		slog.Info("invalid creator id", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid creator id"), w, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(creatorIdRequest, 10, 64)
	if err != nil || id == 0 {
		slog.Info("invalid creator id", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid creator id"), w, http.StatusBadRequest)
		return
	}

	events, err := h.service.GetEventsForMonth(r.Context(), id, date)
	if err != nil {
		slog.Info("unable to get events for day", "caller", pkg.Caller())
		h.writeErrorResponse(err, w, http.StatusServiceUnavailable)
		return
	}

	h.writeResultResponse(events, w)
}

func (h *Calendar) GetEventsForYear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		slog.Info("method not allowed", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("method not allowed"), w, http.StatusBadRequest)
		return
	}

	dateRequest := r.URL.Query().Get("date")
	if dateRequest == "" {
		slog.Info("invalid date", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid date"), w, http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateRequest)
	if err != nil {
		slog.Info("invalid date", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid date"), w, http.StatusBadRequest)
		return
	}

	creatorIdRequest := r.URL.Query().Get("user_id")
	if creatorIdRequest == "" {
		slog.Info("invalid creator id", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid creator id"), w, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(creatorIdRequest, 10, 64)
	if err != nil || id == 0 {
		slog.Info("invalid creator id", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid creator id"), w, http.StatusBadRequest)
		return
	}

	events, err := h.service.GetEventsForYear(r.Context(), id, date)
	if err != nil {
		slog.Info("unable to get events for day", "caller", pkg.Caller())
		h.writeErrorResponse(err, w, http.StatusServiceUnavailable)
		return
	}

	h.writeResultResponse(events, w)
}

func (h *Calendar) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		slog.Info("method not allowed", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("method not allowed"), w, http.StatusBadRequest)
		return
	}

	dateRequest := r.URL.Query().Get("date")
	if dateRequest == "" {
		slog.Info("invalid date", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid date"), w, http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateRequest)
	if err != nil {
		slog.Info("invalid date", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid date"), w, http.StatusBadRequest)
		return
	}

	creatorIdRequest := r.URL.Query().Get("user_id")
	if creatorIdRequest == "" {
		slog.Info("invalid creator id", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid creator id"), w, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(creatorIdRequest, 10, 64)
	if err != nil || id == 0 {
		slog.Info("invalid creator id", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid creator id"), w, http.StatusBadRequest)
		return
	}

	events, err := h.service.GetEventsForDay(r.Context(), id, date)
	if err != nil {
		slog.Info("unable to get events for day", "caller", pkg.Caller())
		h.writeErrorResponse(err, w, http.StatusServiceUnavailable)
		return
	}

	h.writeResultResponse(events, w)
}

func (h *Calendar) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Info("method not allowed", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("method not allowed"), w, http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		slog.Info("invalid content type", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid content type"), w, http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		slog.Info("unable to read request", "caller", pkg.Caller())
		h.writeErrorResponse(err, w, http.StatusBadRequest)
		return
	}

	var request dto.EventCreateRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		slog.Info("unable to unmarshal request", "caller", pkg.Caller())
		h.writeErrorResponse(err, w, http.StatusBadRequest)
		return
	}

	calendarValidator := NewCalendarValidator()
	if err := calendarValidator.Validate(request); err != nil {
		slog.Info("invalid request", "caller", pkg.Caller(), "error", err)
		h.writeErrorResponse(err, w, http.StatusBadRequest)
		return
	}

	eventToUpdate, err := mapper.MapToEvent(request)
	if err != nil {
		slog.Info("unable to map request", "caller", pkg.Caller())
		h.writeErrorResponse(err, w, http.StatusInternalServerError)
		return
	}

	event, err := h.service.Update(r.Context(), eventToUpdate)
	if err != nil {
		if errors.Is(err, calendarerr.ErrEventNotFound) || errors.Is(err, calendarerr.ErrEventOperationNotAllowed) {
			slog.Info("event not found", "caller", pkg.Caller())
			h.writeErrorResponse(err, w, http.StatusServiceUnavailable)
			return
		}
		slog.Info("unable to update event", "caller", pkg.Caller())
		h.writeErrorResponse(err, w, http.StatusInternalServerError)
		return
	}

	h.writeResultResponse(event, w)
}

func (h *Calendar) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Info("method not allowed", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("method not allowed"), w, http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		slog.Info("invalid content type", "caller", pkg.Caller())
		h.writeErrorResponse(errors.New("invalid content type"), w, http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		slog.Info("unable to read request", "caller", pkg.Caller())
		h.writeErrorResponse(err, w, http.StatusInternalServerError)
		return
	}

	var request dto.EventCreateRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		slog.Info("unable to unmarshal request", "caller", pkg.Caller())
		h.writeErrorResponse(err, w, http.StatusBadRequest)
		return
	}

	calendarValidator := NewCalendarValidator()
	if err := calendarValidator.Validate(request); err != nil {
		slog.Info("invalid request", "caller", pkg.Caller(), "error", err)
		h.writeErrorResponse(err, w, http.StatusBadRequest)
		return
	}

	eventToSave, err := mapper.MapToEvent(request)
	if err != nil {
		slog.Info("unable to map request", "caller", pkg.Caller())
		h.writeErrorResponse(err, w, http.StatusInternalServerError)
		return
	}

	event, err := h.service.Add(r.Context(), eventToSave)
	if err != nil {
		slog.Info("unable to add event", "caller", pkg.Caller())
		h.writeErrorResponse(err, w, http.StatusInternalServerError)
		return
	}

	h.writeResultResponse(event, w)
}

func (h *Calendar) writeErrorResponse(err error, w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response, marshalErr := json.Marshal(map[string]string{"error": err.Error()})
	if marshalErr != nil {
		slog.Info("unable to marshal response", "caller", pkg.Caller())
		return
	}
	_, writeErr := w.Write(response)
	if writeErr != nil {
		slog.Info("unable to write response", "caller", pkg.Caller())
		http.Error(w, "Unable to write response", http.StatusBadRequest)
		return
	}
}

func (h *Calendar) writeResultResponse(response any, w http.ResponseWriter) {
	res, err := json.Marshal(map[string]any{"result": response})
	if err != nil {
		slog.Info("unable to marshal response", "caller", pkg.Caller())
		h.writeErrorResponse(err, w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, writeErr := w.Write(res)
	if writeErr != nil {
		slog.Info("unable to write response", "caller", pkg.Caller())
		h.writeErrorResponse(err, w, http.StatusInternalServerError)
		return
	}
}
