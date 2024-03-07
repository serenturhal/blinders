package repo_test

import (
	"testing"

	"blinders/packages/db"
	"blinders/packages/db/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var manager = db.NewMongoManager("mongodb://username:password@localhost:27017/peakee", "peakee")

func TestInsertUser(t *testing.T) {
	user := models.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	newUser, err := manager.Users.InsertNewRawUser(user)
	assert.Nil(t, err)
	assert.NotEqual(t, newUser.ID, primitive.ObjectID{})
	assert.Equal(t, user.ID, primitive.ObjectID{})
}

func TestInsertUserFailedWithDuplicatedFirebaseUID(t *testing.T) {
	user := models.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	_, _ = manager.Users.InsertNewRawUser(user)
	_, err := manager.Users.InsertNewRawUser(user)
	assert.NotNil(t, err)
}

func TestGetUserByFirebaseUID(t *testing.T) {
	user := models.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	user, _ = manager.Users.InsertNewRawUser(user)
	queriedUser, err := manager.Users.GetUserByFirebaseUID(user.FirebaseUID)
	assert.Nil(t, err)
	assert.Equal(t, user, queriedUser)
}

func TestGetUserByID(t *testing.T) {
	user := models.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	user, _ = manager.Users.InsertNewRawUser(user)
	queriedUser, err := manager.Users.GetUserByPrimitiveID(user.ID)
	assert.Nil(t, err)
	assert.Equal(t, user, queriedUser)
}

func TestGetUserByIDNotFound(t *testing.T) {
	_, err := manager.Users.GetUserByPrimitiveID(primitive.NewObjectID())
	assert.NotNil(t, err)
}

func TestGetUserByFirebaseUIDNotFound(t *testing.T) {
	_, err := manager.Users.GetUserByFirebaseUID(primitive.NewObjectID().String())
	assert.NotNil(t, err)
}

func TestUsersRepo_GetUserByUserIDSuccess(t *testing.T) {
	user := models.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	user, _ = manager.Users.InsertNewRawUser(user)
	queriedUser, err := manager.Users.GetUserByUserID(user.ID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, user, queriedUser)

	deleted, err := manager.Users.DropUserByUserID(user.ID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, user, deleted)
}

func TestUsersRepo_GetUserByUserIDFailed(t *testing.T) {
	userID := primitive.NewObjectID()

	deleted, err := manager.Users.DropUserByUserID(userID.Hex())
	if err == nil {
		assert.NotEqual(t, models.User{}, deleted)
	}

	failedGet, err := manager.Users.GetUserByUserID(userID.Hex())
	assert.NotNil(t, err)
	assert.Equal(t, models.User{}, failedGet)
}

func TestUsersRepo_DropUserByUserID(t *testing.T) {
	user := models.User{
		FirebaseUID: primitive.NewObjectID().String(),
	}
	user, _ = manager.Users.InsertNewRawUser(user)

	queriedUser, err := manager.Users.GetUserByUserID(user.ID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, user, queriedUser)

	deleted, err := manager.Users.DropUserByUserID(user.ID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, user, deleted)

	failedDelete, err := manager.Users.DropUserByUserID(user.ID.Hex())
	assert.NotNil(t, err)
	assert.Equal(t, models.User{}, failedDelete)
}
