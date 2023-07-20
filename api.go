package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var errInternalServer = errors.New("internal server error")
var errTooManyRequest = errors.New("too many requests")

type APIServer struct {
	listenAddr string
	store      DbStorage
	cache      CacheStorage
}

func NewAPIServer(listenAddr string, store DbStorage, cache CacheStorage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
		cache:      cache,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	v1Router := router.PathPrefix("/api/v1").Subrouter()

	v1Router.Handle("", perClientRateLimiter(makeHTTPHandleFunc(s.handleGetUrls)))
	v1Router.Handle("/shorten", perClientRateLimiter(makeHTTPHandleFunc(s.handleShorten)))
	v1Router.Handle("/{shortUrl}", perClientRateLimiter(makeHTTPHandleFunc(s.handleRedirect)))

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

	longUrl, err := s.cache.GetLongUrlFromShortUrl(id)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if longUrl == "" {
		url, err := s.store.GetUrlByShortUrl(id)
		if err != nil {
			return err
		}
		longUrl = url.LongUrl

		err = s.cache.SetShortUrlToLongUrl(url.ShortUrl, url.LongUrl)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return WriteJSON(w, http.StatusOK, longUrl)
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

	shortUrl := generateID()

	url, err := NewUrl(shortUrl, req.LongUrl)
	if err != nil {
		return err
	}

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
	return uuid.New().String()[:7]
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
