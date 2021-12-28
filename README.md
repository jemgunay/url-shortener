# url-shortener

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
$ curl -i -XPOST "http://localhost:8080/api/v1/shorten" -d '{"original_url": "https://jemgunay.co.uk"}

HTTP/1.1 200 OK
Date: Tue, 28 Dec 2021 21:25:48 GMT
Content-Length: 99
Content-Type: text/plain; charset=utf-8

{"short_url":"[::1]:8080/yyE7EkqwrmyQJ","short_hash":"yyE7EkqwrmyQJ","original_url":"https://jemgunay.co.uk"}
```

Entering the `short_url` in a browser will result in a redirect to the originally submitted URL.

### CLI Tool

TODO

## Design

URL Shortener
Write an API for a URL shortener that satisfies the following behaviour:
Accepts a URL to be shortened.
Generates a short URL for the original URL.
Accepts a short URL and redirects the caller to the original URL.

Bonus points
Comes with a CLI that can be used to call your service.

Things we'd like to see
Sound design approach that's not overly complicated or over-engineered.
Code that's easy to read and not too "clever".
Sensible tests in place

one way hash using URL as salt - store in DB

prevent someone from traversing all URLs, in case someone posted something sensitive by accident 

Hasher interface doesn't lock you down to the pkg definitions
note on 3rd party short uuid pkg - bad to implement cryptography yourself  