package response

import (
	"encoding/json"
	"github.com/GouthamGuna/learn_to_go/tree/main/go_atm/utils"
	"net/http"
)

func InvalidPayloadResponse(w http.ResponseWriter, err error) {
	utils.HeaderTypeJson(w, http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request payload.", ErrorType: err})
}

func InvaildPathParamResponse(w http.ResponseWriter, err error) {
	utils.HeaderTypeJson(w, http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid path parameter.", ErrorType: err})
}

func UserAlreadyExistsResponse(w http.ResponseWriter, err error) {
	utils.HeaderTypeJson(w, http.StatusConflict)
	json.NewEncoder(w).Encode(ErrorResponse{Error: "User already exists.", ErrorType: err})
}

func UserNotFoundResponse(w http.ResponseWriter, err error) {
	utils.HeaderTypeJson(w, http.StatusNotFound)
	json.NewEncoder(w).Encode(ErrorResponse{Error: "User not found.", ErrorType: err})
}

func InvaildUpdateResponse(w http.ResponseWriter, err error) {
	utils.HeaderTypeJson(w, http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{Error: "Nothing to be changed.", ErrorType: err})
}

func MiniumDepositResponse(w http.ResponseWriter, err error) {
	utils.HeaderTypeJson(w, http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{Error: "The minimum deposit amount is higher than 500.", ErrorType: err})
}

func MiniumWithdrawalResponse(w http.ResponseWriter, err error) {
	utils.HeaderTypeJson(w, http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{Error: "The minimum withdrawal amount is higher than 500.", ErrorType: err})
}	

func InsufficientBalanceResponse(w http.ResponseWriter, err error) {
	utils.HeaderTypeJson(w, http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{Error: "Insufficient balance.", ErrorType: err})
}

func InvalidPinResponse(w http.ResponseWriter, err error) {
	utils.HeaderTypeJson(w, http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid pin.", ErrorType: err})
}