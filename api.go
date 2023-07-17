package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	v1Router := router.PathPrefix("/api/v1").Subrouter()

	v1Router.HandleFunc("/", makeHTTPHandleFunc(s.handleGetUrls))
	v1Router.HandleFunc("/{shortUrl}", makeHTTPHandleFunc(s.handleRedirect)).Methods(http.MethodGet)
	v1Router.HandleFunc("/shorten", makeHTTPHandleFunc(s.handleShorten)).Methods(http.MethodPost)

	log.Println("JSON API server running on port: ", s.listenAddr)

	err := http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		log.Fatal("Error connecting to server", err)
	}
}

func (s *APIServer) handleRedirect(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return fmt.Errorf("method not allowed %s", r.Method)
	}
	id := mux.Vars(r)["shortUrl"]

	// check if it's in cache lookup in db and add

	url, err := s.store.GetUrlByShortUrl(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, url)

}

func (s *APIServer) handleGetUrls(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	skip, _ := strconv.Atoi(r.URL.Query().Get("skip"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if skip < 0 {
		skip = 0
	}

	if limit < 1 || limit > 10 {
		limit = 10
	}

	urls, err := s.store.GetUrls(skip, limit)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, urls)

}

func (s *APIServer) handleShorten(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	req := new(ShortenRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	// // check if longUrl is in db n cache
	// url, err := s.store.GetUrlByLongUrl(req.LongUrl)
	// if err != nil {
	// 	return err
	// }

	// if url != nil {
	// 	return fmt.Errorf("already exist")
	// }

	id := generateID()
	shortUrl := base2Converter(id)

	url, err := NewUrl(shortUrl, req.LongUrl)
	if err != nil {
		return err
	}

	// add to set

	if err := s.store.CreateUrl(url); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusCreated, url)

}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func generateID() string {
	return time.Now().Format("20060102150405")
}

func base2Converter(id string) string {
	// do something
	return id
}

type APIError struct {
	Error string `json:"error"`
}
type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle the error
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}
