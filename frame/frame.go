package frame

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
	cmd := frame.Command.Encode()
	headers := ""
	for key, value := range frame.Headers {
		headers += key + ":" + value + "\n"
	}
	body := string(frame.body)

	frameStream := []byte(cmd + "\n" + headers + "\n" + body)
	frameStream = append(frameStream, 0)

	return []byte(cmd + "\n" + headers + "\n" + body)
}
