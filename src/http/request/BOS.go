package request

import "time"

type DataCounter struct {
	CounterId string `json:"counterId" validate:"required"`
	LastNumber uint `json:"lastNumber" validate:"required,numeric"`
}

type DataBalance struct {
	Account string `json:"account" validate:"required"`
	Currency string `json:"currency" validate:"required,oneof=USD IDR SGD"`
	Balance float64 `json:"balance" validate:"required,number"`
}

type DataRequestHistory struct {
	CounterId string `json:"counterId" validate:"required"`
	Account string `json:"account" validate:"required"`
	Currency string `json:"currency" validate:"required"`
	Amount float64 `json:"amount" validate:"required,min=0"`
	TypeTransaction string `json:"type" validate:"required,oneof=SETOR TARIK TRANSFER"`
	DateTransaction *time.Time `json:"date_transaction"`
}

type RequestTransaction struct {
	FromAccount string `json:"from_account" validate:"required"`
	Details []RequestTransactionDetail `json:"details" validate:"dive"`
}
type RequestTransactionDetail struct {
	ToAccount string `json:"to_account" validate:"required"`
	Currency string `json:"currency" validate:"required"`
	Amount float64 `json:"amount" validate:"required,min=0"`
}

type RequestSelfTransaction struct {
	Account string `json:"account" validate:"required"`
	Details []RequestSelfTransactionDetail `json:"details" validate:"dive"`
}
type RequestSelfTransactionDetail struct {
	Currency string `json:"currency" validate:"required"`
	Amount float64 `json:"amount" validate:"required,min=0"`
}


type RequestHistoryList struct {
	Account string `json:"account" validate:"required"`
	From time.Time `json:"from" validate:"required"`
	To *time.Time `json:"to"`
}

type PrepareRequestTransaction struct {
	Account string `json:"account"`
	Detail []PrepareRequestTransactionDetail `json:"detail"`
}
type PrepareRequestTransactionDetail struct {
	Amount float64 `json:"balance"`
	Transaction string `json:"type_transaction"`
	Currency string `json:"currency"`
}