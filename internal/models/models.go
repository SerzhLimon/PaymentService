package models

type WalletTransaction struct {
	WalletID  string `json:"walletId"`
	Operation string `json:"operation"`
	Amount    int64  `json:"amount"`
}

type GetBalanceResponse struct {
	Amount int64 `json:"amount"`
}
