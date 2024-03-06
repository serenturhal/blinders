package restapi

import (
	"log"
	"net/http"

	"blinders/packages/db/repo"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsersService struct {
	Repo *repo.UsersRepo
}

func NewUsersService(repo *repo.UsersRepo) *UsersService {
	return &UsersService{
		Repo: repo,
	}
}

func (s UsersService) GetUserByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex((id))
	if err != nil {
		log.Println("invalid id:", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid id",
		})
	}

	user, err := s.Repo.GetUserByID(oid)
	if err != nil {
		log.Println("can not get conversation:", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "can not get conversation",
		})
	}

	return ctx.Status(http.StatusOK).JSON(user)
}
