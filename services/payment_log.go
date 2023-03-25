package services

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/leoflalv/roommates-accounts-api/connection"
	"github.com/leoflalv/roommates-accounts-api/models"
	"github.com/leoflalv/roommates-accounts-api/utils"
)

type PaymentLogService struct {
	Db *mongo.Database
}

func (pl *PaymentLogService) GetPaymentLogById(id string) (models.PaymentLog, error) {

	var paymentLog models.PaymentLog
	paymentLogCollection := pl.Db.Collection(connection.PAYMENT_LOGS_COLLECTION)
	objId, err := primitive.ObjectIDFromHex(id)

	filter := bson.M{
		"_id": objId,
		"deletedAt": bson.M{
			"$exists": false,
		}}

	if err == nil {
		err = paymentLogCollection.FindOne(context.TODO(), filter).Decode(&paymentLog)
	}

	return paymentLog, err
}

func (pl *PaymentLogService) GetAllPaymentLogs(mode string, userId string) ([]models.PaymentLog, error) {

	var paymentLogs []models.PaymentLog
	paymentLogCollection := pl.Db.Collection(connection.PAYMENT_LOGS_COLLECTION)

	var idFilter primitive.M

	objId, err := primitive.ObjectIDFromHex(userId)

	switch mode {
	case utils.OnlyPaid:
		idFilter = bson.M{"paidBy.userId": objId}
	case utils.OnlyInvolved:
		idFilter = bson.M{"usersInvolved": bson.M{"$in": []primitive.ObjectID{objId}}}
	default:
		idFilter = bson.M{"$or": []bson.M{
			{"usersInvolved": bson.M{"$in": []primitive.ObjectID{objId}}},
			{"paidBy.userId": objId},
		}}
	}

	filter := bson.D{{Key: "$and",
		Value: []bson.M{
			idFilter,
			{"deletedAt": bson.M{
				"$exists": false,
			}},
		},
	}}

	pointer, err := paymentLogCollection.Find(context.TODO(), filter)

	if err == nil {
		err = pointer.All(context.TODO(), &paymentLogs)
	}

	return paymentLogs, err
}

func (u *PaymentLogService) CreatePaymentLog(paymentLog *models.PaymentLog) (models.PaymentLog, error) {
	var newPaymentLog models.PaymentLog
	paymentLogCollection := u.Db.Collection(connection.PAYMENT_LOGS_COLLECTION)

	var usersInvolved []primitive.ObjectID
	for _, portion := range paymentLog.Portions {
		usersInvolved = append(usersInvolved, portion.UserId)
	}

	paymentLogToCreate := models.PaymentLog{
		Description:   paymentLog.Description,
		Amount:        paymentLog.Amount,
		Portions:      paymentLog.Portions,
		PaidBy:        paymentLog.PaidBy,
		UsersInvolved: usersInvolved,
		CreatedAt:     time.Now(),
	}

	result, err := paymentLogCollection.InsertOne(context.TODO(), paymentLogToCreate)

	if err == nil {
		newPaymentLog = models.PaymentLog{
			ID:            result.InsertedID.(primitive.ObjectID),
			Description:   paymentLogToCreate.Description,
			Amount:        paymentLogToCreate.Amount,
			Portions:      paymentLogToCreate.Portions,
			PaidBy:        paymentLogToCreate.PaidBy,
			UsersInvolved: paymentLogToCreate.UsersInvolved,
			CreatedAt:     paymentLogToCreate.CreatedAt,
		}
	}

	return newPaymentLog, err
}

func (u *PaymentLogService) UpdatePaymentLog(paymentLog *models.PaymentLog) error {
	paymentLogCollection := u.Db.Collection(connection.PAYMENT_LOGS_COLLECTION)

	_, err := paymentLogCollection.ReplaceOne(context.TODO(), bson.M{"_id": paymentLog.ID}, paymentLog)

	return err
}

func (u *PaymentLogService) RemovePaymentLog(id string) error {
	paymentLogCollection := u.Db.Collection(connection.PAYMENT_LOGS_COLLECTION)
	objId, err := primitive.ObjectIDFromHex(id)

	if err == nil {
		_, err = paymentLogCollection.DeleteOne(context.TODO(), bson.M{"_id": objId})
	}

	return err
}
