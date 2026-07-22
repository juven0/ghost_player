package mpv

import (
	"encoding/json"
	"fmt"

	"player/internal/infra/ipc"
)

type Client struct {
	pipe ipc.IPCTransport
}

func NewClient(pipe ipc.IPCTransport) *Client {
	return &Client{pipe: pipe}
}

func (p *Client) SendCommand(command []string) (*ipc.Response, error) {
	req := ipc.Request{Command: command}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failde to marshal request: %w", err)
	}

	rawData, err := p.pipe.Send(data)
	if err != nil {
		return nil, err
	}

	var resp ipc.Response
	if err := json.Unmarshal(rawData, &resp); err != nil {
		return nil, fmt.Errorf("faild to parse response: %w", err)
	}

	if resp.Error != "success" {
		return nil, fmt.Errorf("mpv error: %s ", resp.Error)
	}

	return &resp, nil
}

func (p *Client) Close() error {
	return p.pipe.Close()
}
