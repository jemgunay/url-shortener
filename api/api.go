package api

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/jemgunay/url-shortener/hash"
	"github.com/jemgunay/url-shortener/store"
)

// shortenPayload is the payload expected by the ShortenHandler.
type shortenPayload struct {
	OriginalURL string `json:"original_url"`
}

// shortenResponse is the payload returned by the ShortenHandler. It is composed of the shortenPayload.
type shortenResponse struct {
	ShortURL  string `json:"short_url"`
	ShortHash string `json:"short_hash"`
	shortenPayload
}

// API implements the URL shortener HTTP handlers. It also stores references to a Hasher and Storage for persisting
// short URLs.
type API struct {
	hasher  hash.Hasher
	storage store.Storage
}

// New initialises a new API.
func New(hasher hash.Hasher, storage store.Storage) API {
	return API{
		hasher:  hasher,
		storage: storage,
	}
}

// ShortenHandler takes an original URL payload and stores that URL against a hash. It returns the original URL, the
// hash and the new redirect URL (which is composed of the hash).
func (a API) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	payload := shortenPayload{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("failed to JSON unmarshal request payload: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// generate hash for the given URL
	hashID, err := a.hasher.Hash(payload.OriginalURL)
	if err != nil {
		log.Printf("failed to hash original URL: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// store the new hash against the URL
	if err := a.storage.Set(hashID, payload.OriginalURL); err != nil {
		log.Printf("failed to store URL: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody := shortenResponse{
		shortenPayload: payload,
		ShortHash:      hashID,
	}

	// use the host if available, else extract the local address from the request context
	if r.URL.Host != "" {
		respBody.ShortURL = r.URL.Host + "/" + hashID
	} else {
		srvAddr := r.Context().Value(http.LocalAddrContextKey).(net.Addr)
		respBody.ShortURL = srvAddr.String() + "/" + hashID
	}

	respBytes, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("failed to JDON marshal response payload: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respBytes)
}

// RedirectHandler extracts the hash ID following the URL's final forward slash, does a store lookup for the
// corresponding original URL and performs a 301 Redirect to that URL.
func (a API) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// extract the hash ID from the end of the URL - this adds support for URLs such as "/{hashID}" and "/api/{hashID}"
	urlComponents := strings.Split(r.URL.Path, "/")
	if len(urlComponents) == 0 {
		log.Print("expected at least one URL component")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hashID := urlComponents[len(urlComponents)-1]

	// lookup original URL associated with provided hash ID
	originalURL, err := a.storage.Get(hashID)
	if err != nil {
		if err == store.ErrKeyNotFound {
			log.Printf("URL not found for hash %s: %s", hashID, err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("failed perform store URL lookup: %s", err)
		return
	}

	// perform HTTP redirect to original URL
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}
