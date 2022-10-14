package frame

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type Body string

type Frame struct {
	Command Command
	Headers map[string]string
	Body    Body
}

func (frame *Frame) ToUTF8() []byte {
	var buffer bytes.Buffer

	cmd := frame.Command.Encode()
	buffer.Write([]byte(cmd))

	for key, value := range frame.Headers {
		header := key + ":" + value + "\n"
		buffer.Write([]byte(header))
	}
	buffer.Write([]byte(frame.Body + "\n"))

	return buffer.Bytes()
}

func HeartBeatParser(heartBeatHeader string) (int64, int64, error) {
	heartBeatSplitted := strings.Split(heartBeatHeader, ",")

	if len(heartBeatSplitted) != 2 {
		return -1, -1, errors.New("Invalid HeartBeat Format")
	}

	outHB, err := strconv.Atoi(heartBeatSplitted[0])
	if err != nil {
		return -1, -1, errors.New("Invalid type: heartbeat settings must be in a valid format")
	}

	inHB, err := strconv.Atoi(heartBeatSplitted[1])
	if err != nil {
		return -1, -1, errors.New("Invalid type: heartbeat settings must be in a valid format")
	}

	return int64(outHB), int64(inHB), nil
}
