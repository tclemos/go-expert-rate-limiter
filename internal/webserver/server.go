package webserver

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/tclemos/go-expert-rate-limiter/config"
	internal_middleware "github.com/tclemos/go-expert-rate-limiter/internal/webserver/middleware"
)

func Start(config config.WebServerConfig) {
	r := chi.NewRouter()

	// middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	// r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// rate limiter
	r.Use(internal_middleware.IPTokenRateLimiter(context.Background(), config))

	// REST routes
	r.Get("/", handler)

	// start REST server
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	log.Println("Web server is running on:", addr)
	http.ListenAndServe(addr, r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
