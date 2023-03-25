package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfo struct {
	ID       primitive.ObjectID `json:"userId" bson:"userId"`
	Username string             `json:"userName,omitempty" bson:"userName,omitempty"`
}

func (model UserInfo) GetHash() string {
	return model.ID.Hex()
}

type PaymentPortion struct {
	UserId   primitive.ObjectID `json:"userId" bson:"userId"`
	UserName string             `json:"userName,omitempty" bson:"userName,omitempty"`
	Portion  float64            `json:"portion" bson:"portion"`
}

type PaymentLog struct {
	ID            primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Description   string               `json:"description,omitempty" bson:"description,omitempty"`
	Amount        float64              `json:"amount" bson:"amount"`
	UsersInvolved []primitive.ObjectID `json:"usersInvolved" bson:"usersInvolved"`
	Portions      []PaymentPortion     `json:"portions" bson:"portions"`
	PaidBy        UserInfo             `json:"paidBy" bson:"paidBy"`
	CreatedAt     time.Time            `json:"-" bson:"createdAt"`
	DeletedAt     time.Time            `json:"-" bson:"deletedAt,omitempty"`
}

func (model PaymentLog) GetHash() string {
	return model.ID.Hex()
}

type PaymentLogService interface {
	GetPaymentLogById(id string) (PaymentLog, error)
	GetAllPaymentLogs(mode string, userId string) ([]PaymentLog, error)
	GetPaymentLogsByPayer(id string) ([]PaymentLog, error)
	GetPaymentLogsByUserInvolved(userId string) ([]PaymentLog, error)
	CreatePaymentLog(paymentLog *PaymentLog) (PaymentLog, error)
	UpdatePaymentLog(paymentLog *PaymentLog) error
	RemovePaymentLog(id string) error
}
