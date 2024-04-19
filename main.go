package main

import (
	"log"
	"net/http"
	"trustwallet/controller"
	"trustwallet/models"
)

func main() {
	parser := models.NewParser()        // initiate the parser
	go parser.FetchLatestTransactions() // fire up the observer

	controller := controller.NewController(parser) // pass the parser to the controller

	http.HandleFunc("/block", controller.GetLatestBlockNumber)
	http.HandleFunc("/subscribe", controller.SubscribeAddress)
	http.HandleFunc("/transactions", controller.GetTransactions)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
