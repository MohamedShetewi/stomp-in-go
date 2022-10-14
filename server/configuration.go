package server

import "time"

type Configuration struct {
	topics             []string
	defaultHB          int64
	maxFrameSize       int
	maxTcpConnTime     time.Time
	deadlineForConnect time.Duration
}
