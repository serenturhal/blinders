package auth

import "github.com/gofiber/fiber/v2"

type key string

const UserAuthKey key = "user_auth_key"

func FiberAuthMiddleware(m Manager) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		jwt := ctx.Get("Authorization")
		user, err := m.Verify(jwt)
		if err != nil {
			_ = ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		ctx.Locals(UserAuthKey, user)

		return nil
	}
}
