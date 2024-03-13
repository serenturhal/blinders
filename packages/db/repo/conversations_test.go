package repo_test

import (
	"testing"

	"blinders/packages/db"
	"blinders/packages/db/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var convDBManager = db.NewMongoManager("mongodb://localhost:27017", "blinders")

func TestInsertIndividualConversationSuccess(t *testing.T) {
	user, _ := convDBManager.Users.InsertNewRawUser(
		models.User{FirebaseUID: primitive.NewObjectID().Hex()},
	)
	friend, _ := convDBManager.Users.InsertNewRawUser(
		models.User{FirebaseUID: primitive.NewObjectID().Hex()},
	)
	conv, err := convDBManager.Conversations.InsertIndividualConversation(user.ID, friend.ID)
	assert.Nil(t, err)
	assert.Equal(t, len(conv.Members), 2)
	assert.Equal(t, conv.CreatedBy, user.ID)
	assert.Equal(t, conv.Type, models.IndividualConversation)
	assert.Equal(t, conv.Members[0].UserID, user.ID)
	assert.Equal(t, conv.Members[1].UserID, friend.ID)
}

func TestInsertIndividualConversationFailedWithDuplicatedConversation(t *testing.T) {
	user, _ := convDBManager.Users.InsertNewRawUser(
		models.User{FirebaseUID: primitive.NewObjectID().Hex()},
	)
	friend, _ := convDBManager.Users.InsertNewRawUser(
		models.User{FirebaseUID: primitive.NewObjectID().Hex()},
	)
	_, err := convDBManager.Conversations.InsertIndividualConversation(user.ID, friend.ID)
	assert.Nil(t, err)
	conv, err := convDBManager.Conversations.InsertIndividualConversation(user.ID, friend.ID)
	assert.NotNil(t, err)
	assert.Nil(t, conv)
}
