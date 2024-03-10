package restapi

import (
	"blinders/packages/auth"
	"blinders/packages/db"

	"github.com/gofiber/fiber/v2"
)

type Manager struct {
	App           *fiber.App
	Auth          auth.Manager
	DB            *db.MongoManager
	Users         *UsersService
	Conversations *ConversationsService
	Messages      *MessagesService
}

func NewManager(app *fiber.App, auth auth.Manager, db *db.MongoManager) *Manager {
	return &Manager{
		App:           app,
		Auth:          auth,
		DB:            db,
		Users:         NewUsersService(db.Users),
		Conversations: NewConversationsService(db.Conversations),
		Messages:      NewMessagesService(db.Messages),
	}
}

type InitOptions struct {
	prefix string
}

func (m Manager) InitRoute(options InitOptions) error {
	if options.prefix == "" {
		options.prefix = "/"
	}

	rootRoute := m.App.Group(options.prefix)
	rootRoute.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("hello from Peakee Rest API")
	})

	authorized := rootRoute.Group("/", auth.FiberAuthMiddleware(m.Auth, m.DB.Users))
	authorized.Get("/users/:id", m.Users.GetUserByID)
	authorized.Get("/conversations/:id", m.Messages.GetMessageByID)
	authorized.Get("/messages/:id", m.Messages.GetMessageByID)

	return nil
}
