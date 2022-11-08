package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserInfo struct {
	ID       primitive.ObjectID `json:"userId" bson:"userId"`
	Username string             `json:"userName" bson:"userName"`
}

type PaymentLog struct {
	ID            primitive.ObjectID   `json:"_id" bson:"_id"`
	Description   string               `json:"description,omitempty" bson:"description,omitempty"`
	Amount        float64              `json:"amount" bson:"amount"`
	UsersInvolved []primitive.ObjectID `json:"usersInvolved" bson:"usersInvolved"`
	Portions      []PaymentPortion     `json:"portions" bson:"portions"`
	PaidBy        UserInfo             `json:"paidBy" bson:"paidBy"`
}

type PaymentLogService interface {
	GetPaymentLogById(id string) (PaymentLog, error)
	GetAllPaymentLog() ([]PaymentLog, error)
	GetAllPaymentLogByPayer(id string) ([]PaymentLog, error)
	GetAllPaymentLogsByUserInvolved(userId string) ([]PaymentLog, error)
	CreatePaymentLog(paymentLog *PaymentLog) (string, error)
	UpdatePaymentLog(paymentLog *PaymentLog) (string, error)
	RemovePaymentLog(id string) (string, error)
}
