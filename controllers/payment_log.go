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
	var resp Response[[]models.PaymentLog]

	if err != nil {
		resp = Response[[]models.PaymentLog]{Success: false, Errors: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp = Response[[]models.PaymentLog]{Data: paymentLogs, Success: true}
		w.WriteHeader(http.StatusOK)
	}

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}

// .
// GET payment-log/:id
// .
func (plc PaymentLogController) GetPaymentLogsByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]
	var resp Response[models.PaymentLog]
	paymentLog, err := plc.PaymentLogService.GetPaymentLogById(id)

	if err == mongo.ErrNoDocuments {
		resp = Response[models.PaymentLog]{Success: false, Errors: "No payment log with this id"}
		w.WriteHeader(http.StatusNotFound)
	}

	if err != nil {
		resp = Response[models.PaymentLog]{Success: false, Errors: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
	} else {
		resp = Response[models.PaymentLog]{Data: paymentLog, Success: true}
		w.WriteHeader(http.StatusOK)
	}

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}

// .
// POST payment-log/create
// .
func (plc PaymentLogController) CreatePaymentLogHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var paymentLog models.PaymentLog
	var resp Response[models.PaymentLog]

	if err := json.NewDecoder(r.Body).Decode(&paymentLog); err != nil {
		resp = Response[models.PaymentLog]{Success: false, Errors: "Bad request"}
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse, _ := json.Marshal(resp)
		w.Write(jsonResponse)
	}

	userWhoPaid, _ := plc.UserService.GetUserById(paymentLog.PaidBy.ID.Hex())

	// For each portion information each user is updated
	for i, portion := range paymentLog.Portions {
		if portion.UserId == userWhoPaid.ID {
			break
		}

		involvedUser, _ := plc.UserService.GetUserById(portion.UserId.Hex())
		paymentLog.Portions[i].UserName = involvedUser.Name
		amount := paymentLog.Amount * portion.Portion

		if debt, found := utils.GetItemById(userWhoPaid.ToPay, involvedUser.ID); found {
			newAmount := amount - debt.Amount
			amount = math.Max(newAmount, 0)

			if newAmount > 0 {

				utils.RemoveItemById(&involvedUser.ToCollect, userWhoPaid.ID.Hex())

				newDebtToPay := models.Debt{UserId: userWhoPaid.ID, UserName: userWhoPaid.Name, Amount: amount}
				involvedUser.ToPay = append(involvedUser.ToPay, newDebtToPay)

				utils.RemoveItemById(&userWhoPaid.ToPay, involvedUser.ID.Hex())

				newDebtToCollect := models.Debt{UserId: involvedUser.ID, UserName: involvedUser.Name, Amount: amount}
				userWhoPaid.ToCollect = append(userWhoPaid.ToCollect, newDebtToCollect)
			} else if newAmount == 0 {

				utils.RemoveItemById(&involvedUser.ToCollect, userWhoPaid.ID.Hex())
				utils.RemoveItemById(&userWhoPaid.ToPay, involvedUser.ID.Hex())
			} else {

				debt.Amount = -newAmount
				utils.UpdateItem(&userWhoPaid.ToPay, *debt)

				newDebtToCollect := models.Debt{UserId: userWhoPaid.ID, UserName: userWhoPaid.Name, Amount: -newAmount}
				utils.UpdateItem(&involvedUser.ToCollect, newDebtToCollect)
			}

		} else if debt, found := utils.GetItemById(userWhoPaid.ToCollect, involvedUser.ID); found {
			newDebtToCollect := models.Debt{UserId: debt.UserId, UserName: debt.UserName, Amount: debt.Amount + amount}
			utils.UpdateItem(&userWhoPaid.ToCollect, newDebtToCollect)

			newDebtToPay := models.Debt{UserId: userWhoPaid.ID, UserName: userWhoPaid.Name, Amount: debt.Amount + amount}
			utils.UpdateItem(&involvedUser.ToPay, newDebtToPay)
		} else {
			newDebtToPay := models.Debt{UserId: userWhoPaid.ID, UserName: userWhoPaid.Name, Amount: amount}
			involvedUser.ToPay = append(involvedUser.ToPay, newDebtToPay)

			newDebtToCollect := models.Debt{UserId: involvedUser.ID, UserName: involvedUser.Name, Amount: amount}
			userWhoPaid.ToCollect = append(userWhoPaid.ToCollect, newDebtToCollect)
		}

		plc.UserService.UpdateUser(involvedUser)
	}

	plc.UserService.UpdateUser(userWhoPaid)

	paymentLog.PaidBy.Username = userWhoPaid.Name
	newPaymentLog, _ := plc.PaymentLogService.CreatePaymentLog(&paymentLog)
	resp = Response[models.PaymentLog]{Success: true, Data: newPaymentLog}
	w.WriteHeader(http.StatusOK)

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}

