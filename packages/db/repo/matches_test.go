package repo_test

import (
	"slices"
	"testing"

	"blinders/packages/db"
	"blinders/packages/db/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMatchesRepo_InsertNewRawMatchInfo(t *testing.T) {
	manager := db.NewMongoManager("mongodb://username:password@localhost:27017/peakee", "peakee")
	rawUser := models.MatchInfo{
		FirebaseUID: "firebaseUID",
		Name:        "name",
		Gender:      "male",
		Major:       "student",
		Native:      "vietnamese",
		Country:     "vn",
		Learnings:   []string{},
		Interests:   []string{},
		UserID:      primitive.ObjectID{},
		Age:         0,
	}
	r := manager.Matches
	usr, err := r.InsertNewRawMatchInfo(rawUser)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)

	gotWithFirebaseUID, err := r.GetMatchInfoByFirebaseUID(rawUser.FirebaseUID)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, gotWithFirebaseUID)

	deleted, err := r.DropUserWithFirebaseUID(rawUser.FirebaseUID)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, deleted)
}

func TestMatchesRepo_GetMatchInfoByFirebaseUID(t *testing.T) {
	manager := db.NewMongoManager("mongodb://username:password@localhost:27017/peakee", "peakee")
	rawUser := models.MatchInfo{
		FirebaseUID: "firebaseUID",
		Name:        "name",
		Gender:      "male",
		Major:       "student",
		Native:      "vietnamese",
		Country:     "vn",
		Learnings:   []string{},
		Interests:   []string{},
		UserID:      primitive.ObjectID{},
		Age:         0,
	}
	r := manager.Matches
	usr, err := r.InsertNewRawMatchInfo(rawUser)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)

	gotWithFirebaseUID, err := r.GetMatchInfoByFirebaseUID(rawUser.FirebaseUID)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, gotWithFirebaseUID)

	deleted, err := r.DropUserWithFirebaseUID(rawUser.FirebaseUID)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, deleted)

	gotFailed, err := r.GetMatchInfoByFirebaseUID(rawUser.FirebaseUID)
	assert.NotNil(t, err)
	assert.Equal(t, models.MatchInfo{}, gotFailed)
}

func TestMatchesRepo_GetMatchInfoByUserID(t *testing.T) {
	manager := db.NewMongoManager("mongodb://username:password@localhost:27017/peakee", "peakee")
	rawUser := models.MatchInfo{
		FirebaseUID: "firebaseUID",
		Name:        "name",
		Gender:      "male",
		Major:       "student",
		Native:      "vietnamese",
		Country:     "vn",
		Learnings:   []string{},
		Interests:   []string{},
		UserID:      primitive.ObjectID{},
		Age:         0,
	}
	r := manager.Matches
	usr, err := r.InsertNewRawMatchInfo(rawUser)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)

	gotWithFirebaseUID, err := r.GetMatchInfoByUserID(rawUser.UserID)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, gotWithFirebaseUID)

	deleted, err := r.DropUserWithFirebaseUID(rawUser.FirebaseUID)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, deleted)

	gotFailed, err := r.GetMatchInfoByUserID(rawUser.UserID)
	assert.NotNil(t, err)
	assert.Equal(t, models.MatchInfo{}, gotFailed)
}

func TestMatchesRepo_GetUsersByLanguage(t *testing.T) {
	rawUser := models.MatchInfo{
		FirebaseUID: "firebaseUID",
		Name:        "name",
		Gender:      "male",
		Major:       "student",
		Native:      "vietnamese",
		Country:     "vn",
		Learnings:   []string{"english"},
		Interests:   []string{},
		UserID:      primitive.ObjectID{},
		Age:         0,
	}
	numReturn := uint32(10)
	r := manager.Matches

	usr, err := r.InsertNewRawMatchInfo(rawUser)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)

	got, err := r.GetUsersByLanguage(rawUser.FirebaseUID, numReturn)
	assert.Nil(t, err)

	assert.GreaterOrEqual(t, uint32(len(got)), numReturn)

candidateLoop:
	for _, id := range got {
		candidate, err := r.GetMatchInfoByFirebaseUID(id)
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
	usr, err = r.DropUserWithFirebaseUID(rawUser.FirebaseUID)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)
}
func TestMatchesRepo_DropUserWithFirebaseUID(t *testing.T) {
	manager := db.NewMongoManager("mongodb://username:password@localhost:27017/peakee", "peakee")
	rawUser := models.MatchInfo{
		FirebaseUID: "testID",
		Name:        "name",
		Gender:      "male",
		Major:       "student",
		Native:      "vietnamese",
		Country:     "vn",
		Learnings:   []string{},
		Interests:   []string{},
		UserID:      primitive.NewObjectID(),
		Age:         0,
	}
	r := manager.Matches
	usr, err := r.InsertNewRawMatchInfo(rawUser)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, usr)

	deleted, err := r.DropUserWithFirebaseUID(usr.FirebaseUID)
	assert.Nil(t, err)
	assert.Equal(t, rawUser, deleted)

	failed, err := r.DropUserWithFirebaseUID(usr.FirebaseUID)
	assert.NotNil(t, err)
	assert.Equal(t, models.MatchInfo{}, failed)
}
