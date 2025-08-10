package server

import (
	"L0/internal/config"
	"L0/internal/db/postgres"
	"L0/internal/handlers"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	storage    *postgres.Storage
}

func NewServer(cfg *config.Config, storage *postgres.Storage) *Server {
	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*")

	orderHandler := handlers.NewOrderHandler(storage)

	router.Use(func(c *gin.Context) {
		c.Set("storage", storage)
		c.Next()
	})

	router.GET("/order/:id", orderHandler.GetOrder)
	router.GET("/", orderHandler.IndexPage)
	router.POST("/search", orderHandler.SearchOrder)

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: router,
		},
		storage: storage,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
