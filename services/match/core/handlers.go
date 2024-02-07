package core

import (
	"log"

	"blinders/packages/auth"
	"blinders/packages/match"

	"github.com/gofiber/fiber/v2"
)

func (s *Service) HandleGetMatch(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panic("unexpected error")
	}

	matchs, err := s.Core.Suggest(user.AuthID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"matchs": matchs})
}

func (s *Service) HandleMatch(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok {
		log.Panic("unexpected error")
	}
	toID := ctx.Query("to")
	if toID == "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "must specify id of the user"})
	}

	if err := s.Core.Match(user.AuthID, toID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "cannot process match request, " + err.Error()})
	}

	ctx.Status(fiber.StatusOK)
	return nil
}

// GetUserMatchInformation receive userID and communicate with users service to get user information.
func (s *Service) GetUserMatchInformation(_ string) (match.UserMatch, error) {
	return match.UserMatch{}, nil
}
