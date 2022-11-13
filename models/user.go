package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Debt struct {
	UserId   primitive.ObjectID `json:"userId" bson:"userId"`
	UserName string             `json:"userName" bson:"userName"`
	Amount   float64            `json:"amount" bson:"amount"`
}

type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	ToPay     []Debt             `json:"toPay" bson:"toPay"`
	ToCollect []Debt             `json:"toCollect" bson:"toCollect"`
}

type UserService interface {
	GetUserById(id string) (User, error)
	GetAllUsers() ([]User, error)
	CreateUser(user *User) (User, error)
	UpdateUser(user *User) error
	RemoveUser(id string) (string, error)
}
