package ipc

import (
	"bytes"
	"encoding/json"
	"strings"
)

type Message struct {
	Event string `json:"e,omitempty"`
	Data  any    `json:"d,omitempty"`
}

func MarshalJSON(v any) (string, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(&v); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}
