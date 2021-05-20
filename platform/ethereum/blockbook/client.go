package blockbook

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/trustwallet/blockatlas/pkg/blockatlas"
)

const (
	// apiKey on etherscan ropsten.
	apiKey = "Z85YCX6M4FFKUFKGTWDFYWHNPRZ823J4NI"

	// path
	path = "api"

	// address
	address = "0xFE5Eea7e15d4fD6Ff5F255Eff97803C3e361786D"
)

type Client struct {
	blockatlas.Request
}

func (c *Client) GetTxs(address string) (*Page, error) {
	return c.getTransactions(address, "")
}

func (c *Client) GetTxsWithContract(address, contract string) (*Page, error) {
	return c.getTransactions(address, contract)
}

func (c *Client) GetTokens(address string) ([]Token, error) {
	return c.getTokens(address)
}

func (c *Client) GetCurrentBlockNumber() (int64, error) {
	timestampStr := fmt.Sprintf("%d", time.Now().UTC().Unix())

	query := url.Values{
		"module":    {"block"},
		"action":    {"getblocknobytime"},
		"timestamp": {timestampStr},
		"apikey":    {apiKey},
		"closest":   {"before"},
	}

	var reqResult EtherScan

	err := c.Get(&reqResult, path, query)
	if err != nil {
		return 0, err
	}

	res, err := strconv.ParseInt(reqResult.BlockHeight, 10, 64)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (c *Client) GetBlock(num int64) (Block, error) {
	var b Block

	blockNumStr := fmt.Sprintf("%d", num)

	query := url.Values{
		"module":     {"account"},
		"action":     {"tokentx"},
		"startblock": {blockNumStr},
		"endblock":   {blockNumStr},
		"sort":       {"asc"},
		"apiKey":     {apiKey},
		"address":    {address},
	}

	fmt.Printf("😀 %s \n", query.Encode())

	err := c.Get(&b, path, query)
	if err != nil {
		return b, err
	}

	fmt.Printf("🚨 %+v \n", b)

	return b, nil
}

func (c *Client) getTransactions(address, contract string) (page *Page, err error) {

	query := url.Values{
		"module":     {"account"},
		"address":    {address},
		"apikey":     {apiKey},
		"startblock": {"0"},
		"endblock":   {"999999999"},
		"sort":       {"desc"},
		"page":       {"1"},
		"offset":     {"25"},
	}

	if contract != "" {
		query.Add("action", "tokentx")
	} else {
		query.Add("action", "txlist")
	}

	err = c.Get(&page, path, query)
	return
}

func (c *Client) getTokens(address string) ([]Token, error) {
	// TODO
	return []Token{}, nil
}
