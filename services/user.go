package services

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/leoflalv/roommates-accounts-api/connection"
	"github.com/leoflalv/roommates-accounts-api/models"
)

type UserService struct {
	Db *mongo.Database
}

func (u *UserService) GetUserById(id string) (models.User, error) {

	var user models.User
	userCollection := u.Db.Collection(connection.USERS_COLLECTION)
	objId, err := primitive.ObjectIDFromHex(id)

	if err == nil {
		err = userCollection.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&user)
	}

	return user, err
}
