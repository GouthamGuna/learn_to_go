package users

import (
	"encoding/json"
	"net/http"

	"github.com/GouthamGuna/learn_to_go/tree/main/go_atm/model"
	"github.com/GouthamGuna/learn_to_go/tree/main/go_atm/request"
	"github.com/GouthamGuna/learn_to_go/tree/main/go_atm/response"
	"github.com/GouthamGuna/learn_to_go/tree/main/go_atm/utils"
	"github.com/google/uuid"
)

type Handler struct{}
type Users = model.Users
type balanceStmt = request.BalanceStmt
type depositRequest = request.DepositRequest
type withdrawalRequest = request.WithdrawalRequest

// Global slice to store users
var users = []model.Users{}

const (
	MINIUM_DEPOSIT_AMOUNT    = 500
	MINIUM_WITHDRAWAL_AMOUNT = 500
)

// AddUsers handles the HTTP request for adding users
func (h *Handler) AddUsers(w http.ResponseWriter, r *http.Request) {
	var newUser Users

	err := json.NewDecoder(r.Body).Decode(&newUser)
	//fmt.Println("Incoming JSON Payload:", newUser)

	if err != nil {
		response.InvalidPayloadResponse(w, err)
		return
	}

	if newUser.Min_Deposit.Amount < MINIUM_DEPOSIT_AMOUNT {
		response.MiniumDepositResponse(w, nil)
		return
	}

	if UserDataValidater(newUser.AccountNumber, newUser.Pin) {
		response.UserAlreadyExistsResponse(w, nil)
		return
	}

	// Add the new user to the list
	newUser.ID = uuid.New().String()
	newUser.IsActive = true
	newUser.TotalAmount = newUser.Min_Deposit.Amount
	users = append(users, newUser)

	// Return the added user as a response
	utils.HeaderTypeJson(w, http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	utils.HeaderTypeJson(w, http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"all_users": users})
}

func UserDataValidater(accountNumber int64, pin int) bool {
	for _, u := range users {
		if u.AccountNumber == accountNumber && u.Pin == pin {
			return true
		}
	}
	return false
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var updatedUser Users
	var resourceNotFound bool = true

	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		response.InvalidPayloadResponse(w, err)
		return
	}

	for i, u := range users {
		if u.ID == updatedUser.ID {

			if users[i] == updatedUser {
				response.InvaildUpdateResponse(w, nil)
				return
			}

			users[i] = updatedUser
			resourceNotFound = false
			break
		}
	}

	if resourceNotFound {
		response.UserNotFoundResponse(w, nil)
		return
	}

	utils.HeaderTypeJson(w, http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}

func (h *Handler) FetchBalance(w http.ResponseWriter, r *http.Request) {
	var b_stmt balanceStmt
	var resourceNotFound bool = true

	err := json.NewDecoder(r.Body).Decode(&b_stmt)
	if err != nil {
		response.InvalidPayloadResponse(w, err)
		return
	}

	for _, u := range users {
		if u.Pin == b_stmt.Pin {
			resourceNotFound = false
			utils.HeaderTypeJson(w, http.StatusOK)
			result := model.TotalAmount{
				UserName:      u.UserName,
				AccountNumber: u.AccountNumber,
				TotalAmount:   u.TotalAmount,
				Status: 	  "Successfully fetched balance.",
			}
			json.NewEncoder(w).Encode(result)
			break
		}
	}

	if resourceNotFound {
		response.InvalidPinResponse(w, nil)
		return
	}
}

/*

func (h *Handler) FetchBalance(w http.ResponseWriter, r *http.Request) {
	var resourceNotFound bool = true

	pin := r.URL.Query().Get("pin")
	fmt.Println("Pin = ", pin)
	if pin == "" {
		response.InvaildPathParamResponse(w, errors.New("pin is required"))
		return
	}

	for _, u := range users {
		if u.Pin == pin { // fmt.Sprintf("%d", u.Pin) == pin
			resourceNotFound = false
			utils.HeaderTypeJson(w, http.StatusOK)
			result := model.TotalAmount{
				UserName:      u.UserName,
				AccountNumber: u.AccountNumber,
				TotalAmount:   u.TotalAmount,
				Status: 	  "Successfully fetched balance.",
			}
			json.NewEncoder(w).Encode(result)
			break
		}
	}

	if resourceNotFound {
		response.InvalidPinResponse(w, nil)
		return
	}
}

*/

func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	var depositReq depositRequest
	var resourceNotFound bool = true

	err := json.NewDecoder(r.Body).Decode(&depositReq)
	if err != nil {
		response.InvalidPayloadResponse(w, err)
		return
	}

	if depositReq.Amount < MINIUM_DEPOSIT_AMOUNT {
		response.MiniumDepositResponse(w, nil)
		return
	}

	for i, u := range users {
		if u.Pin == depositReq.Pin {
			users[i].TotalAmount += depositReq.Amount
			resourceNotFound = false
			utils.HeaderTypeJson(w, http.StatusOK)
			result := model.TotalAmount{
				UserName:      users[i].UserName,
				AccountNumber: users[i].AccountNumber,
				TotalAmount:   users[i].TotalAmount,
				Status: 	  "Successfully deposited.",
			}
			json.NewEncoder(w).Encode(result)
			break
		}
	}

	if resourceNotFound {
		response.InvalidPinResponse(w, nil)
		return
	}
}

func (h *Handler) Withdrawal(w http.ResponseWriter, r *http.Request) {
	var withdrawalReq withdrawalRequest
	var resourceNotFound bool = true

	err := json.NewDecoder(r.Body).Decode(&withdrawalReq)
	if err != nil {
		response.InvalidPayloadResponse(w, err)
		return
	}

	if withdrawalReq.Amount < MINIUM_WITHDRAWAL_AMOUNT {
		response.MiniumWithdrawalResponse(w, nil)
		return
	}

	for i, u := range users {
		if u.Pin == withdrawalReq.Pin {
			if users[i].TotalAmount < withdrawalReq.Amount {
				response.InsufficientBalanceResponse(w, nil)
				return
			}
			users[i].TotalAmount -= withdrawalReq.Amount
			resourceNotFound = false
			utils.HeaderTypeJson(w, http.StatusOK)
			result := model.TotalAmount{
				UserName:      users[i].UserName,
				AccountNumber: users[i].AccountNumber,
				TotalAmount:   users[i].TotalAmount,
				Status: 	  "Successfully withdrawn.",
			}
			json.NewEncoder(w).Encode(result)
			break
		}
	}

	if resourceNotFound {
		response.InvalidPinResponse(w, nil)
		return
	}
}
