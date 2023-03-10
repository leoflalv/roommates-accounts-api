package controllers

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/leoflalv/roommates-accounts-api/models"
	"github.com/leoflalv/roommates-accounts-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentLogController struct {
	PaymentLogService models.PaymentLogService
	UserService       models.UserService
}

// .
// GET payment-logs
// .
func (plc PaymentLogController) GetPaymentLogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	paymentLogs, err := plc.PaymentLogService.GetAllPaymentLogs()

	if err != nil {
		utils.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.HttpSuccess(w, http.StatusOK, &paymentLogs)
}

// .
// GET payment-log/:id
// .
func (plc PaymentLogController) GetPaymentLogsByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]
	paymentLog, err := plc.PaymentLogService.GetPaymentLogById(id)

	// Verify if exist a payment log with this id
	if err == mongo.ErrNoDocuments {
		utils.HttpError(w, http.StatusNotFound, "No payment log with this id")
		return
	}

	if err != nil {
		utils.HttpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := utils.Response[models.PaymentLog]{Data: paymentLog, Success: true}
	utils.HttpSuccess(w, http.StatusOK, &resp)
}

// .
// POST payment-log/create
// .
func (plc PaymentLogController) CreatePaymentLogHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var paymentLog models.PaymentLog

	if err := json.NewDecoder(r.Body).Decode(&paymentLog); err != nil {
		utils.HttpError(w, http.StatusBadRequest, "Bad request")
		return
	}

	userWhoPaid, _ := plc.UserService.GetUserById(paymentLog.PaidBy.ID.Hex())

	// For each portion information each user is updated
	for i, portion := range paymentLog.Portions {
		if portion.UserId == userWhoPaid.ID {
			break
		}

		involvedUser, _ := plc.UserService.GetUserById(portion.UserId.Hex())
		paymentLog.Portions[i].UserName = involvedUser.Username
		amount := paymentLog.Amount * portion.Portion

		if debt, found := utils.GetItemById(userWhoPaid.ToPay, involvedUser.ID); found {
			newAmount := amount - debt.Amount
			amount = math.Max(newAmount, 0)

			if newAmount > 0 {

				utils.RemoveItemById(&involvedUser.ToCollect, userWhoPaid.ID.Hex())

				newDebtToPay := models.Debt{UserId: userWhoPaid.ID, UserName: userWhoPaid.Username, Amount: amount}
				involvedUser.ToPay = append(involvedUser.ToPay, newDebtToPay)

				utils.RemoveItemById(&userWhoPaid.ToPay, involvedUser.ID.Hex())

				newDebtToCollect := models.Debt{UserId: involvedUser.ID, UserName: involvedUser.Username, Amount: amount}
				userWhoPaid.ToCollect = append(userWhoPaid.ToCollect, newDebtToCollect)
			} else if newAmount == 0 {
				utils.RemoveItemById(&involvedUser.ToCollect, userWhoPaid.ID.Hex())
				utils.RemoveItemById(&userWhoPaid.ToPay, involvedUser.ID.Hex())
			} else {
				debt.Amount = -newAmount
				utils.UpdateItem(&userWhoPaid.ToPay, *debt)

				newDebtToCollect := models.Debt{UserId: userWhoPaid.ID, UserName: userWhoPaid.Username, Amount: -newAmount}
				utils.UpdateItem(&involvedUser.ToCollect, newDebtToCollect)
			}

		} else if debt, found := utils.GetItemById(userWhoPaid.ToCollect, involvedUser.ID); found {
			newDebtToCollect := models.Debt{UserId: debt.UserId, UserName: debt.UserName, Amount: debt.Amount + amount}
			utils.UpdateItem(&userWhoPaid.ToCollect, newDebtToCollect)

			newDebtToPay := models.Debt{UserId: userWhoPaid.ID, UserName: userWhoPaid.Username, Amount: debt.Amount + amount}
			utils.UpdateItem(&involvedUser.ToPay, newDebtToPay)
		} else {
			newDebtToPay := models.Debt{UserId: userWhoPaid.ID, UserName: userWhoPaid.Username, Amount: amount}
			involvedUser.ToPay = append(involvedUser.ToPay, newDebtToPay)

			newDebtToCollect := models.Debt{UserId: involvedUser.ID, UserName: involvedUser.Username, Amount: amount}
			userWhoPaid.ToCollect = append(userWhoPaid.ToCollect, newDebtToCollect)
		}

		plc.UserService.UpdateUser(involvedUser)
	}

	plc.UserService.UpdateUser(userWhoPaid)
	paymentLog.PaidBy.Username = userWhoPaid.Username
	newPaymentLog, _ := plc.PaymentLogService.CreatePaymentLog(&paymentLog)
	resp := utils.Response[models.PaymentLog]{Success: true, Data: newPaymentLog}
	utils.HttpSuccess(w, http.StatusCreated, &resp)
}

