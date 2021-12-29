# url-shortener

[![CircleCI](https://circleci.com/gh/jemgunay/url-shortener/tree/master.svg?style=svg)](https://circleci.com/gh/jemgunay/url-shortener/tree/master)

A URL shortener service.

## Usage

### Build & Run

Run app (default port 8080):
```bash
$ go run cmd/server/server.go
$ go run cmd/server/server.go -port=8080
```

Run tests:
```bash
$ go test -race ./...
```

### URL Shorten & Redirects

Shorten a URL:
```bash
$ curl -i -XPOST "http://localhost:8080/api/v1/shorten" -d '{"original_url": "https://jemgunay.co.uk"}'

HTTP/1.1 200 OK
Date: Tue, 28 Dec 2021 21:25:48 GMT
Content-Length: 99
Content-Type: text/plain; charset=utf-8

{"short_url":"[::1]:8080/yyE7EkqwrmyQJ","short_hash":"yyE7EkqwrmyQJ","original_url":"https://jemgunay.co.uk"}
```

Entering the `short_url` in a browser will result in a redirect to the originally submitted URL.

### CLI Tool

Shorten:
```bash
$ go run cmd/cli/cli.go -addr="http://localhost:8080" -operation="shorten" -original_url="https://jemgunay.co.uk"
2021/12/29 20:42:21 {"short_url":"[::1]:8080/yyE7RYV14457E","short_hash":"yyE7RYV14457E","original_url":"https://jemgunay.co.uk"}
```
Lookup:
```bash
$ go run cmd/cli/cli.go -addr="http://localhost:8080" -operation="lookup" -hash="yyE7RYV14457E"
2021/12/29 20:44:00 yyE7RYV14457E redirects to https://jemgunay.co.uk
```

## Design Notes

The `github.com/speps/go-hashids/v2` package was chosen for generating hashes as it is a standardised and peer-reviewed algorithm/implementation; it is generally bad practise to implement cryptography yourself if you're not a cryptography specialist. However, these implementation specifics were abstracted behind a `Hasher` interface so that other hashing implementations can be plugged in.

Similarly, a `Storage` interface fronts the map-driven K/V store so that other storage (such as persistent storage, i.e. SQL/flat file) mediums can be implemented and easily swapped out. 