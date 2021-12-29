package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	addr := flag.String("addr", "http://localhost:8080", "the server instance to connect to")
	operation := flag.String("operation", "", "the target operation (shorten/lookup)")
	originalURL := flag.String("original_url", "", "the original URL to shorten")
	hash := flag.String("hash", "", "the hash to lookup")
	flag.Parse()

	var err error
	switch *operation {
	case "shorten":
		err = shorten(*addr, *originalURL)
	case "lookup":
		err = lookup(*addr, *hash)
	default:
		err = fmt.Errorf("unsupported operation arg: %s", *operation)
	}

	if err != nil {
		log.Printf("operation %s failed: %s", *operation, err)
		os.Exit(1)
	}
}

// define HTTP client with timeout and redirect disabled
var httpClient = http.Client{
	Timeout: time.Second * 5,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// shorten performs a request to the shorten handler to generate a redirect hash for the given URL.
func shorten(addr, originalURL string) error {
	req, err := http.NewRequest(http.MethodPost, addr+"/api/v1/shorten", strings.NewReader(`{"original_url": "`+originalURL+`"}`))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %s", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}

	// print the successful shorten result
	log.Printf("%s", respBody)
	return nil
}

// lookup performs a request to the redirect handler and extracts the redirect URL from the response.
func lookup(addr, hash string) error {
	req, err := http.NewRequest(http.MethodGet, addr+"/"+hash, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %s", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMovedPermanently {
		if resp.StatusCode == http.StatusNotFound {
			return errors.New("no URL found for the provided hash")
		}
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// print resulting original URL
	locationHeader := resp.Header.Get("Location")
	log.Printf("%s redirects to %s", hash, locationHeader)
	return nil
}
