package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/suhas-developer07/notification-service/cmd/producer"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/notify", producer.ProduceMessage).Methods("POST")

	http.ListenAndServe(":8080", router)

}
