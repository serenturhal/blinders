package wschat

import (
	"testing"

	"blinders/packages/db"
	"blinders/packages/db/models"
	"blinders/packages/session"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	if app == nil {
		InitApp(
			session.NewManager(redis.NewClient(&redis.Options{Addr: "localhost:6379"})),
			db.NewMongoManager("mongodb://localhost:27017", "blinders"),
		)
	}
}

func TestSendMessageFailedWithWrongPayload(t *testing.T) {
	_, err := HandleSendMessage(
		primitive.NewObjectID().Hex(),
		primitive.NewObjectID().Hex(),
		UserSendMessagePayload{
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        "hello world",
			ConversationID: "wrongID",
			ResolveID:      "resolveID",
		})

	assert.NotNil(t, err)

	_, err = HandleSendMessage(
		primitive.NewObjectID().Hex(),
		primitive.NewObjectID().Hex(),
		UserSendMessagePayload{
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        "hello world",
			ConversationID: primitive.NewObjectID().Hex(),
			ResolveID:      "resolveID",
		})

	assert.NotNil(t, err)
}

func TestSendMessageFailedWithConversationNotFound(t *testing.T) {
	_, err := HandleSendMessage(
		primitive.NewObjectID().Hex(),
		primitive.NewObjectID().Hex(),
		UserSendMessagePayload{
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        "hello world",
			ConversationID: primitive.NewObjectID().Hex(),
		})

	assert.NotNil(t, err)
}

func TestSendMessageFailedWithUserIsNotMember(t *testing.T) {
	conversation, _ := app.DB.Conversations.InsertNewConversation(models.Conversation{})
	_, err := HandleSendMessage(
		primitive.NewObjectID().Hex(),
		primitive.NewObjectID().Hex(),
		UserSendMessagePayload{
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        "hello world",
			ConversationID: conversation.ID.Hex(),
		})

	assert.NotNil(t, err)
}

func TestSendMessageWithNoError(t *testing.T) {
	user, _ := app.DB.Users.InsertNewRawUser(models.User{})
	conversation, _ := app.DB.Conversations.InsertNewRawConversation(models.Conversation{
		Members: []models.Member{{UserID: user.ID}},
	})
	_, err := HandleSendMessage(
		user.ID.Hex(),
		primitive.NewObjectID().Hex(),
		UserSendMessagePayload{
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        "hello world",
			ConversationID: conversation.ID.Hex(),
		})

	assert.Nil(t, err)
}

func TestSendMessageFailedWithInvalidMessageToReply(t *testing.T) {
	_, err := HandleSendMessage(
		primitive.NewObjectID().Hex(),
		primitive.NewObjectID().Hex(),
		UserSendMessagePayload{
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        "hello world",
			ConversationID: primitive.NewObjectID().Hex(),
			ReplyTo:        primitive.NewObjectID().Hex(),
		})

	assert.NotNil(t, err)
}

func TestSendMessageWithValidMessageToReply(t *testing.T) {
	user, _ := app.DB.Users.InsertNewRawUser(models.User{})
	conversation, _ := app.DB.Conversations.InsertNewRawConversation(models.Conversation{
		Members: []models.Member{{UserID: user.ID}},
	})
	message, _ := app.DB.Messages.InsertNewRawMessage(models.Message{
		ConversationID: conversation.ID,
	})

	_, err := HandleSendMessage(
		user.ID.Hex(),
		primitive.NewObjectID().Hex(),
		UserSendMessagePayload{
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        "hello world",
			ConversationID: conversation.ID.Hex(),
			ReplyTo:        message.ID.Hex(),
		})

	assert.Nil(t, err)
}
