package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Debt struct {
	UserId   primitive.ObjectID `json:"userId" bson:"userId"`
	UserName string             `json:"userName" bson:"userName"`
	Amount   float64            `json:"amount" bson:"amount"`
}

func (model Debt) GetHash() string {
	return model.UserId.Hex()
}
