package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	terminator := bytes.Index(data, []byte("\r\n"))
	if terminator == -1 {
		return 0, false, nil
	}

	line := string(data[:terminator])

	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		return 0, false, errors.New("invalid field-line")
	}

	if strings.HasSuffix(parts[0], " ") {
		return 0, false, errors.New("invalid field-line format: whitespace between field-name and colon")
	}

	name := strings.ToLower(strings.TrimSpace(parts[0]))
	value := strings.TrimSpace(parts[1])

	if len(name) < 1 {
		return 0, false, errors.New("invalid field-name: empty field-name")
	}

	for _, r := range name {
		if ('a' <= r && r <= 'z') ||
			('A' <= r && r <= 'Z') ||
			('0' <= r && r <= '9') {
			continue
		}

		switch r {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			continue
		}

		return 0, false, errors.New("invalid field-name characters")
	}

	if existing, ok := h[name]; !ok {
		h[name] = value
	} else {
		h[name] = fmt.Sprintf("%s, %s", existing, value)
	}

	return terminator + 2, true, nil
}