// .
// DELETE user/delete/:id
// .
func (plc PaymentLogController) DeletePaymentLogHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]
	var resp Response[struct{}]
	paymentLog, err := plc.PaymentLogService.GetPaymentLogById(id)

	if err == mongo.ErrNoDocuments {
		resp = Response[struct{}]{Success: false, Errors: "No documents with this id"}
		w.WriteHeader(http.StatusNotFound)
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
			newDebtToPay := models.Debt{UserId: involvedUser.ID, UserName: involvedUser.Name, Amount: debt.Amount + amount}
			utils.UpdateItem(&userWhoPaid.ToPay, newDebtToPay)

			newDebtToCollect := models.Debt{UserId: userWhoPaid.ID, UserName: userWhoPaid.Name, Amount: debt.Amount + amount}
			utils.UpdateItem(&involvedUser.ToCollect, newDebtToCollect)
		} else if debt, found := utils.GetItemById(userWhoPaid.ToCollect, involvedUser.ID); found {
			newAmount := amount - debt.Amount
			amount = math.Max(newAmount, 0)

			if amount > 0 {
				utils.RemoveItemById(&userWhoPaid.ToCollect, involvedUser.ID.Hex())
				newDebtToPay := models.Debt{UserId: involvedUser.ID, UserName: involvedUser.Name, Amount: debt.Amount + amount}
				userWhoPaid.ToPay = append(userWhoPaid.ToPay, newDebtToPay)

				utils.RemoveItemById(&involvedUser.ToPay, userWhoPaid.ID.Hex())
				newDebtToCollect := models.Debt{UserId: userWhoPaid.ID, UserName: userWhoPaid.Name, Amount: debt.Amount + amount}
				involvedUser.ToCollect = append(involvedUser.ToCollect, newDebtToCollect)
			} else if newAmount == 0 {
				utils.RemoveItemById(&userWhoPaid.ToCollect, involvedUser.ID.Hex())
				utils.RemoveItemById(&involvedUser.ToPay, involvedUser.ID.Hex())
			} else {
				debt.Amount = -newAmount
				utils.UpdateItem(&userWhoPaid.ToCollect, *debt)

				newDebtToPay := models.Debt{UserId: userWhoPaid.ID, UserName: userWhoPaid.Name, Amount: -newAmount}
				utils.UpdateItem(&involvedUser.ToPay, newDebtToPay)
			}

		} else {
			newDebtToPay := models.Debt{UserId: involvedUser.ID, UserName: involvedUser.Name, Amount: amount}
			userWhoPaid.ToPay = append(userWhoPaid.ToPay, newDebtToPay)

			newDebtToCollect := models.Debt{UserId: userWhoPaid.ID, UserName: userWhoPaid.Name, Amount: amount}
			involvedUser.ToCollect = append(involvedUser.ToCollect, newDebtToCollect)
		}
		plc.UserService.UpdateUser(involvedUser)

	}
	plc.UserService.UpdateUser(userWhoPaid)

	paymentLog.PaidBy.Username = userWhoPaid.Name
	paymentLog.DeletedAt = time.Now()

	error := plc.PaymentLogService.UpdatePaymentLog(&paymentLog)

	if error != nil {
		resp = Response[struct{}]{Success: false}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		resp = Response[struct{}]{Success: true}
		w.WriteHeader(http.StatusOK)
	}

	jsonResponse, _ := json.Marshal(resp)
	w.Write(jsonResponse)
}

//
// UPDATE user/update
//
// func (uc PaymentLogController) UpdatePaymentLogHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	var user models.PaymentLog
// 	var resp Response[models.PaymentLog]

// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		resp = Response[models.PaymentLog]{Success: false, Errors: "Bad request"}
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	err := uc.PaymentLogService.UpdatePaymentLog(&user)

// 	if err == mongo.ErrNoDocuments {
// 		resp = Response[models.PaymentLog]{Success: false, Errors: "No documents with this id"}
// 		w.WriteHeader(http.StatusNotFound)
// 	}

// 	if err != nil {
// 		resp = Response[models.PaymentLog]{Success: false, Errors: "Something went wrong"}
// 		w.WriteHeader(http.StatusInternalServerError)
// 	} else {
// 		resp = Response[models.PaymentLog]{Data: user, Success: true}
// 		w.WriteHeader(http.StatusOK)
// 	}

// 	jsonResponse, _ := json.Marshal(resp)
// 	w.Write(jsonResponse)
// }
