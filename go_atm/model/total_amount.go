package model

type TotalAmount struct {
	UserName      string `json:"user_name"`
	AccountNumber int64  `json:"account_number"`
	TotalAmount   int64  `json:"total_amount"`
	Status        string `json:"status"`
}
