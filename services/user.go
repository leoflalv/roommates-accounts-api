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

func (u *UserService) GetUserById(id string) (*models.User, error) {
	var user models.User
	userCollection := u.Db.Collection(connection.USERS_COLLECTION)
	objId, err := primitive.ObjectIDFromHex(id)

	if err == nil {
		err = userCollection.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&user)
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
	}

	return &user, err
}

func (u *UserService) GetAllUsers() ([]models.User, error) {

	var users []models.User
	userCollection := u.Db.Collection(connection.USERS_COLLECTION)

	pointer, err := userCollection.Find(context.TODO(), bson.M{})

	if err == nil {
		err = pointer.All(context.TODO(), &users)
	}

	return users, err
}

func (u *UserService) CreateUser(user *models.User) (models.User, error) {
	var newUser models.User
	userCollection := u.Db.Collection(connection.USERS_COLLECTION)

	result, err := userCollection.InsertOne(context.TODO(), user)

	if err == nil {
		newUser = models.User{
			ID:        result.InsertedID.(primitive.ObjectID),
			Name:      user.Name,
			ToPay:     user.ToPay,
			ToCollect: user.ToCollect,
		}
	}

	return newUser, err
}

func (u *UserService) UpdateUser(user *models.User) error {
	userCollection := u.Db.Collection(connection.USERS_COLLECTION)

	_, err := userCollection.ReplaceOne(context.TODO(), bson.M{"_id": user.ID}, user)

	return err
}

func (u *UserService) RemoveUser(id string) (string, error) {
	userCollection := u.Db.Collection(connection.USERS_COLLECTION)
	objId, err := primitive.ObjectIDFromHex(id)

	if err == nil {
		_, err = userCollection.DeleteOne(context.TODO(), bson.M{"_id": objId})
	}

	return id, err
}
