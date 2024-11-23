package response

type ResponseCounter struct {
	CounterId string `json:"counterId"`
	LastNumber uint `json:"lastNumber"`
}

type ResponseAllSaldo struct {
	Account string `json:"account"`
	Wallet []ResponseSaldoOnly `json:"wallet"`
}

type ResponseSaldoOnly struct {
	Currency string `json:"currency"`
	Balance float64 `json:"balance"`
}

type ResponseHistoryList struct {
	TransactionId string `json:"transaction_id"`
	Account string `json:"account"`
	Currency string `json:"currency"`
	TypeTransaction string `json:"type_transaction"`
	TransactionDate string `json:"transaction_date"`
	Amount float64 `json:"amount"`
}