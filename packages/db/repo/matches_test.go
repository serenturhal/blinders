package repo_test

import (
	"fmt"
	"testing"

	"blinders/packages/db"
	"blinders/packages/db/models"
	"blinders/packages/db/repo"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	type fields struct {
		Col *mongo.Collection
	}
	type args struct {
		doc models.MatchInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.MatchInfo
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name:   "success",
			fields: fields{Col: manager.Matches.Col},
			args:   args{doc: rawUser},
			want:   rawUser,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repo.MatchesRepo{
				Col: tt.fields.Col,
			}
			got, err := r.InsertNewRawMatchInfo(tt.args.doc)
			if !tt.wantErr(t, err, fmt.Sprintf("InsertNewRawMatchInfo(%v)", tt.args.doc)) {
				return
			}
			assert.Equalf(t, tt.want, got, "InsertNewRawMatchInfo(%v)", tt.args.doc)

			gotWithFirebaseUID, err := r.GetUserByFirebaseUID(tt.args.doc.FirebaseUID)
			assert.Nil(t, err)
			assert.Equal(t, tt.args.doc, gotWithFirebaseUID)

			deleted, err := r.DropUserWithFirebaseUID(tt.args.doc.FirebaseUID)
			assert.Nil(t, err)
			assert.Equal(t, tt.args.doc, deleted)
		})
	}
}
