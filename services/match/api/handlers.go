package matchapi

import (
	"context"
	"encoding/json"
	"time"

	"blinders/packages/auth"
	"blinders/packages/db/models"

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

func (s *Service) HandleAddMatchUser(ctx *fiber.Ctx) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx.Context(), time.Second*5)
	defer cancel()
	user := new(models.MatchInfo)
	if err := json.Unmarshal(ctx.Body(), user); err != nil {
		return err
	}

	if err := s.Core.AddUserMatch(ctxWithTimeout, *user); err != nil {
		return err
	}
	return nil
}
