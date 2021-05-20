package blockbook

import (
	"strconv"

	"github.com/trustwallet/blockatlas/pkg/blockatlas"
)

func (s *Transaction) GetStatus() (blockatlas.Status, string) {

	confirmations, err := strconv.Atoi(s.Confirmations)
	if err != nil {
		return blockatlas.StatusError, ""
	}

	if confirmations > 0 {
		return blockatlas.StatusCompleted, ""
	}

	return blockatlas.StatusPending, ""
}

func (t *Transaction) FromAddress() string {
	return t.From
}

func (t *Transaction) GetFee() string {
	gasPrice, err := strconv.Atoi(t.GasPrice)
	if err != nil {
		return "0"
	}

	gasUsed, err := strconv.Atoi(t.GasUsed)
	if err != nil {
		return "0"
	}

	fees := gasUsed * gasPrice

	status, _ := t.GetStatus()
	if status != blockatlas.StatusPending {
		return strconv.Itoa(fees)
	}

	gasLimit, err := strconv.Atoi(t.Gas)
	if err != nil {
		return "0"
	}

	fees = gasPrice * gasLimit

	return strconv.Itoa(fees)
}

func (t *Transaction) ToAddress() string {
	return t.To
}

func GetDirection(address, from, to string) blockatlas.Direction {
	if address == from && address == to {
		return blockatlas.DirectionSelf
	}
	if address == from {
		return blockatlas.DirectionOutgoing
	}
	return blockatlas.DirectionIncoming
}
