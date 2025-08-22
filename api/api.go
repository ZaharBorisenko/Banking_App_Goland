package api

import (
	"encoding/json"
	"github.com/ZaharBorisenko/Banking_App_Goland/dto"
	"github.com/ZaharBorisenko/Banking_App_Goland/storage"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}
type ApiMessage struct {
	Message string `json:"message"`
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
	store      storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleGetAccount)).Methods(http.MethodGet)
	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleCreateAccount)).Methods(http.MethodPost)

	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(s.handleGetAccountByID)).Methods(http.MethodGet)
	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(s.handleDeleteAccount)).Methods(http.MethodDelete)

	log.Println("Server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleGetAccount(writer http.ResponseWriter, request *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(writer, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(writer http.ResponseWriter, request *http.Request) error {
	id, _ := uuid.Parse(mux.Vars(request)["id"])
	account, err := s.store.GetAccountById(id)

	if err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(writer http.ResponseWriter, request *http.Request) error {
	createAccount := dto.CreateAccountRequest{}
	if err := json.NewDecoder(request.Body).Decode(&createAccount); err != nil {
		return err
	}

	account := dto.NewAccount(createAccount.FirstName, createAccount.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(writer, http.StatusCreated, &account)
}
func (s *APIServer) handleDeleteAccount(writer http.ResponseWriter, request *http.Request) error {
	id, _ := uuid.Parse(mux.Vars(request)["id"])
	err := s.store.DeleteAccount(id)
	if err != nil {
		return err
	}
	return WriteJSON(writer, http.StatusOK, ApiMessage{Message: "deletion is successful!"})
}

func (s *APIServer) handleTransfer(writer http.ResponseWriter, request *http.Request) error {
	return nil
}
