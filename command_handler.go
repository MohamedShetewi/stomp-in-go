package main

import (
	"math"
	"strconv"
	"strings"
)

type CommandHandler func(server *Server, conn *Client, frame *Frame)

var CommandHandlerMap = map[Command]CommandHandler{
	CONNECT: connectHandler,
}

func connectHandler(server *Server, client *Client, frame *Frame) {
	if ok, _ := server.hasClient(client); !ok {
		//TODO send error message
		return
	}
	hearBeatSettings := strings.Split(frame.headers["heart-beat"], ",")
	if len(hearBeatSettings) != 2 {
		//TODO send error message
		return
	}

	clientOutHB := hearBeatSettings[0]
	if heartbeat, err := strconv.Atoi(clientOutHB); err == nil {
		if heartbeat == 0 {
			(*client).outHB = -1
		} else {
			(*client).outHB = int64(math.Max(float64(server.config.defaultHB), float64(heartbeat)))
		}
	} else {
		//TODO send error message
		return
	}
	clientInHB := hearBeatSettings[1]
	if heartbeat, err := strconv.Atoi(clientInHB); err == nil {
		if heartbeat == 0 {
			(*client).inHB = -1
		} else {
			(*client).inHB = int64(int(math.Max(float64(server.config.defaultHB), float64(heartbeat))))
		}
	} else {
		//TODO send error message
		return
	}
}
