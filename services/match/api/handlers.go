package matchapi

import (
	"context"
	"encoding/json"
	"time"

	"blinders/packages/auth"
	"blinders/packages/match"

	"github.com/gofiber/fiber/v2"
)

func (s *Service) HandleGetMatch(ctx *fiber.Ctx) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx.Context(), time.Second*5)
	defer cancel()

	user, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	userID := user.AuthID

	matchs, err := s.Core.Suggest(ctxWithTimeout, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"matchs": matchs})
}

func (s *Service) HandleMatch(ctx *fiber.Ctx) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx.Context(), time.Second*5)
	defer cancel()

	user, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	fromID := user.AuthID

	toID := ctx.Query("to")
	if toID == "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "must specify id of the user"})
	}

	if fromID == "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "must specify id of the user"})
	}

	if err := s.Core.Match(ctxWithTimeout, fromID, toID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "cannot process match request, " + err.Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (s *Service) HandleAddMatchUser(ctx *fiber.Ctx) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx.Context(), time.Second*5)
	defer cancel()
	user := new(match.UserMatch)
	if err := json.Unmarshal(ctx.Body(), user); err != nil {
		return err
	}

	if err := s.Core.AddUserMatch(ctxWithTimeout, *user); err != nil {
		return err
	}
	return nil
}
