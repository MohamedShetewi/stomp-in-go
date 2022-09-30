package main

type Header struct {
	key   string
	value string
}

func (header Header) toString() string {
	return header.key + ":" + header.value
}

type Body string

type Frame[T Command] struct {
	command T
	headers []*Header
	body    Body
}

func (frame *Frame[T]) toUTF8() []byte {
	command := Command(frame.command).encode()
	headers := ""
	for _, header := range frame.headers {
		headers += header.toString() + "\n"
	}
	body := string(frame.body)

	frameStream := []byte(command + "\n" + headers + "\n" + body)
	frameStream = append(frameStream, 0)

	return []byte(command + "\n" + headers + "\n" + body)
}
