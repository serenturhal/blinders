package repo_test

import (
	"slices"
	"testing"

	"blinders/packages/db"
	"blinders/packages/db/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mongoManager = db.NewMongoManager("mongodb://username:password@localhost:27017/peakee", "peakee")

func TestMatchesRepo_InsertNewRawMatchInfo(t *testing.T) {
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
	r := mongoManager.Matches
	usr, err := r.InsertNewRawMatchInfo(rawUser)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)

	gotWithUserID, err := r.GetMatchInfoByUserID(rawUser.UserID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, rawUser, gotWithUserID)

	deleted, err := r.DropUserWithUserID(rawUser.UserID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, rawUser, deleted)
}

func TestMatchesRepo_GetMatchInfoByFirebaseUID(t *testing.T) {
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
	r := mongoManager.Matches
	usr, err := r.InsertNewRawMatchInfo(rawUser)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)

	gotWithUserID, err := r.GetMatchInfoByUserID(rawUser.UserID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, rawUser, gotWithUserID)

	deleted, err := r.DropUserWithUserID(rawUser.UserID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, rawUser, deleted)

	gotFailed, err := r.GetMatchInfoByUserID(rawUser.UserID.Hex())
	assert.NotNil(t, err)
	assert.Equal(t, models.MatchInfo{}, gotFailed)
}

func TestMatchesRepo_GetMatchInfoByUserID(t *testing.T) {
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
	r := mongoManager.Matches
	usr, err := r.InsertNewRawMatchInfo(rawUser)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)

	gotWithUserID, err := r.GetMatchInfoByUserID(rawUser.UserID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, rawUser, gotWithUserID)

	deleted, err := r.DropUserWithUserID(rawUser.UserID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, rawUser, deleted)

	gotFailed, err := r.GetMatchInfoByUserID(rawUser.UserID.Hex())
	assert.NotNil(t, err)
	assert.Equal(t, models.MatchInfo{}, gotFailed)
}

func TestMatchesRepo_GetUsersByLanguage(t *testing.T) {
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
	r := mongoManager.Matches

	usr, err := r.DropUserWithUserID(rawUser.UserID.Hex())
	if err != nil {
		assert.Equal(t, models.MatchInfo{}, usr)
	} else {
		assert.NotEmpty(t, usr)
	}

	failedGot, err := r.GetUsersByLanguage(rawUser.UserID.Hex(), 10)
	assert.NotNil(t, err)
	assert.Len(t, failedGot, 0)

	usr, err = r.InsertNewRawMatchInfo(rawUser)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)

	got, err := r.GetUsersByLanguage(rawUser.UserID.Hex(), numReturn)
	assert.Nil(t, err)

	assert.GreaterOrEqual(t, uint32(len(got)), numReturn)

candidateLoop:
	for _, id := range got {
		candidate, err := r.GetMatchInfoByUserID(id)
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
	usr, err = r.DropUserWithUserID(rawUser.UserID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)
}

func TestMatchesRepo_DropUserWithFirebaseUID(t *testing.T) {
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
	r := mongoManager.Matches
	usr, err := r.InsertNewRawMatchInfo(rawUser)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)

	deleted, err := r.DropUserWithUserID(usr.UserID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, rawUser, deleted)

	failed, err := r.DropUserWithUserID(usr.UserID.Hex())
	assert.NotNil(t, err)
	assert.Equal(t, models.MatchInfo{}, failed)
}
