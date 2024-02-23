package repo

import (
	"context"
	"log"
	"time"

	"blinders/packages/db/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Users struct {
	Col *mongo.Collection
}

func NewUsers(col *mongo.Collection) *Users {
	ctx, cal := context.WithTimeout(context.Background(), time.Second*5)
	defer cal()

	_, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"firebaseUID": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Println("can not create index for firebaseUID", err)
		return nil
	}

	return &Users{
		Col: col,
	}
}

func (r *Users) InsertNewUser(user models.User) (models.User, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second*5)
	defer cal()

	user.ID = primitive.NewObjectID()
	now := primitive.NewDateTimeFromTime(time.Now())
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := r.Col.InsertOne(ctx, user)

	return user, err
}

func (r *Users) GetUserByID(id primitive.ObjectID) (models.User, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second*5)
	defer cal()

	var user models.User
	err := r.Col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)

	return user, err
}

func (r *Users) GetUserByFirebaseUID(uid string) (models.User, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second*5)
	defer cal()

	var user models.User
	err := r.Col.FindOne(ctx, bson.M{"firebaseUID": uid}).Decode(&user)

	return user, err
}
