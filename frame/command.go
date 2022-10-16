package frame

import "errors"

type Command interface {
	Encode() string
}

type ClientCommand string
type ServerCommand string

const (
	SEND        ClientCommand = "SEND"
	SUBSCRIBE   ClientCommand = "SUBSCRIBE"
	UNSUBSCRIBE ClientCommand = "UNSUBSCRIBE"
	CONNECT     ClientCommand = "CONNECT"
	DISCONNECT  ClientCommand = "DISCONNECT"
	STOMP       ClientCommand = "STOMP"
	BEGIN       ClientCommand = "BEGIN"
	COMMIT      ClientCommand = "COMMIT"
	ABORT       ClientCommand = "ABORT"
)

func (command ClientCommand) Encode() string {
	return string(command)
}

const (
	CONNECTED ServerCommand = "CONNECTED"
	MESSAGE   ServerCommand = "MESSAGE"
	RECEIPT   ServerCommand = "RECEIPT"
	ERROR     ServerCommand = "ERROR"
)

func (command ServerCommand) Encode() string {
	return string(command)
}

func SupportedCommands(cmd string) (Command, error) {
	if val, ok := supportedCommands[cmd]; ok {
		return val, nil
	}
	return ERROR, errors.New("unsupported command")
}

var supportedCommands = map[string]Command{
	"SEND":        SEND,
	"SUBSCRIBE":   SUBSCRIBE,
	"UNSUBSCRIBE": UNSUBSCRIBE,
	"CONNECT":     CONNECT,
	"DISCONNECT":  DISCONNECT,
	"STOMP":       STOMP,
	"CONNECTED":   CONNECTED,
	"MESSAGE":     MESSAGE,
	"RECEIPT":     RECEIPT,
	"ERROR":       ERROR,
	"Begin":       BEGIN,
	"COMMIT":      COMMIT,
	"ABORT":       ABORT,
}
