package server

import (
	"github.com/gofiber/fiber/v2"
	"go-sca/src/server/middleware"
)

type TopicHandler interface {
	ApplyHandlers(*fiber.App)
}

type Server struct {
	app  *fiber.App
	addr string
}

func NewServer(addr string, handlers ...TopicHandler) *Server {
	app := fiber.New()
	app.Use(middleware.RequestLogger())
	app.Use(middleware.ErrLogger())

	for _, handler := range handlers {
		handler.ApplyHandlers(app)
	}

	return &Server{
		app:  app,
		addr: addr,
	}
}

func (s *Server) Listen() error {
	return s.app.Listen(s.addr)
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
