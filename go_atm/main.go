package main

import (
	"log"
	"net/http"

	service "github.com/GouthamGuna/learn_to_go/tree/main/go_atm/service"
)

const SERVER_PORT = ":8083"

func main() {
	userService := &service.Handler{}

	router := http.NewServeMux()
	router.HandleFunc("POST /create-user", userService.AddUsers)
	router.HandleFunc("GET /get-all-user", userService.GetAllUsers)
	router.HandleFunc("PUT /update-user", userService.UpdateUser)
	router.HandleFunc("GET /balance", userService.FetchBalance)
	router.HandleFunc("POST /deposit", userService.Deposit)
	router.HandleFunc("POST /withdrawal", userService.Withdrawal)

	server := http.Server{
		Addr:    SERVER_PORT,
		Handler: router,
	}

	log.Printf("Server started at port %s\n", SERVER_PORT)
	server.ListenAndServe()
}
