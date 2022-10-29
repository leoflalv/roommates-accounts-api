package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID   primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}

type UserService interface {
	GetUserById(id string) (User, error)
	GetAllUsers() ([]User, error)
	CreateUser(user *User) (User, error)
	UpdateUser(user *User) (User, error)
	RemoveUser(id string) (string, error)
}
