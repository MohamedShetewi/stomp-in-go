package server

import (
	"github.com/MohamedShetewi/stomp-in-go/frame"
	"strconv"
)

func sendError(client *Client, headers map[string]string, msg string) {
	cmd := frame.ERROR
	headers["content-type"] = "text/plain"
	headers["content-length"] = strconv.Itoa(len(msg))

	frm := frame.Frame{
		Command: cmd,
		Headers: headers,
		Body:    frame.Body(msg),
	}
	client.sendChan <- frm.ToUTF8()
}

func sendMessage(client *Client, headers map[string]string, body []byte) {
	cmd := frame.MESSAGE

	frm := frame.Frame{
		Command: cmd,
		Headers: headers,
		Body:    frame.Body(body),
	}
	client.sendChan <- frm.ToUTF8()
}
