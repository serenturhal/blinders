package exploreapi

import (
	"encoding/json"
	"fmt"

	"blinders/packages/auth"
	"blinders/packages/db/models"

	"github.com/gofiber/fiber/v2"
)

func (s *Service) HandleGetMatch(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok || user == nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get user"})
	}

	matchs, err := s.Core.Suggest(user.AuthID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"matchs": matchs})
}

// HandleAddUserMatch will add match-information to match db
func (s *Service) HandleAddUserMatch(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if !ok || user == nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot get user"})
	}

	userMatch := new(models.MatchInfo)
	if err := json.Unmarshal(ctx.Body(), userMatch); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "service: match information required in body",
		})
	}
	// currently user.AuthID is firebaseUID
	if userMatch.FirebaseUID != user.AuthID {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Insufficient permissions. The requester and the user in the request body must match.",
		})
	}

	info, err := s.Core.AddUserMatchInformation(*userMatch)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Errorf("service: cannot add user information, %v", err).Error(),
		})
	}
	s.MatchCh <- info
	return ctx.Status(fiber.StatusOK).JSON(info)
}
