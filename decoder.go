package main

import "strings"

func Decode(msg string) (*Frame, ErrorCode) {
	msgSplit := strings.Split(msg, "\n")

	if len(msgSplit) < 1 {
		return nil, EmptyMessage
	}

	var command Command
	if val, ok := SupportedCommands[msgSplit[0]]; ok {
		command = val
	} else {
		return nil, UnsupportedCommand
	}

	headers, newOffset, errCode := decodeHeaders(msgSplit, 1)
	if errCode != OK {
		return nil, errCode
	}

	body := decodeBody(msgSplit, newOffset)
	return &Frame{
		command: command,
		headers: headers,
		body:    Body(body),
	}, OK
}

func decodeHeaders(frame []string, offset int) (headers map[string]string, newOffset int, code ErrorCode) {
	var i = offset
	for ; i < len(frame) && frame[i] != ""; i += 1 {
		header := strings.Split(frame[i], ":")
		if len(header) != 2 {
			return nil, -1, UnsupportedHeaderFormat
		}
		headers[header[0]] = header[1]
	}
	return headers, i, OK
}

func decodeBody(frame []string, offset int) string {
	var body string
	for i := offset; i < len(frame) && frame[i] != ""; i += 1 {
		body += frame[i]
	}
	return body
}
