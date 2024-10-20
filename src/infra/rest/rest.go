package rest

import (
	"context"
	"encoding/json"
	"eventstore/app"
	"eventstore/core"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type restServer struct {
	address string
	service *app.EventStoreService
	router  *mux.Router
}

func NewRestServer(address string, service *app.EventStoreService) *restServer {
	server := &restServer{
		address: address,
		service: service,
		router:  mux.NewRouter(),
	}
	server.configureRoutes()
	return server
}

func (s *restServer) configureRoutes() {
	handler := &Handler{Service: s.service}
	s.router.HandleFunc("/event", handler.handleWriteEvent).Methods("PUT")
	s.router.HandleFunc("/events", handler.handleWriteEvents).Methods("PUT")
	s.router.HandleFunc("/events", handler.handleReadEvents).Methods("GET")
}

// Start launches the HTTP server
func (s *restServer) Start() error {
	log.Printf("Starting server on %s...", s.address)
	return http.ListenAndServe(s.address, s.router)
}

// Handler manages incoming requests
type Handler struct {
	Service *app.EventStoreService
}

func (h *Handler) handleWriteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var event core.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		// http.Error(w, "Invalid input", http.StatusBadRequest)

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		msg := map[string]string{
			"msg": "Invalid input",
		}
		json.NewEncoder(w).Encode(msg)
		log.Println(err)
		return
	}
	if err := h.Service.AppendEvent(ctx, event); err != nil {
		http.Error(w, "Error saving event", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	msg := map[string]string{
		"msg": "OK",
	}
	json.NewEncoder(w).Encode(msg)
}

func (h *Handler) handleWriteEvents(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var events []core.Event
	if err := json.NewDecoder(r.Body).Decode(&events); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Println(err)
		return
	}
	if err := h.Service.AppendEvents(ctx, events); err != nil {
		http.Error(w, "Error saving events", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleReadEvents(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	// Retrieve 'start' and 'limit' query parameters
	startStr := r.URL.Query().Get("start")
	limitStr := r.URL.Query().Get("limit")

	var start, limit int
	var err error

	if startStr != "" {
		start, err = strconv.Atoi(startStr)
		if err != nil {
			http.Error(w, "Invalid 'start' parameter", http.StatusBadRequest)
			log.Println(err)
			return
		}
	} else {
		start = 0 // Default start
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid 'limit' parameter", http.StatusBadRequest)
			log.Println(err)
			return
		}
	} else {
		limit = 0 // No limit
	}

	events, err := h.Service.GetEvents(ctx, start, limit)
	if err != nil {
		http.Error(w, "Error retrieving events", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Set response header and encode events as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
