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
		sendResJsonContent("datasource/json/videos.json", response)
	})

	http.HandleFunc("/stream/{filename}", handleToStream)

	http.ListenAndServe(sERVER_PORT, nil)
}
