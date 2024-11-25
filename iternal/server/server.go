package server

import (
	"Yakudza/pkg/config"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"sync"
)

// Router - сущность маршрутизатора, содержит приватные поля для работы исключительно внутри пакета
type Router struct {
	r   *mux.Router
	mu  sync.Mutex
	cfg *config.Config
}

// New - создает новый роутер для маршрутизации
func New(cfg *config.Config) *http.Server {
	router := &Router{
		r:   mux.NewRouter(),
		mu:  sync.Mutex{},
		cfg: cfg,
	}

	return router.loadEndpoints()
}

func (route *Router) loadEndpoints() *http.Server {
	addr := fmt.Sprintf(":%d", route.cfg.Server.Port)

	//Эндпоинты auth
	authRoute := route.r.PathPrefix("/auth").Subrouter()
	authRoute.Use(cors.Default().Handler, route.publicMiddleware)

	//Эндпоинты links
	linksRoute := route.r.PathPrefix("/links").Subrouter()
	linksRoute.Use(cors.Default().Handler, route.authMiddleware)

	//Аутентификация
	{
		authRoute.HandleFunc("/login", route.Login).Methods(http.MethodPost, http.MethodOptions)
	}

	//Линки
	{
		linksRoute.HandleFunc("", route.GetLinks).Methods(http.MethodGet, http.MethodOptions)
		linksRoute.HandleFunc("", route.CreateLink).Methods(http.MethodPost, http.MethodOptions)
		linksRoute.HandleFunc("/{id:[0-9]+}", route.GetLink).Methods(http.MethodGet, http.MethodOptions)
		linksRoute.HandleFunc("/{id:[0-9]+}", route.UpdateLink).Methods(http.MethodPut, http.MethodOptions)
		linksRoute.HandleFunc("/{id:[0-9]+}", route.DeleteLinkByID).Methods(http.MethodDelete, http.MethodOptions)
	}

	route.r.Use(cors.Default().Handler, mux.CORSMethodMiddleware(route.r))

	// CORS обработчик
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders: []string{"Content-Type", "application/json"},
	})
	handler := crs.Handler(route.r)

	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: route.cfg.Server.Timeout,
		ReadTimeout:  route.cfg.Server.Timeout,
		IdleTimeout:  route.cfg.Server.Timeout,
		Handler:      cors.AllowAll().Handler(handler),
	}

	return srv
}
