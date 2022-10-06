package main

type ErrorCode uint

const (
	OK ErrorCode = iota
	EmptyMessage
	UnsupportedCommand
	UnsupportedHeaderFormat
)

type Header struct {
	key   string
	value string
}

func (header Header) toString() string {
	return header.key + ":" + header.value
}

type Body string

type Frame struct {
	command Command
	headers map[string]string
	body    Body
}

func (frame *Frame) toUTF8() []byte {
	command := frame.command.encode()
	headers := ""
	for key, value := range frame.headers {
		headers += key + ":" + value + "\n"
	}
	body := string(frame.body)

	frameStream := []byte(command + "\n" + headers + "\n" + body)
	frameStream = append(frameStream, 0)

	return []byte(command + "\n" + headers + "\n" + body)
}
