package main

import (
	"log"
	"net/http"
)

func main() {
	router := createRouter()
	log.Fatal(http.ListenAndServe(":12345", router))
}
