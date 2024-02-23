package repo_test

import (
	"testing"

	"blinders/packages/db"
	"blinders/packages/db/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var manager = db.NewMongoManager("mongodb://localhost:27017", "blinders")

func TestInsertUser(t *testing.T) {
	user := models.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	newUser, err := manager.Users.InsertNewUser(user)
	assert.Nil(t, err)
	assert.NotEqual(t, newUser.ID, primitive.ObjectID{})
	assert.Equal(t, user.ID, primitive.ObjectID{})
}

func TestInsertUserFailedWithDuplicatedFirebaseUID(t *testing.T) {
	user := models.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	_, _ = manager.Users.InsertNewUser(user)
	_, err := manager.Users.InsertNewUser(user)
	assert.NotNil(t, err)
}

func TestGetUserByFirebaseUID(t *testing.T) {
	user := models.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	user, _ = manager.Users.InsertNewUser(user)
	queriedUser, err := manager.Users.GetUserByFirebaseUID(user.FirebaseUID)
	assert.Nil(t, err)
	assert.Equal(t, user, queriedUser)
}

func TestGetUserByID(t *testing.T) {
	user := models.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	user, _ = manager.Users.InsertNewUser(user)
	queriedUser, err := manager.Users.GetUserByID(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, user, queriedUser)
}

func TestGetUserByIDNotFound(t *testing.T) {
	_, err := manager.Users.GetUserByID(primitive.NewObjectID())
	assert.NotNil(t, err)
}

func TestGetUserByFirebaseUIDNotFound(t *testing.T) {
	_, err := manager.Users.GetUserByFirebaseUID(primitive.NewObjectID().String())
	assert.NotNil(t, err)
}
