# go-middlewares
[![Go Reference](https://pkg.go.dev/badge/github.com/TwiN/go-middlewares.svg)](https://pkg.go.dev/github.com/TwiN/go-middlewares)

Collection of middlewares in Go

Each middleware is in a separate package, and each has its own implementation.

## Middlewares
### accesslogs
HTTP access logs middleware

Router-wide access logs:
```go
package main

import (
	"log"
	"net/http"

	"github.com/TwiN/go-middlewares/accesslogs"
)

var exampleHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
})

func main() {
	router := http.NewServeMux()
	router.Handle("/health", exampleHandler)
	router.Handle("/hello", exampleHandler)
	router.Handle("/world", exampleHandler)
	routerWithAccessLogs := accesslogs.New().WithColors().SkipPathsWithDots().IgnorePaths("/health").Handler(router)
	log.Fatal(http.ListenAndServe(":8080", routerWithAccessLogs))
}
```

You can also use it on a per-handler basis if you wish:
```go
package main

import (
	"log"
	"net/http"

	"github.com/TwiN/go-middlewares/accesslogs"
)

var exampleHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
})

func main() {
	accessLogs := accesslogs.New().WithColors()
	http.Handle("/hello", accessLogs.Handler(exampleHandler))
	http.Handle("/world", accessLogs.Handler(exampleHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

