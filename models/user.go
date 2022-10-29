package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

type UserService interface {
	GetUserById(id string) (User, error)
	GetAllUser() ([]User, error)
	CreateUser(user *User) (string, error)
	UpdateUser(user *User) (string, error)
	RemoveUser(id string) (string, error)
}
