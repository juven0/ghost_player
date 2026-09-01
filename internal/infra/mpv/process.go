package mpv

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type Process struct {
	cmd    *exec.Cmd
	ctx    context.Context
	cancel context.CancelFunc
	pipe   string
}

func NewProcess(ctx context.Context) *Process {
	ctx, cancel := context.WithCancel(ctx)
	return &Process{
		ctx:    ctx,
		cancel: cancel,
		pipe:   generatePipe(),
	}
}

func (p Process) PipePath() string {
	return p.pipe
}

func (p *Process) Start(streamURL string) error {
	p.cmd = exec.CommandContext(p.ctx, "mpv",
		streamURL,
		"--no-video",
		"--ytdl-format=bestaudio",
		fmt.Sprintf("--input-ipc-server=%s", p.pipe),
		"--idle=yes",
		"--keep-open=yes",
	)

	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start mpv: %w", err)
	}

	return nil
}

func (p *Process) Wait() error {
	return p.cmd.Wait()
}

func (p *Process) Kill() error {
	p.cancel()
	if p.cmd.Process != nil {
		return p.cmd.Process.Kill()
	}

	return nil
}

func generatePipe() string {
	pid := os.Getpid()
	if runtime.GOOS == "windows" {
		return fmt.Sprintf(`\\.\pipe\mpvsocket_%d`, pid)
	}

	return fmt.Sprintf("/tmp/mpvsocket_%d", pid)
}
