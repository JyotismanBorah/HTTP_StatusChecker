package main

import (
	"fmt"
	"main/handler"
	"net/http"
)

func main() {
	fmt.Println("Starting server...")
	http.HandleFunc("/websites", handler.ReqHandler)
	go handler.UpdateStatus(handler.GetMap())
	http.ListenAndServe("127.0.0.1:8080", nil)
}
