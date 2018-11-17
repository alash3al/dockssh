package main

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/gliderlabs/ssh"
	"github.com/kr/pty"
)

func handler(session ssh.Session) {
	io.WriteString(session, "\n")
	io.WriteString(session, "\tWelcome to Dockssh ^^\n")
	io.WriteString(session, fmt.Sprintf("\tDefault Dockssh Command is '%s'\n", *flagEntryPoint))
	io.WriteString(session, fmt.Sprintf("\tYour current container is '%s'", session.User()))
	io.WriteString(session, "\n\n")

	// check whether the current session supports PTY or not.
	ptyReq, winCh, isPty := session.Pty()
	if !isPty {
		io.WriteString(session, "Your SSH client didn't request a PTY\n")
		return
	}

	cmd := exec.Command("docker", "exec", "-it", session.User(), *flagEntryPoint)
	cmd.Env = append(cmd.Env, session.Environ()...)
	cmd.Env = append(cmd.Env, "TERM="+ptyReq.Term)

	terminal, err := pty.Start(cmd)
	if err != nil {
		io.WriteString(session, err.Error())
		return
	}

	go func() {
		for win := range winCh {
			setWinsize(terminal, win.Width, win.Height)
		}
	}()

	go func() {
		io.Copy(terminal, session)
	}()

	io.Copy(session, terminal)
}
