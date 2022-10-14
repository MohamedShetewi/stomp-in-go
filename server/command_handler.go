package server

import (
	"github.com/MohamedShetewi/stomp-in-go/frame"
	"github.com/MohamedShetewi/stomp-in-go/utils"
)

type commandHandler func(server *Server, conn *Client, frame *frame.Frame)

var commandHandlerMap = map[frame.Command]commandHandler{
	frame.CONNECT: connectHandler,
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