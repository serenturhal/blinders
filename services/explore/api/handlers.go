package exploreapi

import (
	"context"
	"time"

	"blinders/packages/auth"

	"github.com/gofiber/fiber/v2"
)

func (s *Service) HandleGetMatch(ctx *fiber.Ctx) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx.Context(), time.Second*5)
	defer cancel()

	user, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok || user == nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get user"})
	}

	matchs, err := s.Core.Suggest(ctxWithTimeout, user.AuthID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"matchs": matchs})
}
