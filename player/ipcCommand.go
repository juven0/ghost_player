package player

import (
	"bufio"
	"encoding/json"
	"fmt"

	"github.com/Microsoft/go-winio"
)

func (p *Player) sendCommand(command []string) (map[string]interface{}, error) {
	if p.pipe == "" {
		return nil, fmt.Errorf("pipe not initialized")
	}

	conn, err := winio.DialPipe(p.pipe, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to pipe: %w", err)
	}
	defer conn.Close()

	request := map[string]interface{}{
		"command": command,
	}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failde to mashal command: %w", err)
	}

	_, err = conn.Write(append(requestJSON, '\n'))
	if err != nil {
		return nil, fmt.Errorf("failed to write command: %w", err)
	}

	reader := bufio.NewReader(conn)
	responseJSON, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(responseJSON, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if errMsg, ok := response["error"].(string); ok && errMsg != "success" {
		return nil, fmt.Errorf("mpv error: %s", errMsg)
	}

	return response, nil
}

func (p *Player) GetPercentPos() (int, error) {
	response, err := p.sendCommand([]string{"get_property", "percent-pos"})
	if err != nil {
		return 0, err
	}

	percentPos, ok := response["data"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid percent-pos response")
	}

	return int(percentPos), nil
}
