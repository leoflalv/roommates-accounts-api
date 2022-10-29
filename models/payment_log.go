package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PaymentLog struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Portions    []PaymentPortion   `bson:"inline"`
}

type PaymentLogService interface {
	GetPaymentLogById(id string) (PaymentLog, error)
	GetAllPaymentLog() ([]PaymentLog, error)
	CreatePaymentLog(paymentLog *PaymentLog) (string, error)
	UpdatePaymentLog(paymentLog *PaymentLog) (string, error)
	RemovePaymentLog(id string) (string, error)
}
