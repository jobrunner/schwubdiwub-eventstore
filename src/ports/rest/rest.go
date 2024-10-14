package rest

import (
	"encoding/json"
	"eventstore/app"
	"eventstore/domain"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// RestServer struct encapsulates the HTTP server and router
type RestServer struct {
	address string
	repo    domain.EventRepository
	router  *mux.Router
}

// NewRestServer creates a new instance of RestServer
func NewRestServer(address string, repo domain.EventRepository) *RestServer {
	return &RestServer{
		address: address,
		repo:    repo,
		router:  mux.NewRouter(),
	}
}

// ConfigureRoutes sets up the routes for the HTTP server
func (s *RestServer) ConfigureRoutes() {
	service := app.NewEventService(s.repo)
	handler := &Handler{Service: service}

	// Define routes
	s.router.HandleFunc("/events", handler.handlePost).Methods("POST")
	s.router.HandleFunc("/events", handler.handleGet).Methods("GET")
}

// Start launches the HTTP server
func (s *RestServer) Start() error {
	log.Printf("Starting server on %s...", s.address)
	return http.ListenAndServe(s.address, s.router)
}

// Handler struct manages incoming requests
type Handler struct {
	Service *app.EventService
}

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {
	var event domain.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Println(err)
		return
	}
	if err := h.Service.AppendEvent(event); err != nil {
		http.Error(w, "Error saving event", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
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

	events, err := h.Service.GetEvents(start, limit)
	if err != nil {
		http.Error(w, "Error retrieving events", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Set response header and encode events as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
