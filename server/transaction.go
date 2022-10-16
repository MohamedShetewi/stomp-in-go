package server

import "bytes"

type transactionStatus int

const (
	PENDING = iota
	COMMITTED
	ABORTED
)

type transaction struct {
	id          string
	status      transactionStatus
	content     bytes.Buffer
	destination string
}

func findTX(txID string, transactions []*transaction) (bool, *transaction) {
	for _, tx := range transactions {
		if tx.id == txID {
			return true, tx
		}
	}
	return false, nil
}
