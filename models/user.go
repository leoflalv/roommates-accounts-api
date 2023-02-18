package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	ToPay     []Debt             `json:"toPay,omitempty" bson:"toPay,omitempty"`
	ToCollect []Debt             `json:"toCollect,omitempty" bson:"toCollect,omitempty"`
}

func (model User) GetHash() string {
	return model.ID.Hex()
}

type UserService interface {
	GetUserById(id string) (*User, error)
	GetAllUsers() ([]User, error)
	CreateUser(user *User) (User, error)
	UpdateUser(user *User) error
	RemoveUser(id string) error
}
