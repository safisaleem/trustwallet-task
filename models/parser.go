package models

import (
	"fmt"
	"sync"
	"time"
	"trustwallet/helpers"
	"trustwallet/interfaces"
)

type Parser struct {
	subscribed        map[string]bool
	transactions      []interfaces.Transaction
	latestBlockNumber string
	transactionsLock  sync.RWMutex
	blockNumberLock   sync.RWMutex
}

func NewParser() *Parser {
	return &Parser{
		subscribed: make(map[string]bool),
	}
}

func (p *Parser) GetLatestBlockNumber() string {
	p.blockNumberLock.RLock()
	defer p.blockNumberLock.RUnlock()
	return p.latestBlockNumber
}

func (p *Parser) Subscribe(address string) bool {
	p.subscribed[helpers.ToLowerCase(address)] = true
	fmt.Println(p.subscribed)
	return true
}

func (p *Parser) GetTransactions(address string) []interfaces.Transaction {
	p.transactionsLock.RLock()
	defer p.transactionsLock.RUnlock()

	lowerCaseAddress := helpers.ToLowerCase(address)

	var filteredTransactions []interfaces.Transaction
	for _, tx := range p.transactions {
		if tx.From == lowerCaseAddress || tx.To == lowerCaseAddress {
			filteredTransactions = append(filteredTransactions, tx)
		}
	}
	return filteredTransactions
}

// This is the watcher for transactions. It constantly pings the ethereum network
// every second to fetch the latest block and stores the transactions for the
// subscribed addresses in memory
func (p *Parser) FetchLatestTransactions() {
	for {
		// here we get the latest block number
		blockNumber, err := helpers.FetchLatestBlockNumber()
		if err != nil {
			fmt.Println("Error fetching latest block number:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// if the block is already fetched, we skip
		if p.GetLatestBlockNumber() == blockNumber {
			time.Sleep(1 * time.Second)
			continue
		}

		// update the latest block number
		p.blockNumberLock.Lock()
		p.latestBlockNumber = blockNumber
		p.blockNumberLock.Unlock()

		transactions, err := helpers.GetTransactionsInBlock(blockNumber)
		if err != nil {
			fmt.Println("Error fetching transactions for block:", blockNumber, err)
			time.Sleep(1 * time.Second)
			continue
		}

		var filteredTransactions []interfaces.Transaction
		for _, tx := range transactions {

			if p.subscribed[tx.To] || p.subscribed[tx.From] {
				filteredTransactions = append(filteredTransactions, tx)
			}
		}

		p.transactionsLock.Lock()
		p.transactions = append(p.transactions, filteredTransactions...)
		p.transactionsLock.Unlock()

		time.Sleep(1 * time.Second)
	}
}
