package rest

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func (server *Server) Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,                             // Log API request calls
		middleware.DefaultCompress,                    // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes,                    // Redirect slashes to no slash URL versions
		middleware.Recoverer,                          // Recover from panics without crashing server
		middleware.Timeout(60*time.Second),            // Timeout requests after 60 seconds
	)
	chiCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-Auth-Token", "*"},
		Debug:            false,
	})
	router.Use(chiCors.Handler)
	router.Post("/api/collections", ErrorWrapper(server.CreateCollection))
	router.Get("/api/collections", ErrorWrapper(server.GetCollections))
	router.Get("/api/collections/{collectionName}", ErrorWrapper(server.GetSchema))
	router.Get("/ping", ErrorWrapper(PingPong))
	return router
}

func PingPong(w http.ResponseWriter, r *http.Request) (int, error) {
	w.Write([]byte("pong"))
	return http.StatusOK, nil
}
