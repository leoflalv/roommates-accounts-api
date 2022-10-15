package models

type PaymentLog struct {
	ID          string           `json:"id"`
	Description string           `json:"description"`
	Portions    []PaymentPortion `json:"portions"`
}
