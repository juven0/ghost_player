package mpv

import (
	"context"
	"encoding/json"
	"fmt"

	"player/internal/ports"
)

type MpvPlayer struct {
	client  *Client
	process *Process
}

func NewMpvPlayer(client *Client, process *Process) ports.Player {
	return &MpvPlayer{
		client:  client,
		process: process,
	}
}

func (m *MpvPlayer) Connect(ctx context.Context) error {
	return m.client.pipe.Connect(ctx, m.process.PipePath())
}

func (m *MpvPlayer) Play() error {
	_, err := m.client.SendCommand([]string{"set_property", "pause", "false"})
	return err
}

func (m *MpvPlayer) Pause() error {
	_, err := m.client.SendCommand([]string{"set_proterty", "pause", "true"})
	return err
}

func (m *MpvPlayer) Stop() error {
	_, err := m.client.SendCommand([]string{"stop"})
	return err
}

func (m *MpvPlayer) Resume() error {
	_, err := m.client.SendCommand([]string{"set_property", "pause", "false"})
	return err
}

func (m *MpvPlayer) SetVolume(value int) error {
	_, err := m.client.SendCommand([]string{"set_property", "volume", fmt.Sprintf("%d", value)})
	return err
}

func (m *MpvPlayer) GetVolume() (int, error) {
	resp, err := m.client.SendCommand([]string{"get_property", "volume"})
	if err != nil {
		return 0, err
	}
	var vol float64
	if err := json.Unmarshal(resp.Data, &vol); err != nil {
		return 0, fmt.Errorf("invalid volume response")
	}
	return int(vol), nil
}

func (m *MpvPlayer) Seek(seconds float64, mode ports.SeekMode) error {
	_, err := m.client.SendCommand([]string{"seek", fmt.Sprintf("%f", seconds), string(mode)})
	return err
}

func (m *MpvPlayer) SetMute(mute bool) error {
	_, err := m.client.SendCommand([]string{"set_property", "mute", fmt.Sprintf("%s", mute)})
	return err
}

func (m *MpvPlayer) GetMute() (bool, error) {
	resp, err := m.client.SendCommand([]string{"get_property", "mute"})
	if err != nil {
		return false, err
	}
	var mute bool
	err = json.Unmarshal(resp.Data, &mute)
	if err != nil {
		return false, fmt.Errorf("invalide mute response")
	}
	return mute, nil
}

func (m *MpvPlayer) StartPlay(ctx context.Context, streamUrl string) error {
	err := m.process.Start(streamUrl)
	if err != nil {
		return err
	}
	return m.Connect(ctx)
}

func (m *MpvPlayer) GetPercentPos() (int, error) {
	resp, err := m.client.SendCommand([]string{"get_property", "percent-pos"})
	if err != nil {
		return 0, err
	}
	var pos float64
	if err := json.Unmarshal(resp.Data, &pos); err != nil {
		return 0, fmt.Errorf("invalid percent-pos response")
	}

	return int(pos), nil
}

func (m *MpvPlayer) Close() error {
	m.client.SendCommand([]string{"quit"})
	m.process.Kill()
	return m.client.Close()
}
