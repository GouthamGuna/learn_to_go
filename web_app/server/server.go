package server

import (
	"fmt"
	"net/http"
)

var IP_ADDRESS string = "http://127.0.0.1"
var SERVER_PORT string = ":2427"

func StartServer() {

	fmt.Printf("Hello Lunar! \nServer Start and Running on ::: %s%s", IP_ADDRESS, SERVER_PORT)

	const songsDir = "datasource/Chinna_Chinna_Kangal.mp3"

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

	http.HandleFunc("/music", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, songsDir)
	})

	http.HandleFunc("/video", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "datasource/Arupadai Muruga _ Murugan.mp4")
	})

	http.HandleFunc("/mix", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "datasource/mix.json")
	})

	http.ListenAndServe(SERVER_PORT, nil)
}
