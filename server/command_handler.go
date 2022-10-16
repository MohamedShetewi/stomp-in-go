package server

import (
	"github.com/MohamedShetewi/stomp-in-go/frame"
	"github.com/MohamedShetewi/stomp-in-go/utils"
)

type commandHandler func(server *Server, conn *Client, frame *frame.Frame)

var commandHandlerMap = map[frame.Command]commandHandler{
	frame.CONNECT:     connectHandler,
	frame.SUBSCRIBE:   subscribeHandler,
	frame.UNSUBSCRIBE: unsubscribeHandler,
	frame.SEND:        sendHandler,
	frame.BEGIN:       beginHandler,
}

func connectHandler(server *Server, client *Client, frm *frame.Frame) {
	if ok, _ := server.HasClient(client); !ok {
		sendError(client, frm.Headers,
			"Connection Error: connection is already established!")
		return
	}
	clientOutHB, clientInHB, err := frame.HeartBeatParser(frm.Headers["heart-beat"])

	if err != nil {
		sendError(client, frm.Headers, err.Error())
		return
	}

	if clientOutHB == 0 {
		client.outHB = -1
	} else {
		client.outHB = utils.Max(server.config.defaultHB, clientOutHB)
	}

	if clientInHB == 0 {
		client.inHB = -1
	} else {
		client.inHB = utils.Max(server.config.defaultHB, clientInHB)
	}
}

func subscribeHandler(server *Server, client *Client, frm *frame.Frame) {
	destination, ok := frm.Headers["destination"]
	if !ok {
		sendError(client, nil, "destination field is required in the subscribe frame")
		return
	}
	if _, ok := server.destinations[destination]; !ok {
		sendError(client, nil, destination+"is not supported in the server")
		return
	}
	newSubscriber := &subscriber{
		client:       client,
		isSubscribed: true,
	}
	server.destinationsLock.Lock()
	subList := server.destinations[destination]
	newSubList := append(subList, newSubscriber)
	server.destinations[destination] = newSubList
	server.destinationsLock.Unlock()
}

func unsubscribeHandler(server *Server, client *Client, frm *frame.Frame) {
	destination, ok := frm.Headers["destination"]
	if !ok {
		sendError(client, nil, "destination field is required in the subscribe frame")
		return
	}
	if _, ok := server.destinations[destination]; !ok {
		sendError(client, nil, destination+"is not supported in the server")
		return
	}
	subscribers := server.destinations[destination]
	for _, sub := range subscribers {
		if sub.client == client {
			sub.isSubscribed = false
		}
	}
	// TODO remove unsubscribed clients after a threshold
}

func sendHandler(server *Server, client *Client, frm *frame.Frame) {
	destination, ok := frm.Headers["destination"]
	if !ok {
		sendError(client, make(map[string]string), "destination field is required in the subscribe frame")
		return
	}
	if _, ok := server.destinations[destination]; !ok {
		sendError(client, make(map[string]string), destination+"is not supported in the server")
		return
	}
	headers := map[string]string{
		"destination": destination,
	}
	msg := &frame.Frame{
		Command: frame.MESSAGE,
		Headers: headers,
		Body:    frm.Body,
	}
	for _, subscriber := range server.destinations[destination] {
		if subscriber.isSubscribed {
			subscriber.client.sendChan <- msg.ToUTF8()
		}
	}
}

func beginHandler(server *Server, client *Client, frm *frame.Frame) {
	txID, ok := frm.Headers["transaction"]
	if !ok {
		sendError(client, make(map[string]string), "cannot find transaction id in the begin frame")
	}
	for _, transaction := range client.transactions {
		if transaction.id == txID {
			sendError(client, make(map[string]string), "transaction with "+txID+" already exists")
			return
		}
	}
	newTx := &transaction{
		id:      txID,
		status:  PENDING,
		content: "",
	}
	client.transactions = append(client.transactions, newTx)
}
