package frame

import (
	"errors"
	"strings"
)

func Decode(msg string) (*Frame, error) {
	msgSplit := strings.Split(msg, "\n")

	if len(msgSplit) < 1 {
		return nil, errors.New("empty frame")
	}

	var cmd Command
	if val, err := SupportedCommands(msgSplit[0]); err != nil {
		cmd = val
	} else {
		return nil, errors.New("unsupported command")
	}

	headers, newOffset, err := decodeHeaders(msgSplit, 1)
	if err != nil {
		return nil, err
	}

	body := decodeBody(msgSplit, newOffset)
	return &Frame{
		Command: cmd,
		Headers: headers,
		Body:    Body(body),
	}, nil
}

func decodeHeaders(frame []string, offset int) (map[string]string, int, error) {
	headers := make(map[string]string)
	for ; offset < len(frame) && frame[offset] != ""; offset += 1 {
		header := strings.Split(frame[offset], ":")
		if len(header) != 2 {
			return nil, -1, errors.New("unsupported header format")
		}
		headers[header[0]] = header[1]
	}
	return headers, offset, nil
}

func decodeBody(frame []string, offset int) string {
	var body string
	for i := offset; i < len(frame) && frame[i] != ""; i += 1 {
		body += frame[i]
	}
	return body
}
