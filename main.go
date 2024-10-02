package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type URLType struct {
	FullUrlName string `json:"fullurlname"`
}

var PORT = ":3001"

var urlStore = make(map[string]string) // In-memory store for shortened URLs

func main() {
	r := chi.NewRouter()

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("oh no! wrong route!!"))
	})

	r.Post("/short", handleShorten) // POST method for creating shortened URL
	r.Get("/shorten/{shortKey}", handleRedirect) // GET method for redirection

	log.Println("Server starting on ",PORT)
	log.Fatal(http.ListenAndServe(PORT, r))
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	var fullurl URLType
	if err := json.NewDecoder(r.Body).Decode(&fullurl); err != nil || !validLink(fullurl.FullUrlName) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	shortKey := generateShortKey()
	urlStore[shortKey] = fullurl.FullUrlName // Store the mapping
	shortenedURL := fmt.Sprintf("http://localhost%s/shorten/%s",PORT, shortKey)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"shortened_url": shortenedURL})
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := chi.URLParam(r, "shortKey")
	fullURL, exists := urlStore[shortKey] // Look up the full URL
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, fullURL, http.StatusFound) // Redirect to the full URL
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 8

	rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

func validLink(link string) bool {
	r, err := regexp.Compile("^(http|https)://")
	if err != nil {
		return false
	}
	link = strings.TrimSpace(link)
	log.Printf("Checking for valid link: %s", link)
	return r.MatchString(link)
}