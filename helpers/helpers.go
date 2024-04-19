package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	interfaces "trustwallet/interfaces"
)

const ethUrl string = "https://cloudflare-eth.com"

func ToLowerCase(input string) string {
	var result strings.Builder
	for _, char := range input {
		if char >= 'A' && char <= 'Z' {
			result.WriteRune(char + 32)
		} else {
			result.WriteRune(char)
		}
	}
	return result.String()
}

func FetchLatestBlockNumber() (string, error) {
	payload := strings.NewReader(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`)

	req, err := http.NewRequest("POST", ethUrl, payload)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return "", err
	}

	return responseData["result"].(string), nil
}

func GetTransactionsInBlock(blockNumber string) ([]interfaces.Transaction, error) {
	payload := strings.NewReader(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["%s", true],"id":1}`, blockNumber))

	req, err := http.NewRequest("POST", ethUrl, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var blockData map[string]interface{}
	if err := json.Unmarshal(body, &blockData); err != nil {
		return nil, err
	}

	transactionsData, ok := blockData["result"].(map[string]interface{})["transactions"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to parse transactions data")
	}

	var transactions []interfaces.Transaction
	for _, txData := range transactionsData {

		// to can be null
		to, ok := txData.(map[string]interface{})["to"].(string)
		if !ok {
			to = ""
			continue
		}

		tx := interfaces.Transaction{
			Hash:        txData.(map[string]interface{})["hash"].(string),
			From:        txData.(map[string]interface{})["from"].(string),
			To:          to,
			Value:       txData.(map[string]interface{})["value"].(string),
			BlockNumber: txData.(map[string]interface{})["blockNumber"].(string),
			Gas:         txData.(map[string]interface{})["gas"].(string),
			GasPrice:    txData.(map[string]interface{})["gasPrice"].(string),
			Nonce:       txData.(map[string]interface{})["nonce"].(string),
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}
