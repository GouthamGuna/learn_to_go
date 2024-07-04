package server

import (
	"fmt"
	"net/http"
)

const ip_ADDRESS string = "http://127.0.0.1"
const sERVER_PORT string = ":2427"

func StartServer() {
	fmt.Printf("Hello Lunar! \nServer Start and Running on ::: %s%s", ip_ADDRESS, sERVER_PORT)

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		DownloadHandler(response, request)
		SetResHTMLContent("public/index.html", response)
	})

	http.HandleFunc("/upload", func(response http.ResponseWriter, request *http.Request) {
		DownloadHandler(response, request)
	})

	http.HandleFunc("/receive", ReceiveHandler)

	http.HandleFunc("/auth/login/", func(response http.ResponseWriter, request *http.Request) {

		request.ParseForm()
		userName := request.FormValue("userName")
		secret := request.FormValue("secret")

		fmt.Print("userName :", userName)
		fmt.Print("\nsecret :", secret)

		if userName == "" || secret == "" {
			SetResHTMLContent("public/badRequest.html", response)
		}

		if userName == "admin" && secret == "Kandha_Saranam" {
			fmt.Fprintf(response, "Hello user")
		} else {
			SetResHTMLContent("public/html/badRequest.html", response)
		}
	})

	http.HandleFunc("/greet/", func(reponse http.ResponseWriter, request *http.Request) {
		userName := request.URL.Path[len("/greet/"):]

		if userName == "" {
			fmt.Fprintf(reponse, "Hello user your using guest mode...")
		} else {
			fmt.Fprintf(reponse, "Hello %s\n", userName)
		}
	})

	http.ListenAndServe(sERVER_PORT, nil)
}
