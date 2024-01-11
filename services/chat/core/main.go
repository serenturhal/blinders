package chatcore

import (
	"blinders/packages/auth"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	App  *fiber.App
	Auth auth.Manager
}

func (s Service) InitRoute() {
	chat := s.App.Group("/chat")
	chat.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("hello from chat service")
	})

	authorized := chat.Group("/", auth.FiberAuthMiddleware(s.Auth))
	authorized.Post("/message", handlePostMessage)
}
