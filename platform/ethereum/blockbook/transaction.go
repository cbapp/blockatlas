package blockbook

import (
	"strconv"

	Address "github.com/trustwallet/blockatlas/pkg/address"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/golibs/coin"
)

func (c *Client) GetTransactions(address string, coinIndex uint) (blockatlas.TxPage, error) {
	page, err := c.GetTxs(address)
	if err != nil {
		return nil, err
	}
	return NormalizePage(page, address, "", coinIndex), nil
}

func (c *Client) GetTokenTxs(address, token string, coinIndex uint) (blockatlas.TxPage, error) {
	page, err := c.GetTxsWithContract(address, token)
	if err != nil {
		return nil, err
	}
	return NormalizePage(page, address, token, coinIndex), nil
}

func NormalizePage(srcPage *Page, address, token string, coinIndex uint) blockatlas.TxPage {
	var txs []blockatlas.Tx
	normalizedAddr := Address.EIP55Checksum(address)
	var normalizedToken string
	if token != "" {
		normalizedToken = Address.EIP55Checksum(token)
	}
	for _, srcTx := range srcPage.Transactions {
		tx := normalizeTxWithAddress(&srcTx, normalizedAddr, normalizedToken, coinIndex)
		txs = append(txs, tx)
	}
	return txs
}

func normalizeTx(srcTx *Transaction, coinIndex uint) blockatlas.Tx {
	blockTime, err := strconv.ParseInt(srcTx.TimeStamp, 10, 64)
	if err != nil {
		return blockatlas.Tx{}
	}

	blockHeight, err := strconv.ParseInt(srcTx.BlockNumber, 10, 64)
	if err != nil {
		return blockatlas.Tx{}
	}

	nonce, err := strconv.ParseUint(srcTx.Nonce, 10, 64)
	if err != nil {
		return blockatlas.Tx{}
	}

	status, errReason := srcTx.GetStatus()
	normalized := blockatlas.Tx{
		ID:       srcTx.TxID,
		Coin:     coinIndex,
		From:     srcTx.FromAddress(),
		To:       srcTx.ToAddress(),
		Fee:      blockatlas.Amount(srcTx.GetFee()),
		Date:     blockTime,
		Block:    normalizeBlockHeight(blockHeight),
		Status:   status,
		Error:    errReason,
		Sequence: nonce,
	}
	fillMeta(&normalized, srcTx, coinIndex)
	return normalized
}

func normalizeTxWithAddress(srcTx *Transaction, address, token string, coinIndex uint) blockatlas.Tx {
	normalized := normalizeTx(srcTx, coinIndex)
	normalized.Direction = GetDirection(address, normalized.From, normalized.To)
	fillMetaWithAddress(&normalized, srcTx, address, token, coinIndex)
	return normalized
}

func normalizeBlockHeight(height int64) uint64 {
	if height < 0 {
		return uint64(0)
	}
	return uint64(height)
}

func fillMeta(final *blockatlas.Tx, tx *Transaction, coinIndex uint) {
	if ok := fillTokenTransfer(final, tx, coinIndex); !ok {
		fillTransferOrContract(final, tx, coinIndex)
	}
}

func fillMetaWithAddress(final *blockatlas.Tx, tx *Transaction, address, token string, coinIndex uint) {
	if ok := fillTokenTransferWithAddress(final, tx, address, token, coinIndex); !ok {
		fillTransferOrContract(final, tx, coinIndex)
	}
}

func fillTokenTransfer(final *blockatlas.Tx, tx *Transaction, coinIndex uint) bool {

	decimal, err := strconv.Atoi(tx.TokenDecimal)
	if err != nil {
		return false
	}

	final.Meta = blockatlas.TokenTransfer{
		Name:     tx.TokenName,
		Symbol:   tx.TokenSymbol,
		TokenID:  tx.TokenName, // not specified on etherscan
		Decimals: uint(decimal),
		Value:    blockatlas.Amount(tx.Value),
		From:     tx.From,
		To:       tx.To,
	}

	return true
}

func fillTokenTransferWithAddress(final *blockatlas.Tx, tx *Transaction, address, token string, coinIndex uint) bool {
	decimal, err := strconv.Atoi(tx.TokenDecimal)
	if err != nil {
		return false
	}

	if tx.To == address || tx.From == address {
		// filter token if specified
		if token != "" {
			if token != tx.TokenName {
				return false
			}
		}
		direction := GetDirection(address, tx.From, tx.To)
		metadata := blockatlas.TokenTransfer{
			Name:     tx.TokenName,
			Symbol:   tx.TokenSymbol,
			TokenID:  tx.TokenName,
			Decimals: uint(decimal),
			Value:    blockatlas.Amount(tx.Value),
		}
		if direction == blockatlas.DirectionSelf {
			metadata.From = address
			metadata.To = address
		} else if direction == blockatlas.DirectionOutgoing {
			metadata.From = address
			metadata.To = tx.To
		} else {
			metadata.From = tx.From
			metadata.To = address
		}
		final.Direction = direction
		final.Meta = metadata
	}
	return true
}

func fillTransferOrContract(final *blockatlas.Tx, tx *Transaction, coinIndex uint) {
	gasUsed, err := strconv.Atoi(tx.GasUsed)
	if err != nil {
		return
	}

	if gasUsed == 21000 {
		final.Meta = blockatlas.Transfer{
			Value:    blockatlas.Amount(tx.Value),
			Symbol:   coin.Coins[coinIndex].Symbol,
			Decimals: coin.Coins[coinIndex].Decimals,
		}
		return
	}

	final.Meta = blockatlas.ContractCall{
		Input: "0x",
		Value: tx.Value,
	}
}
