package suggestcore

import (
	"blinders/packages/auth"
	"blinders/packages/suggest"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	App       *fiber.App
	Auth      auth.Manager
	Suggester suggest.Suggester
}

func (s Service) InitRoute() {
	chat := s.App.Group("/suggest")
	chat.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("hello from suggest service")
	})

	authorized := chat.Group("/", auth.FiberAuthMiddleware(s.Auth))
	authorized.Post("/text", s.HandleChatSuggestion())
	authorized.Post("/chat", s.HandleTextSuggestion())
}
