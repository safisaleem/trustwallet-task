For this task here are some details:

The code has 4 key components;
1. The Parser
2. The observer (its the FetchLatestTransactions() function in the parser), which runs as a goroutine in main
3. The rest api interface for external use - this is for simplicity. ideally we could use a queue/websocket/external webhook to send events to the notifcation service
4. Helper functions to make calls to the JSONRPC address


The Parser contains:
- an in memory list of subscribed addresses
- an in memory list of transactions that have already been fetched
- the latest block number
- GetLatestBlockNumber() which returns the stored latest block number
- Subscribe(address string) which adds a new address to the subscription list
- GetTransactions(address string) which returns the transactions that have been fetched
- FetchLatestTransactions() which is the observer that runs every second to check if the current block is the latest block on the ethereum network, and if not, then it updates the latest block number and then fetches all the transactions from that block, filters them and adds the transactions from the subscribed list into our in memory stored transactions

There are three rest endpoints that allow the use of this parser.
For simplicity, I have used rest endpoints, which could be pinged by the notification service to check for new transactions.

Alternately, the parser could be modified and a queue or a websocket or external webhook could be added where it emits events of new transactions so the notification service can use it.
