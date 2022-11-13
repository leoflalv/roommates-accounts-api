package services

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/leoflalv/roommates-accounts-api/connection"
	"github.com/leoflalv/roommates-accounts-api/models"
)

type PaymentLogService struct {
	Db *mongo.Database
}

func (pl *PaymentLogService) GetPaymentLogById(id string) (models.PaymentLog, error) {

	var paymentLog models.PaymentLog
	paymentLogCollection := pl.Db.Collection(connection.PAYMENT_LOGS_COLLECTION)
	objId, err := primitive.ObjectIDFromHex(id)

	if err == nil {
		err = paymentLogCollection.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&paymentLog)
	}

	return paymentLog, err
}

func (pl *PaymentLogService) GetAllPaymentLogs() ([]models.PaymentLog, error) {

	var paymentLogs []models.PaymentLog
	paymentLogCollection := pl.Db.Collection(connection.PAYMENT_LOGS_COLLECTION)

	pointer, err := paymentLogCollection.Find(context.TODO(), bson.M{})

	if err == nil {
		err = pointer.All(context.TODO(), &paymentLogs)
	}

	return paymentLogs, err
}

type Test struct {
	UserId primitive.ObjectID `json:"_userId" bson:"_userId"`
}

func (pl *PaymentLogService) GetPaymentsByUserInvolved(id string) ([]models.PaymentLog, error) {

	var paymentLogs []models.PaymentLog
	paymentLogCollection := pl.Db.Collection(connection.PAYMENT_LOGS_COLLECTION)
	objId, err := primitive.ObjectIDFromHex(id)
	usersIds := [1]primitive.ObjectID{objId}

	if err == nil {
		cursor, queryErr := paymentLogCollection.Find(context.TODO(), bson.M{"usersInvolved": bson.M{"$in": usersIds}})

		if queryErr == nil {
			queryErr = cursor.All(context.TODO(), &paymentLogs)
		} else {
			err = queryErr
		}
	}

	return paymentLogs, err
}

func (pl *PaymentLogService) GetPaymentsLogsByPayer(id string) ([]models.PaymentLog, error) {

	var paymentLogs []models.PaymentLog
	paymentLogCollection := pl.Db.Collection(connection.PAYMENT_LOGS_COLLECTION)
	objId, err := primitive.ObjectIDFromHex(id)

	if err == nil {
		cursor, queryErr := paymentLogCollection.Find(context.TODO(), bson.M{"paidBy.userId": objId})

		if queryErr == nil {
			queryErr = cursor.All(context.TODO(), &paymentLogs)
		} else {
			err = queryErr
		}
	}

	return paymentLogs, err
}

func (u *PaymentLogService) CreatePaymentLog(paymentLog *models.PaymentLog) (models.PaymentLog, error) {
	var newPaymentLog models.PaymentLog
	paymentLogCollection := u.Db.Collection(connection.PAYMENT_LOGS_COLLECTION)

	result, err := paymentLogCollection.InsertOne(context.TODO(), paymentLog)

	if err == nil {
		newPaymentLog = models.PaymentLog{
			ID:          result.InsertedID.(primitive.ObjectID),
			Description: paymentLog.Description,
			Amount:      paymentLog.Amount,
			Portions:    paymentLog.Portions,
			PaidBy:      paymentLog.PaidBy,
		}
	}

	return newPaymentLog, err
}

func (u *PaymentLogService) UpdatePaymentLog(paymentLog models.PaymentLog) (models.PaymentLog, error) {
	paymentLogCollection := u.Db.Collection(connection.PAYMENT_LOGS_COLLECTION)

	_, err := paymentLogCollection.ReplaceOne(context.TODO(), bson.M{"_id": paymentLog.ID}, paymentLog)

	return paymentLog, err
}

func (u *PaymentLogService) RemovePaymentLog(id string) (string, error) {
	paymentLogCollection := u.Db.Collection(connection.PAYMENT_LOGS_COLLECTION)
	objId, err := primitive.ObjectIDFromHex(id)

	if err == nil {
		_, err = paymentLogCollection.DeleteOne(context.TODO(), bson.M{"_id": objId})
	}

	return id, err
}
