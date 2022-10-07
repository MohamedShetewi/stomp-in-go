package main

type Command interface {
	encode() string
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
)

func (command ClientCommand) encode() string {
	return string(command)
}

const (
	CONNECTED ServerCommand = "CONNECTED"
	MESSAGE   ServerCommand = "MESSAGE"
	RECEIPT   ServerCommand = "RECEIPT"
	ERROR     ServerCommand = "ERROR"
)

func (command ServerCommand) encode() string {
	return string(command)
}

var SupportedCommands = map[string]Command{
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
}
