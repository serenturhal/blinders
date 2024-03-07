package restapi

import (
	"log"
	"net/http"

	"blinders/packages/db/repo"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConversationsService struct {
	Repo *repo.ConversationsRepo
}

func NewConversationsService(repo *repo.ConversationsRepo) *ConversationsService {
	return &ConversationsService{
		Repo: repo,
	}
}

func (s ConversationsService) GetConversationByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("invalid id:", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid id",
		})
	}

	conversation, err := s.Repo.GetConversationByID(oid)
	if err != nil {
		log.Println("can not get conversation:", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"error": "can not get conversation",
		})
	}

	return ctx.Status(http.StatusOK).JSON(conversation)
}
