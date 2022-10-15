package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PaymentPortion struct {
	UserId   primitive.ObjectID `json:"_userId" bson:"_userId"`
	UserName string             `json:"userName" bson:"userName"`
	Portion  float64            `json:"portion" bson:"portion"`
}
