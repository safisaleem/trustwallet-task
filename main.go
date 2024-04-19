package main

import (
	"log"
	"net/http"
	"trustwallet/controller"
	"trustwallet/models"
)

func main() {
	parser := models.NewParser()
	go parser.FetchLatestTransactions()

	controller := controller.NewController(parser)

	http.HandleFunc("/block", controller.GetLatestBlockNumber)
	http.HandleFunc("/subscribe", controller.SubscribeAddress)
	http.HandleFunc("/transactions", controller.GetTransactions)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
