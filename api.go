package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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
	v1Router := router.PathPrefix("/api/v1").Subrouter()

	v1Router.HandleFunc("/", handleWelcomePage)
	v1Router.HandleFunc("/{id}", handleRedirect).Methods(http.MethodGet)
	v1Router.HandleFunc("/shorten", handleShorten).Methods(http.MethodPost)

	log.Println("JSON API server running on port: ", s.listenAddr)

	err := http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		log.Fatal("Error connecting to server")
	}
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	// return nil

}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	req := new(ShortenRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		// return err
	}

	// return nil
}

func handleWelcomePage(w http.ResponseWriter, r *http.Request) {
	// return nil
}
