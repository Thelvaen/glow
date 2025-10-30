package ui

import (
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/fyne-io/terminal"
)

type Session struct {
	Title string
	Cmd   *exec.Cmd
	Pty   *os.File
	Term  *terminal.Terminal
}

func NewSession(title, cmdStr, cwd string, env map[string]string) (*Session, error) {
	t := terminal.New()
	cmd := exec.Command("/bin/bash", "-lc", cmdStr)
	if cwd != "" {
		cmd.Dir = cwd
	}
	if len(env) > 0 {
		for k, v := range env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	// Wire IO
	go func() { defer ptmx.Close(); _, _ = io.Copy(ptmx, t) }()
	go func() { defer ptmx.Close(); _, _ = io.Copy(t, ptmx) }()

	// Resize sync
	resizeCh := make(chan terminal.Config, 1)
	t.AddListener(resizeCh)
	go func() {
		var rows, cols uint
		for cfg := range resizeCh {
			if cfg.Rows == rows && cfg.Columns == cols {
				continue
			}
			rows, cols = cfg.Rows, cfg.Columns
			_ = pty.Setsize(ptmx, &pty.Winsize{Rows: uint16(rows), Cols: uint16(cols)})
		}
	}()

	return &Session{Title: title, Cmd: cmd, Pty: ptmx, Term: t}, nil
}

func (s *Session) Close() {
	if s.Pty != nil {
		_ = s.Pty.Close()
	}
	if s.Cmd != nil && s.Cmd.Process != nil {
		_ = s.Cmd.Process.Kill()
	}
}
