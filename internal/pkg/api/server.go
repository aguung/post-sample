package api

import (
	"fmt"
	"log"

	"post/internal/pkg/config"
	"post/internal/pkg/database"
	"post/internal/pkg/logger"
	"post/internal/router"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	cfg := config.LoadConfig()
	logger.InitLogger(cfg.App.Env)
	database.Connect(cfg)

	r := router.Init(cfg)

	return &Server{
		router: r,
	}
}

func (s *Server) Run(port int) {
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Server starting on port %d", port)
	if err := s.router.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
