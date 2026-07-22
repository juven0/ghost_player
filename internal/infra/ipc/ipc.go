package ipc

import (
	"context"
	"encoding/json"
)

type IPCTransport interface {
	Connect(ctx context.Context, pipePath string) error
	Send(data []byte) ([]byte, error)
	Close() error
}

type Request struct {
	Command   []string `json:"command"`
	RequestId int      `json:"request_id, omitempty"`
}

type Response struct {
	Error     string          `json:"error"`
	Data      json.RawMessage `json:"data"`
	RequestId int             `json:"request_id, omitempty"`
}

type Event struct {
	Event string          `json:"event"`
	Name  string          `json:"name, omitempty"`
	Data  json.RawMessage `json:"data, omitempty"`
}
