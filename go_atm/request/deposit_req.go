package request

type DepositRequest struct {
	Pin    int `json:"pin"`
	Amount int64 `json:"amount"`
}
