package server

import (
	"github.com/MohamedShetewi/stomp-in-go/frame"
	"math"
	"strconv"
	"strings"
)

type commandHandler func(server *Server, conn *Client, frame *frame.Frame)

var commandHandlerMap = map[frame.Command]commandHandler{
	frame.CONNECT: connectHandler,
}

func connectHandler(server *Server, client *Client, frame *frame.Frame) {
	if ok, _ := server.HasClient(client); !ok {
		sendError(client, frame.Headers,
			"Connection Error: connection is already established!")
		return
	}
	hearBeatSettings := strings.Split(frame.Headers["heart-beat"], ",")
	if len(hearBeatSettings) != 2 {
		//TODO send error message
		return
	}

	clientOutHB := hearBeatSettings[0]
	if heartbeat, err := strconv.Atoi(clientOutHB); err == nil {
		if heartbeat == 0 {
			client.outHB = -1
		} else {
			client.outHB = int64(math.Max(float64(server.config.defaultHB), float64(heartbeat)))
		}
	} else {
		sendError(client, frame.Headers,
			"Invalid type: heartbeat settings must be in a valid format")
		return
	}
	clientInHB := hearBeatSettings[1]
	if heartbeat, err := strconv.Atoi(clientInHB); err == nil {
		if heartbeat == 0 {
			client.inHB = -1
		} else {
			client.inHB = int64(int(math.Max(float64(server.config.defaultHB), float64(heartbeat))))
		}
	} else {
		sendError(client, frame.Headers,
			"Invalid type: heartbeat settings must be in a valid format")
		return
	}
}
