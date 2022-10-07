package frame

import "bytes"

type ErrorCode uint

const (
	OK ErrorCode = iota
	EmptyMessage
	UnsupportedCommand
	UnsupportedHeaderFormat
)

type Body string

type Frame struct {
	Command Command
	Headers map[string]string
	body    Body
}

func (frame *Frame) ToUTF8() []byte {
	var buffer bytes.Buffer

	cmd := frame.Command.Encode()
	buffer.Write([]byte(cmd))

	for key, value := range frame.Headers {
		header := key + ":" + value + "\n"
		buffer.Write([]byte(header))
	}
	buffer.Write([]byte(frame.body + "\n"))

	return buffer.Bytes()
}
