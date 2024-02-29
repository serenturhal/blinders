package matchapi

import (
	"blinders/packages/auth"
	"blinders/packages/match"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	App  *fiber.App
	Auth auth.Manager
	Core match.Matcher
}

func (s Service) InitRoute() {
	s.App.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("service healthy")
	})

	matchRoute := s.App.Group("/match", auth.FiberAuthMiddleware(s.Auth))
	matchRoute.Get("/suggest", s.HandleGetMatch)
}
