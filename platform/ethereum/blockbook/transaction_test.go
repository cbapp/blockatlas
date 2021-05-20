package blockbook

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
)

func TestNormalizePage(t *testing.T) {
	type args struct {
		srcPage   string
		address   string
		token     string
		coinIndex uint
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test normalize etherscan txs",
			args: args{
				srcPage: `{
					"result": [
						{
							"blockNumber": "9814035",
							"timeStamp": "1615437064",
							"hash": "0x14d4277bfe58b11882a90b6ca41c1e08954cf57f57b2835585ac93118867c74f",
							"nonce": "157",
							"blockHash": "0x378dd4699d36e63e0e254eecf13f5514bbe3250e03a844f27bdd16b24330bfd0",
							"from": "0x074c02d3bb12d5d8d22ea361cf547afe49ba0480",
							"to": "0xfe5eea7e15d4fd6ff5f255eff97803c3e361786d",
							"contractAddress": "0x0d045476b552982f167f382ef505d2512bddc3b3",
							"value": "9300",
							"tokenName": "TetherToken",
							"tokenSymbol": "USDT",
							"tokenDecimal": "6",
							"transactionIndex": "10",
							"gas": "69123",
							"gasPrice": "1708100000",
							"gasUsed": "46082",
							"cumulativeGasUsed": "1371076",
							"input": "deprecated",
							"confirmations": "456344"
						}
					]
				}`,
				address:   "0xfe5eea7e15d4fd6ff5f255eff97803c3e361786d",
				token:     "",
				coinIndex: 60,
			},
			want: `[
					{
						"id": "0x14d4277bfe58b11882a90b6ca41c1e08954cf57f57b2835585ac93118867c74f",
						"coin": 60,
						"from": "0x074c02d3bb12d5d8d22ea361cf547afe49ba0480",
						"to": "0xfe5eea7e15d4fd6ff5f255eff97803c3e361786d",
						"fee": "78712664200000",
						"date": 1615437064,
						"block": 9814035,
						"status": "completed",
						"sequence": 157,
						"type": "contract_call",
						"direction": "incoming",
						"memo": "",
						"metadata": {
							"input": "0x",
							"name": "TetherToken",
							"symbol": "USDT",
							"token_id": "0x0d045476b552982f167f382ef505d2512bddc3b3",
							"decimals": 6,
							"value": "9300",
							"from": "0x7d8bf18C7cE84b3E175b339c4Ca93aEd1dD166F1",
							"to": "0x7d8bf18C7cE84b3E175b339c4Ca93aEd1dD166F1"
						}
					}
				  ]`,
		},
	}
	for _, tt := range tests {
		var page Page
		var txPage blockatlas.TxPage
		err := json.Unmarshal([]byte(tt.args.srcPage), &page)
		assert.Nil(t, err)
		err = json.Unmarshal([]byte(tt.want), &txPage)
		assert.Nil(t, err)
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizePage(&page, tt.args.address, tt.args.token, tt.args.coinIndex)
			gotJson, err := json.Marshal(got)
			assert.Nil(t, err)
			gotTxPage, err := json.Marshal(txPage)
			assert.Nil(t, err)
			if string(gotJson) != string(gotTxPage) {
				t.Errorf("NormalizePage() = %v, want %v", string(gotJson), string(gotTxPage))
			}
		})
	}
}
