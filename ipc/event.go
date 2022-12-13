package ipc

import (
	"bytes"
	"encoding/json"
	"strings"
)

type OperationType int

const (
	IpcDispatch OperationType = iota
	IpcReady
	IpcExit
)

type Message struct {
	Op OperationType `json:"op"`
	T  string        `json:"t,omitempty"`
	D  any           `json:"d,omitempty"`
}

func MarshalJSON(v interface{}) (string, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(&v); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}
