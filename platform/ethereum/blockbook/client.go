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

// GetBlock fetches all transactions that belongs to the block number specified.
func (c *Client) GetBlock(num int64) (Block, error) {
	var block Block

	// As etherscan doesn't allow us to fetch token transactions in bulk for a single block.
	// we're gonna need to loop through the addresses that interest us.
	addresses := []string{
		"0xFE5Eea7e15d4fD6Ff5F255Eff97803C3e361786D",
		"0xd447d5EAAafd87f15C276EbD10D02337A64aAAd7",
		"0xb82b9174D8b79f9b59C89656156a5A0382AF6ecA",
		"0x074C02D3BB12D5D8d22EA361CF547aFE49BA0480",
		"0x778a928df0bD1DCe363d3C220a315bA16B7e7ddD",
		"0xe734955227d53A1C17c94C55862f2D880408F5f7",
		"0x7Eea51C5a47D4a361fcAD0af7259a4f754989c97",
		"0x26c73EC1Efc273a6b51b2A3d0C589De4a095B961",
		"0x2eAd0EFDE0983bA1aE7B794d49fe5C89040Dea50",
		"0xeAe7679e84e65114F1433be2EEf6eBB0C1279C32",
		"0xac3bD78e33EE4406dC21250A0e321E3637F601bF",
		"0x1Cd2201811749Ec3658c662225F4fCeF17764D26",
		"0x28Ec7DF4Cd88AcDfB3D00C33D61b572BeBa55872",
		"0x296744107B1a83E1f817C2b9051E7338Cb85E076",
		"0xAf1cE0D74b05f3C8A8A5D386C1094dA2b50dfb9b",
		"0x0F8d631aEf351Ad0b6A6d4878002729c14863768",
		"0x69e075d969ad795eb359D2F0910Ec5B9D3Ae1385",
		"0x235Cfb4Ff8f0da903907be760d5A1AD28743D7d9",
		"0x6236966220585C2792f279867dc03380A78836bc",
		"0xfB0aE2042d023FbA129Dbe8035d54E8BbF794093"}

	blockNumStr := fmt.Sprintf("%d \n", num)

	for _, address := range addresses {
		var b Block

		query := url.Values{
			"module":     {"account"},
			"action":     {"tokentx"},
			"startblock": {blockNumStr},
			"endblock":   {blockNumStr},
			"sort":       {"asc"},
			"apiKey":     {apiKey},
			"address":    {address},
		}

		err := c.Get(&b, path, query)
		if err != nil {
			return b, err
		}

		block.Transactions = append(block.Transactions, b.Transactions...)

		// Etherscan only allows 5 calls per sec/IP
		time.Sleep(300 * time.Millisecond)
	}

	return block, nil
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
