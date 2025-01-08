package request

type WithdrawalRequest struct {
	Pin    int   `json:"pin"`
	Amount int64 `json:"amount"`
}
