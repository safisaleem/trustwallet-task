package controller

import (
	"encoding/json"
	"net/http"
	"trustwallet/interfaces"
	"trustwallet/models"
)

type Controller struct {
	parser *models.Parser
}

func NewController(parser *models.Parser) *Controller {
	return &Controller{parser: parser}
}

func (c *Controller) GetLatestBlockNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	blockNumber := c.parser.GetLatestBlockNumber()

	var blockNumberResponse interfaces.BlockNumber = interfaces.BlockNumber{
		BlockNumber: blockNumber,
	}

	responseJSON, err := json.Marshal(blockNumberResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(responseJSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *Controller) SubscribeAddress(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var input interfaces.AddressInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid Input", http.StatusBadRequest)
		return
	}

	subscribed := c.parser.Subscribe(input.Address)

	if subscribed {
		w.WriteHeader(http.StatusOK)
		return
	} else {
		http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
		return
	}
}

func (c *Controller) GetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	addressInput := r.URL.Query().Get("address")
	if addressInput == "" {
		http.Error(w, "Missing address", http.StatusBadRequest)
		return
	}

	transactions := c.parser.GetTransactions(addressInput)

	responseJSON, err := json.Marshal(transactions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(responseJSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
