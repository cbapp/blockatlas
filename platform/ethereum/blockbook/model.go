package blockbook

import (
	"math/big"

	"github.com/trustwallet/golibs/tokentype"
)

type Page struct {
	Transactions []Transaction `json:"result,omitempty"`
	Tokens       []Token       `json:"tokens,omitempty"`
}

type NodeInfo struct {
	Blockbook *Blockbook `json:"blockbook"`
	Backend   *Backend   `json:"backend"`
}

type Blockbook struct {
	BestHeight int64 `json:"bestHeight"`
}

type EtherScan struct {
	BlockHeight string `json:"result"`
}

type Backend struct {
	Blocks int64 `json:"blocks"`
}

type Block struct {
	Transactions []Transaction `json:"result"`
}

type Transaction struct {
	BlockNumber       string `json:"blockNumber,omitempty"`
	TimeStamp         string `json:"timeStamp,omitempty"`
	TxID              string `json:"hash,omitempty"`
	Nonce             string `json:"nonce,omitempty"`
	BlockHash         string `json:"blockHash,omitempty"`
	From              string `json:"from,omitempty"`
	ContractAddress   string `json:"contractAddress,omitempty"`
	To                string `json:"to,omitempty"`
	Value             string `json:"value,omitempty"`
	TokenName         string `json:"tokenName,omitempty"`
	TokenSymbol       string `json:"tokenSymbol,omitempty"`
	TokenDecimal      string `json:"tokenDecimal,omitempty"`
	TransactionIndex  string `json:"transactionIndex,omitempty"`
	Gas               string `json:"gas,omitempty"`
	GasPrice          string `json:"gasPrice,omitempty"`
	GasUsed           string `json:"gasUsed,omitempty"`
	CumulativeGasUsed string `json:"cumulativeGasUsed,omitempty"`
	Input             string `json:"input,omitempty"`
	Confirmations     string `json:"confirmations,omitempty"`

	// Always empty, doesn't exist on etherscan
	TokenTransfers []TokenTransfer `json:"tokenTransfers,omitempty"`
}

type Output struct {
	Value     string   `json:"value,omitempty"`
	Addresses []string `json:"addresses"`
}

type TokenTransfer struct {
	Decimals uint   `json:"decimals"`
	From     string `json:"from"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	To       string `json:"to"`
	Token    string `json:"token"`
	Type     string `json:"type"`
	Value    string `json:"value"`
}

// Token contains info about tokens held by an address
type Token struct {
	Balance  string         `json:"balance,omitempty"`
	Contract string         `json:"contract"`
	Decimals uint           `json:"decimals"`
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	Type     tokentype.Type `json:"type"`
}

// EthereumSpecific contains ethereum specific transaction data
type EthereumSpecific struct {
	Status   int      `json:"status"` // -1 pending, 0 Fail, 1 OK
	Nonce    uint64   `json:"nonce"`
	GasLimit *big.Int `json:"gasLimit"`
	GasUsed  *big.Int `json:"gasUsed"`
	GasPrice string   `json:"gasPrice"`
	Data     string   `json:"data,omitempty"`
}