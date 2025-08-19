package api

import (
	"encoding/json"
	"fmt"
	"github.com/ZaharBorisenko/Banking_App_Goland/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
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
	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(s.handleGetAccount))

	log.Println("Server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(writer http.ResponseWriter, request *http.Request) error {
	if request.Method == "GET" {
		return s.handleGetAccount(writer, request)
	}
	if request.Method == "POST" {
		return s.handleCreateAccount(writer, request)
	}
	if request.Method == "DELETE" {
		return s.handleDeleteAccount(writer, request)
	}

	return fmt.Errorf("method not allowed %s", request.Method)
}

func (s *APIServer) handleGetAccount(writer http.ResponseWriter, request *http.Request) error {
	id, _ := uuid.Parse(mux.Vars(request)["id"])
	fmt.Println(id)
	return WriteJSON(writer, http.StatusOK, &models.Account{})
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
