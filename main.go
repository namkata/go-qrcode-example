package main

import (
	"net/http"
	"qr-code-generator/handlers"
)

func main() {
	http.HandleFunc("/generate", handlers.HandleRequest)
	http.ListenAndServe(":8080", nil)
}
