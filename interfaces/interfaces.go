package interfaces

type Transaction struct {
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	BlockNumber string `json:"blockNumber"`
	Gas         string `json:"gas"`
	GasPrice    string `json:"gasPrice"`
	Nonce       string `json:"nonce"`
}

type BlockNumber struct {
	BlockNumber string `json:"blockNumber"`
}

type AddressInput struct {
	Address string `json:"address"`
}
