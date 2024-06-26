package server

import (
	"errors"

	"github.com/bukharney/giga-chat/configs"
	_repo "github.com/bukharney/giga-chat/modules/repositories"
	"github.com/bukharney/giga-chat/server/ws"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	App *gin.Engine
	Cfg *configs.Configs
	DB  *sqlx.DB
}

func NewServer(db *sqlx.DB, cfg *configs.Configs) *Server {
	return &Server{
		App: gin.Default(),
		DB:  db,
		Cfg: cfg,
	}
}

func (s *Server) Run() error {
	s.App.Use(cors.New(
		cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
		},
	))

	err := s.MapHandlers()
	if err != nil {
		return errors.New("failed to map handlers")
	}

	hub := ws.NewHub()
	go hub.Run()

	chatRepo := _repo.NewChatRepo(s.DB)
	s.App.GET("/ws/:roomId", func(c *gin.Context) {
		ws.ServeWS(c, hub, chatRepo)
	})

	err = s.App.Run()
	if err != nil {
		return errors.New("failed to run server")
	}

	return nil
}
