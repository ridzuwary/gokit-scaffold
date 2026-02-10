package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func New(port int, logger *log.Logger) *http.Server {
	mux := http.NewServeMux()
	registerRoutes(mux)

	return &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           requestLogger(mux, logger),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}

func requestLogger(next http.Handler, logger *log.Logger) http.Handler {
	if logger == nil {
		logger = log.Default()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("method=%s path=%s remote=%s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
