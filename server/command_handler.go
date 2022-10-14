package server

import (
	"github.com/MohamedShetewi/stomp-in-go/frame"
	"github.com/MohamedShetewi/stomp-in-go/utils"
)

type commandHandler func(server *Server, conn *Client, frame *frame.Frame)

var commandHandlerMap = map[frame.Command]commandHandler{
	frame.CONNECT:   connectHandler,
	frame.SUBSCRIBE: subscribeHandler,
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
	if _, ok := server.subscribers[destination]; !ok {
		sendError(client, nil, destination+"is not supported in the server")
		return
	}
	server.subscribersLock.Lock()
	subList := server.subscribers[destination]
	newSubList := append(subList, client)
	server.subscribers[destination] = newSubList
	server.subscribersLock.Unlock()
}
