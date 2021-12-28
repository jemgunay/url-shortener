package api

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/jemgunay/url-shortener/hash"
	"github.com/jemgunay/url-shortener/store"
)

type shortenPayload struct {
	OriginalURL string `json:"original_url"`
}

type shortenResponse struct {
	ShortURL  string `json:"short_url"`
	ShortHash string `json:"short_hash"`
	shortenPayload
}

type API struct {
	hasher  hash.Hasher
	storage store.Storage
}

func New(hasher hash.Hasher, storage store.Storage) API {
	return API{
		hasher:  hasher,
		storage: storage,
	}
}

func (a API) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	payload := shortenPayload{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash, err := a.hasher.Hash(payload.OriginalURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := a.storage.Set(hash, payload.OriginalURL); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody := shortenResponse{
		shortenPayload: payload,
		ShortHash:      hash,
	}

	if r.URL.Host != "" {
		respBody.ShortURL = r.URL.Host + "/" + hash
	} else {
		srvAddr := r.Context().Value(http.LocalAddrContextKey).(net.Addr)
		respBody.ShortURL = srvAddr.String() + "/" + hash
	}

	respBytes, err := json.Marshal(respBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respBytes)
}

//
func (a API) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	urlComponents := strings.Split(r.URL.Path, "/")
	if len(urlComponents) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hashID := urlComponents[len(urlComponents)-1]

	longURL, err := a.storage.Get(hashID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}
