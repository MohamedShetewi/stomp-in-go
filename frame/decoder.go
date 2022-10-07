package frame

import (
	"strings"
)

func Decode(msg string) (*Frame, ErrorCode) {
	msgSplit := strings.Split(msg, "\n")

	if len(msgSplit) < 1 {
		return nil, EmptyMessage
	}

	var cmd Command
	if val, err := SupportedCommands(msgSplit[0]); err != nil {
		cmd = val
	} else {
		return nil, UnsupportedCommand
	}

	headers, newOffset, errCode := decodeHeaders(msgSplit, 1)
	if errCode != OK {
		return nil, errCode
	}

	body := decodeBody(msgSplit, newOffset)
	return &Frame{
		Command: cmd,
		Headers: headers,
		Body:    Body(body),
	}, OK
}

func decodeHeaders(frame []string, offset int) (map[string]string, int, ErrorCode) {
	headers := make(map[string]string)
	for ; offset < len(frame) && frame[offset] != ""; offset += 1 {
		header := strings.Split(frame[offset], ":")
		if len(header) != 2 {
			return nil, -1, UnsupportedHeaderFormat
		}
		headers[header[0]] = header[1]
	}
	return headers, offset, OK
}

func decodeBody(frame []string, offset int) string {
	var body string
	for i := offset; i < len(frame) && frame[i] != ""; i += 1 {
		body += frame[i]
	}
	return body
}
