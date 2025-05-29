package infrastructure

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/frencius/loan-service/application"
	"github.com/frencius/loan-service/controller"
	"github.com/go-chi/chi"
)

type HTTPServer struct {
	Server *http.Server
}

func (hs *HTTPServer) Close() {
	ctxShutdown, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer func() {
		cancel()
	}()

	log.Println("shutting down http server")
	err := hs.Server.Shutdown(ctxShutdown)
	if err != nil {
		log.Println("http server shut down error", err)
	}
	log.Println("http server shut down is successful")
}

func RunHTTPServer(app *application.App) *HTTPServer {
	hs := &HTTPServer{}

	router := setupRouter(app)
	hs.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.AppHTTPPort),
		Handler: router,
	}

	go func(hs *HTTPServer) {
		err := hs.Server.ListenAndServe()
		if err != nil {
			log.Println("failed to start http server", err)
		}
	}(hs)

	log.Printf("http server started at port %v", hs.Server.Addr)

	return hs
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// enable cors
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Session-Key")

		if r.Method != http.MethodOptions {
			next.ServeHTTP(w, r)
		}
	})
}

func setupRouter(app *application.App) *chi.Mux {
	router := chi.NewRouter()

	healthCheckController := controller.NewHealthCheckController(app)

	// middleware
	router.Use(CORS)
	router.Get("/health-checks", healthCheckController.Ping)
	router.Route("/v1", func(r chi.Router) {

	})

	return router
}
