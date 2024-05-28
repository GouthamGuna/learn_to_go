package main

import (
	"fmt"
	"net/http"
)

var IP_ADDRESS string = "http://127.0.0.1"
var SERVER_PORT string = ":2427"

func main() {

	fmt.Printf("Hello Lunar! \nServer Start and Running on ::: %s%s", IP_ADDRESS, SERVER_PORT)

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(response, "Hello Lunar From Golang Server...")
	})

	http.HandleFunc("/greet/", func(reponse http.ResponseWriter, request *http.Request) {
		userName := request.URL.Path[len("/greet/"):]

		if userName == "" {
			fmt.Fprintf(reponse, "Hello user your using guest mode...")
		} else {
			fmt.Fprintf(reponse, "Hello %s\n", userName)
		}

	})

	http.ListenAndServe(SERVER_PORT, nil)
}
