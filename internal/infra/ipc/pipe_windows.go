// go:build windows

package ipc

import (
	"bufio"
	"context"
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"
)

type PipeWindows struct {
	conn net.Conn
}

func NewPipeWindows() IPCTransport {
	return &PipeWindows{}
}

func (p *PipeWindows) Connect(ctx context.Context, pipePath string) error {
	conn, err := winio.DialPipe(pipePath, nil)
	if err != nil {
		return fmt.Errorf("faild to connect to pipe: %w", err)
	}

	p.conn = conn

	return nil
}

func (p *PipeWindows) Send(data []byte) ([]byte, error) {
	_, err := p.conn.Write(append(data, '\n'))
	if err != nil {
		return nil, fmt.Errorf("faildeto write data: %w", err)
	}
	reader := bufio.NewReader(p.conn)
	response, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}

	return response, nil
}

func (p *PipeWindows) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}

	return nil
}