// .
// DELETE user/delete/:id
// .
func (plc PaymentLogController) DeletePaymentLogHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]
	paymentLog, err := plc.PaymentLogService.GetPaymentLogById(id)

	if err == mongo.ErrNoDocuments {
		utils.HttpError(w, http.StatusNotFound, "No documents with this id")
		return
	}

	userWhoPaid, _ := plc.UserService.GetUserById(paymentLog.PaidBy.ID.Hex())

	for _, portion := range paymentLog.Portions {
		if portion.UserId == userWhoPaid.ID {
			break
		}

		involvedUser, _ := plc.UserService.GetUserById(portion.UserId.Hex())
		amount := paymentLog.Amount * portion.Portion

		if debt, found := utils.GetItemById(userWhoPaid.ToPay, involvedUser.ID); found {
			newDebtToPay := models.Debt{UserId: involvedUser.ID, UserName: involvedUser.Username, Amount: debt.Amount + amount}
			utils.UpdateItem(&userWhoPaid.ToPay, newDebtToPay)

			newDebtToCollect := models.Debt{UserId: userWhoPaid.ID, UserName: userWhoPaid.Username, Amount: debt.Amount + amount}
			utils.UpdateItem(&involvedUser.ToCollect, newDebtToCollect)
		} else if debt, found := utils.GetItemById(userWhoPaid.ToCollect, involvedUser.ID); found {
			newAmount := amount - debt.Amount
			amount = math.Max(newAmount, 0)

			if amount > 0 {
				utils.RemoveItemById(&userWhoPaid.ToCollect, involvedUser.ID.Hex())
				newDebtToPay := models.Debt{UserId: involvedUser.ID, UserName: involvedUser.Username, Amount: debt.Amount + amount}
				userWhoPaid.ToPay = append(userWhoPaid.ToPay, newDebtToPay)

				utils.RemoveItemById(&involvedUser.ToPay, userWhoPaid.ID.Hex())
				newDebtToCollect := models.Debt{UserId: userWhoPaid.ID, UserName: userWhoPaid.Username, Amount: debt.Amount + amount}
				involvedUser.ToCollect = append(involvedUser.ToCollect, newDebtToCollect)
			} else if newAmount == 0 {
				utils.RemoveItemById(&userWhoPaid.ToCollect, involvedUser.ID.Hex())
				utils.RemoveItemById(&involvedUser.ToPay, involvedUser.ID.Hex())
			} else {
				debt.Amount = -newAmount
				utils.UpdateItem(&userWhoPaid.ToCollect, *debt)

				newDebtToPay := models.Debt{UserId: userWhoPaid.ID, UserName: userWhoPaid.Username, Amount: -newAmount}
				utils.UpdateItem(&involvedUser.ToPay, newDebtToPay)
			}

		} else {
			newDebtToPay := models.Debt{UserId: involvedUser.ID, UserName: involvedUser.Username, Amount: amount}
			userWhoPaid.ToPay = append(userWhoPaid.ToPay, newDebtToPay)

			newDebtToCollect := models.Debt{UserId: userWhoPaid.ID, UserName: userWhoPaid.Username, Amount: amount}
			involvedUser.ToCollect = append(involvedUser.ToCollect, newDebtToCollect)
		}
		plc.UserService.UpdateUser(involvedUser)

	}
	plc.UserService.UpdateUser(userWhoPaid)

	paymentLog.PaidBy.Username = userWhoPaid.Username
	paymentLog.DeletedAt = time.Now()

	error := plc.PaymentLogService.UpdatePaymentLog(&paymentLog)

	if error != nil {
		utils.HttpError(w, http.StatusInternalServerError, error.Error())
		return
	}

	resp := utils.Response[struct{}]{Success: true}
	utils.HttpSuccess(w, http.StatusOK, &resp)
}

//
// UPDATE user/update
//
// func (uc PaymentLogController) UpdatePaymentLogHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	var user models.PaymentLog
// 	var resp utils.Response[models.PaymentLog]

// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		resp = utils.Response[models.PaymentLog]{Success: false, Errors: "Bad request"}
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	err := uc.PaymentLogService.UpdatePaymentLog(&user)

// 	if err == mongo.ErrNoDocuments {
// 		resp = utils.Response[models.PaymentLog]{Success: false, Errors: "No documents with this id"}
// 		w.WriteHeader(http.StatusNotFound)
// 	}

// 	if err != nil {
// 		resp = utils.Response[models.PaymentLog]{Success: false, Errors: "Something went wrong"}
// 		w.WriteHeader(http.StatusInternalServerError)
// 	} else {
// 		resp = utils.Response[models.PaymentLog]{Data: user, Success: true}
// 		w.WriteHeader(http.StatusOK)
// 	}

// 	jsonResponse, _ := json.Marshal(resp)
// 	w.Write(jsonResponse)
// }
