package server

type transactionStatus int

const (
	PENDING = iota
	COMMITTED
	ABORTED
)

type transaction struct {
	id      string
	status  transactionStatus
	content string
}
