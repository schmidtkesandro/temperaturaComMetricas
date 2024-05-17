package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/cep", handleCEP).Methods("POST")

	log.Println("Service B is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
