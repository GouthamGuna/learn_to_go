package model

type Users struct {
	ID            string  `json:"id"`
	UserName      string  `json:"username"`
	AccountNumber int64   `json:"account_number"`
	Email         string  `json:"email"`
	PhoneNumber   int64   `json:"phone_number"`
	Pin           int     `json:"pin"`
	IsActive      bool    `json:"is_active"`
	AccountType   string  `json:"account_type"`
	Min_Deposit   Deposit `json:"min_deposit"`
	TotalAmount   int64   `json:"total_amount"`
}
