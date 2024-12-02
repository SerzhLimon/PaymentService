package models

type WalletTransaction struct {
	WalletID  string `json:"wallet_id"`
	Operation string `json:"operation"`
	Amount    int64  `json:"amount"`
}

type GetBalanceResponse struct {
	Amount float64 `json:"balance"`
}
