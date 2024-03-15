package restapi

import (
	"net/http"

	"blinders/packages/auth"

	"github.com/gofiber/fiber/v2"
)

func ValidateUserIDParam(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")
	userAuth := ctx.Locals(auth.UserAuthKey).(*auth.UserAuth)
	if userAuth.ID != userID {
		return ctx.Status(http.StatusForbidden).JSON(&fiber.Map{
			"error": "insufficient permissions",
		})
	}

	return ctx.Next()
}
