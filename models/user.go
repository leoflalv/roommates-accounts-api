package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	Username  string             `json:"username" bson:"username"`
	Password  []byte             `json:"-" bson:"password"`
	ToPay     []Debt             `json:"toPay" bson:"toPay"`
	ToCollect []Debt             `json:"toCollect" bson:"toCollect"`
}

func (model User) GetHash() string {
	return model.ID.Hex()
}

type UserService interface {
	GetUserById(id string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetAllUsers() ([]User, error)
	CreateUser(user *User) (User, error)
	UpdateUser(user *User) error
	RemoveUser(id string) error
}
