package server

import "time"

type Configuration struct {
	defaultHB          int64
	maxFrameSize       int
	maxTcpConnTime     time.Time
	deadlineForConnect time.Duration
}
