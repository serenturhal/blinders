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

type UsersRepo struct {
	Col *mongo.Collection
}

func NewUsersRepo(col *mongo.Collection) *UsersRepo {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	_, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"firebaseUID": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Println("can not create index for firebaseUID:", err)
		return nil
	}

	return &UsersRepo{
		Col: col,
	}
}

func (r *UsersRepo) InsertNewUser(u models.User) (models.User, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	_, err := r.Col.InsertOne(ctx, u)

	return u, err
}

// this function creates new ID and time and insert the document to database
func (r *UsersRepo) InsertNewRawUser(u models.User) (models.User, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	u.ID = primitive.NewObjectID()
	now := primitive.NewDateTimeFromTime(time.Now())
	u.CreatedAt = now
	u.UpdatedAt = now

	_, err := r.Col.InsertOne(ctx, u)

	return u, err
}

func (r *UsersRepo) GetUserByID(id primitive.ObjectID) (models.User, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	var user models.User
	err := r.Col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)

	return user, err
}

func (r *UsersRepo) GetUserByFirebaseUID(uid string) (models.User, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	var user models.User
	err := r.Col.FindOne(ctx, bson.M{"firebaseUID": uid}).Decode(&user)

	return user, err
}

func (r *UsersRepo) DeleteUserByID(userID primitive.ObjectID) (models.User, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	usr := models.User{}
	err := r.Col.FindOneAndDelete(ctx, bson.M{"_id": userID}).Decode(&usr)
	return usr, err
}
