package auth

import (
	"fmt"

	"blinders/packages/db/repo"

	"github.com/gofiber/fiber/v2"
)

type key string

const UserAuthKey key = "user_auth_key"

func FiberAuthMiddleware(m Manager, userRepo *repo.UsersRepo) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		jwt := ctx.Get("Authorization")
		user, err := m.Verify(jwt)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		// currently, user.AuthID is firebaseUID
		usr, err := userRepo.GetUserByFirebaseUID(user.AuthID)
		if err != nil {
			fmt.Println(err)
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "cannot get user " + err.Error(),
			})
		}

		user.ID = usr.ID.Hex()
		ctx.Locals(UserAuthKey, user)

		return ctx.Next()
	}
}
