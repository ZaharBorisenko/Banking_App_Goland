package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if err := f(writer, request); err != nil {
			WriteJSON(writer, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// APIServer ===================== APIServer =====================
type APIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleAccount))

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(writer http.ResponseWriter, request *http.Request) error {
	return nil
}

func (s *APIServer) handleGetAccount(writer http.ResponseWriter, request *http.Request) error {
	return nil
}

func (s *APIServer) handleCreateAccount(writer http.ResponseWriter, request *http.Request) error {
	return nil
}
func (s *APIServer) handleDeleteAccount(writer http.ResponseWriter, request *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(writer http.ResponseWriter, request *http.Request) error {
	return nil
}
