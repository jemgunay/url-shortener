package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	hashstub "github.com/jemgunay/url-shortener/hash/stub"
	"github.com/jemgunay/url-shortener/store"
)

func TestAPI_ShortenHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		reqURL     string
		reqBody    string
		hashVal    string
		hashErr    error
		respStatus int
		respBody   string
	}{
		{
			name:       "success_shorten",
			method:     http.MethodPost,
			reqURL:     "/api/v1/shorten",
			reqBody:    `{"original_url": "https://jemgunay.co.uk"}`,
			hashVal:    "123456",
			hashErr:    nil,
			respStatus: http.StatusOK,
			respBody:   `{"short_url":"localhost:8080/123456","short_hash":"123456","original_url":"https://jemgunay.co.uk"}`,
		},
		{
			name:       "invalid_method",
			method:     http.MethodGet,
			reqURL:     "/api/v1/shorten",
			reqBody:    `{"original_url": "https://jemgunay.co.uk"}`,
			hashVal:    "",
			hashErr:    nil,
			respStatus: http.StatusMethodNotAllowed,
			respBody:   "",
		},
		{
			name:       "hash_error",
			method:     http.MethodPost,
			reqURL:     "/api/v1/shorten",
			reqBody:    `{"original_url": "https://jemgunay.co.uk"}`,
			hashVal:    "",
			hashErr:    errors.New("error creating hash"),
			respStatus: http.StatusInternalServerError,
			respBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set up hasher, store and feed them into the API
			hashStub := hashstub.Stub{
				Val: tt.hashVal,
				Err: tt.hashErr,
			}
			storeStub := store.New()
			handlers := New(hashStub, storeStub)

			// configure the request and response writer
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, tt.reqURL, bytes.NewBuffer([]byte(tt.reqBody)))
			r.URL.Host = "localhost:8080"

			handlers.ShortenHandler(w, r)

			if w.Code != tt.respStatus {
				t.Fatalf("unexpected status, expected %d, got %d", tt.respStatus, w.Code)
			}

			respBody := w.Body.String()
			if respBody != tt.respBody {
				t.Fatalf("unexpected body, expected %s, got %s", tt.respBody, respBody)
			}
		})
	}
}

func TestAPI_RedirectHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		reqURL       string
		storePairs   map[string]string
		respStatus   int
		respLocation string
	}{
		{
			name:         "success_shorten",
			method:       http.MethodGet,
			reqURL:       "/123456",
			storePairs:   map[string]string{"123456": "https://jemgunay.co.uk"},
			respStatus:   http.StatusMovedPermanently,
			respLocation: "https://jemgunay.co.uk",
		},
		{
			name:         "invalid_method",
			method:       http.MethodPost,
			reqURL:       "/123456",
			storePairs:   nil,
			respStatus:   http.StatusMethodNotAllowed,
			respLocation: "",
		},
		{
			name:         "hash_not_found",
			method:       http.MethodGet,
			reqURL:       "/123456",
			storePairs:   nil,
			respStatus:   http.StatusNotFound,
			respLocation: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// store and API - hasher not required for this handler
			storeStub := store.New()
			for k, v := range tt.storePairs {
				storeStub.Set(k, v)
			}
			handlers := New(nil, storeStub)

			// configure the request and response writer
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, tt.reqURL, nil)

			handlers.RedirectHandler(w, r)

			if w.Code != tt.respStatus {
				t.Fatalf("unexpected status, expected %d, got %d", tt.respStatus, w.Code)
			}

			locationHeader := w.Header().Get("Location")
			if locationHeader != tt.respLocation {
				t.Fatalf("unexpected location header, expected %s, got %s", tt.respLocation, locationHeader)
			}
		})
	}
}
