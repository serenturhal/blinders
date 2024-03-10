package auth

import (
	"log"
	"strings"

	"blinders/packages/db/repo"

	"github.com/gofiber/fiber/v2"
)

type key string

const UserAuthKey key = "user_auth_key"

func FiberAuthMiddleware(m Manager, userRepo *repo.UsersRepo) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := ctx.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid jwt, missing bearer token",
			})
		}

		jwt := strings.Split(auth, " ")[1]
		user, err := m.Verify(jwt)
		if err != nil {
			log.Println("failed to verify jwt", err)
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "failed to verify jwt",
			})
		}

		// currently, user.AuthID is firebaseUID
		usr, err := userRepo.GetUserByFirebaseUID(user.AuthID)
		if err != nil {
			log.Println("failed to get user", err)
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "failed to get user",
			})
		}

		user.ID = usr.ID.Hex()
		ctx.Locals(UserAuthKey, user)

		return ctx.Next()
	}
}
