package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/jemgunay/url-shortener/api"
	"github.com/jemgunay/url-shortener/hash"
	"github.com/jemgunay/url-shortener/store"
)

func main() {
	port := flag.Int("port", 8080, "the HTTP server port")
	flag.Parse()

	// create hasher and handler instances
	hasher := hash.New()
	storage := store.New()
	apiHandlers := api.New(hasher, storage)

	// hook up HTTP handlers
	http.HandleFunc("/api/v1/shorten", apiHandlers.ShortenHandler)
	http.HandleFunc("/", apiHandlers.RedirectHandler)

	// start HTTP server
	log.Printf("HTTP server starting on port %d", *port)
	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	log.Printf("HTTP server shut down: %s", err)
}
