package accesslogs

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/TwiN/go-color"
)

type Middleware struct {
	enableColors        bool
	skipPathsWithDots   bool
	ignoredPaths        []string
	ignoredPathPrefixes []string
}

// New creates a new instance of the accesslogs middleware
func New() *Middleware {
	return &Middleware{}
}

// WithColors prevents the middleware from logging requests with a dot in the path (such as static assets)
func (m *Middleware) WithColors() *Middleware {
	m.enableColors = true
	return m
}

// SkipPathsWithDots prevents the middleware from logging requests with a dot in the path (such as static assets)
func (m *Middleware) SkipPathsWithDots() *Middleware {
	m.skipPathsWithDots = true
	return m
}

// IgnorePaths sets which paths should be ignored by the middleware
func (m *Middleware) IgnorePaths(paths ...string) *Middleware {
	m.ignoredPaths = paths
	return m
}

// IgnorePathPrefixes is like Ignore, but ignores all paths that start with the given prefixes
func (m *Middleware) IgnorePathPrefixes(pathPrefixes ...string) *Middleware {
	m.ignoredPathPrefixes = pathPrefixes
	return m
}

func (m Middleware) Handler(next http.Handler) http.Handler {
	return m.HandlerFunc(next.ServeHTTP)
}

func (m Middleware) HandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if m.checkIfPathIsIgnored(r.URL.String()) {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		next.ServeHTTP(w, r)
		if m.enableColors {
			log.Printf("[accesslogs] m=%s p=%s t=%s", color.InBlue(r.Method), color.InYellow(r.URL.String()), color.InGreen(time.Since(start).Milliseconds()))
		} else {
			log.Printf("[accesslogs] m=%s p=%s t=%d", r.Method, r.URL.String(), time.Since(start).Microseconds())
		}
	}
}

func (m Middleware) checkIfPathIsIgnored(path string) bool {
	if m.skipPathsWithDots && strings.Contains(path, ".") {
		return true
	}
	// XXX: Could be optimized by sorting the array on Ignore(paths) & using binary search, but likely overkill since most people will only have a few paths to ignore
	for _, ignoredPath := range m.ignoredPaths {
		if path == ignoredPath {
			return true
		}
	}
	for _, ignoredPathPrefix := range m.ignoredPathPrefixes {
		if strings.HasPrefix(path, ignoredPathPrefix) {
			return true
		}
	}
	return false
}
