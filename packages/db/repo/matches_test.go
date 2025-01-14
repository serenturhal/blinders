package repo_test

import (
	"slices"
	"testing"

	"blinders/packages/db/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInsertNewRawMatchInfo(t *testing.T) {
	rawUser := models.MatchInfo{
		UserID:    primitive.NewObjectID(),
		Name:      "name",
		Gender:    "male",
		Major:     "student",
		Native:    "vietnamese",
		Country:   "vn",
		Learnings: []string{},
		Interests: []string{},
		Age:       0,
	}
	r := manager.Matches
	usr, err := r.InsertNewRawMatchInfo(rawUser)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)

	gotWithUserID, err := r.GetMatchInfoByUserID(rawUser.UserID)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, gotWithUserID)

	deleted, err := r.DropMatchInfoByUserID(rawUser.UserID)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, deleted)
}

func TestGetMatchInfoByFirebaseUID(t *testing.T) {
	rawUser := models.MatchInfo{
		UserID:    primitive.NewObjectID(),
		Name:      "name",
		Gender:    "male",
		Major:     "student",
		Native:    "vietnamese",
		Country:   "vn",
		Learnings: []string{},
		Interests: []string{},
		Age:       0,
	}
	r := manager.Matches
	usr, err := r.InsertNewRawMatchInfo(rawUser)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)

	gotWithUserID, err := r.GetMatchInfoByUserID(rawUser.UserID)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, gotWithUserID)

	deleted, err := r.DropMatchInfoByUserID(rawUser.UserID)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, deleted)

	gotFailed, err := r.GetMatchInfoByUserID(rawUser.UserID)
	assert.NotNil(t, err)
	assert.Equal(t, models.MatchInfo{}, gotFailed)
}

func TestGetMatchInfoByUserID(t *testing.T) {
	rawUser := models.MatchInfo{
		UserID:    primitive.NewObjectID(),
		Name:      "name",
		Gender:    "male",
		Major:     "student",
		Native:    "vietnamese",
		Country:   "vn",
		Learnings: []string{},
		Interests: []string{},
		Age:       0,
	}
	r := manager.Matches
	usr, err := r.InsertNewRawMatchInfo(rawUser)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)

	gotWithUserID, err := r.GetMatchInfoByUserID(rawUser.UserID)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, gotWithUserID)

	deleted, err := r.DropMatchInfoByUserID(rawUser.UserID)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, deleted)

	gotFailed, err := r.GetMatchInfoByUserID(rawUser.UserID)
	assert.NotNil(t, err)
	assert.Equal(t, models.MatchInfo{}, gotFailed)
}

func TestGetUsersByLanguage(t *testing.T) {
	rawUser := models.MatchInfo{
		UserID:    primitive.NewObjectID(),
		Name:      "name",
		Gender:    "male",
		Major:     "student",
		Native:    "vietnamese",
		Country:   "vn",
		Learnings: []string{"english"},
		Interests: []string{},
		Age:       0,
	}
	numReturn := uint32(10)
	r := manager.Matches

	usr, err := r.DropMatchInfoByUserID(rawUser.UserID)
	if err != nil {
		assert.Equal(t, models.MatchInfo{}, usr)
	} else {
		assert.NotEmpty(t, usr)
	}

	failedGot, err := r.GetUsersByLanguage(rawUser.UserID, 10)
	assert.NotNil(t, err)
	assert.Len(t, failedGot, 0)

	usr, err = r.InsertNewRawMatchInfo(rawUser)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)

	got, err := r.GetUsersByLanguage(rawUser.UserID, numReturn)
	assert.Nil(t, err)

	assert.GreaterOrEqual(t, numReturn, uint32(len(got)))

candidateLoop:
	for _, id := range got {
		oid, err := primitive.ObjectIDFromHex(id)
		assert.Nil(t, err)
		assert.False(t, oid.IsZero())

		candidate, err := r.GetMatchInfoByUserID(oid)
		assert.Nil(t, err)
		assert.NotNil(t, candidate)
		// at here, candidate must be learning same language with curr user or natively speak the language that current
		// user is learning as well as learning language that current user is natively speak.
		for _, language := range candidate.Learnings {
			if slices.Contains[[]string, string](usr.Learnings, language) {
				// user and candidate learning same language
				continue candidateLoop
			}
		}
		assert.Contains(t, usr.Learnings, candidate.Native)
		assert.Contains(t, candidate.Learnings, usr.Native)
	}
	usr, err = r.DropMatchInfoByUserID(rawUser.UserID)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)
}

func TestDropUserWithFirebaseUID(t *testing.T) {
	rawUser := models.MatchInfo{
		UserID:    primitive.NewObjectID(),
		Name:      "name",
		Gender:    "male",
		Major:     "student",
		Native:    "vietnamese",
		Country:   "vn",
		Learnings: []string{},
		Interests: []string{},
		Age:       0,
	}
	r := manager.Matches
	usr, err := r.InsertNewRawMatchInfo(rawUser)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)

	deleted, err := r.DropMatchInfoByUserID(usr.UserID)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, deleted)

	failed, err := r.DropMatchInfoByUserID(usr.UserID)
	assert.NotNil(t, err)
	assert.Equal(t, models.MatchInfo{}, failed)
}
