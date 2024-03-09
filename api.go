package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	Store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		Store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleAccount))

	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(s.handleGetAccountById))

	log.Println("JSON API server running on port:", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}

	if r.Method == "POST" {
		return s.handleCreatecount(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeteleAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.Store.GetAccounts()
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]

	fmt.Println(id)

	// account := NewAccount("Udin", "Jamal")

	return writeJSON(w, http.StatusOK, &Account{})
}

func (s *APIServer) handleCreatecount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeteleAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// encode json to response
func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handling the error
			writeJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
