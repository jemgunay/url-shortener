# url-shortener
A URL shortener service.

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


1. one way hash using URL as salt - store in DB
2. reversable - 

prevent someone from traversing all URLs, in case someone posted something sensitive by accident 


Hasher interface doesn't lock you down to the pkg definitions

## Usage

```bash
go test -race ./...
```

```bash
curl -i -XPOST "http://localhost:8080/api/v1/shorten" -d '{"original_url": "http://jemgunay.co.uk"}
```