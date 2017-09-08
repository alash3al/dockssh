// Dockssh is a tiny tool that exposes your docker container to the world
// behind a ssh server, so you can access any of them using your ssh client.
// It uses redis as password storage for the containers.
// Each container has a name, and each name has a pass,
// it uses this format for redis keys to store passwords (dockssh:$CONTAINER_NAME:pass).
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"
	"unsafe"
)

import (
	"github.com/gliderlabs/ssh"
	"github.com/go-redis/redis"
	"github.com/kr/pty"
)

var (
	LISTEN_ADDR     = flag.String("listen-addr", "0.0.0.0:22022", "The ssh listening address")
	HOST_KEY_FILE   = flag.String("host-key", path.Join(os.Getenv("HOME"), ".ssh/id_rsa"), "The host key, if you left it empty, it would be generated everytime")
	DEFAULT_COMMAND = flag.String("default-command", "/bin/sh", "The default command to execute")
	REDIS_ADDR      = flag.String("redis-addr", "localhost:6379", "The redis host address")
	REDIS_PASSWORD  = flag.String("redis-password", "", "The redis password")
	REDIS_DB        = flag.Int("redis-db", 0, `The redis database number (default "0")`)
)

var (
	RedisConn *redis.Client
)

func main() {
	flag.Parse()
	RedisConn = redis.NewClient(&redis.Options{
		Addr:     *REDIS_ADDR,
		Password: *REDIS_PASSWORD,
		DB:       *REDIS_DB,
	})
	if _, err := RedisConn.Ping().Result(); err != nil {
		log.Fatal(err)
	}
	ssh.Handle(handler)
	log.Fatal(ssh.ListenAndServe(
		*LISTEN_ADDR,
		nil,
		ssh.HostKeyFile(*HOST_KEY_FILE),
		ssh.PasswordAuth(checkPassword),
	))
}

func handler(session ssh.Session) {
	io.WriteString(session, "\n")
	io.WriteString(session, "\tWelcome to Dockssh ^^\n")
	io.WriteString(session, fmt.Sprintf("\tDefault Dockssh Command is '%s'\n", *DEFAULT_COMMAND))
	io.WriteString(session, fmt.Sprintf("\tYour current container is '%s'", session.User()))
	io.WriteString(session, "\n\n")

	// check whether the current session supports PTY or not.
	_, winCh, isPty := session.Pty()
	if !isPty {
		io.WriteString(session, "Your SSH client didn't request a PTY\n")
		return
	}

	// initialize the default command,
	// and append the session environment vars to it.
	cmd := exec.Command("docker", "exec", "-it", session.User(), *DEFAULT_COMMAND)
	cmd.Env = append(cmd.Env, session.Environ()...)

	// starts the previous cmd in a PTY device.
	terminal, err := pty.Start(cmd)
	if err != nil {
		io.WriteString(session, "Invalid params specified")
		return
	}

	// a goroutine that handles window size changes.
	go func() {
		for win := range winCh {
			setWinsize(terminal, win.Width, win.Height)
		}
	}()

	// a goroutine that reads from
	// the session and write to the PTY device.
	go func() {
		io.Copy(terminal, session)
	}()

	// wait for data from the PTY device and pipe it to the session.
	// **a blocking call** so that the current session keeps running.
	io.Copy(session, terminal)
}

// this is our ssh.PasswordHandler
func checkPassword(ctx ssh.Context, pass string) bool {
	return RedisConn.Get(fmt.Sprintf("dockssh:%s:pass", ctx.User())).Val() == pass
}

// this function uses low-level system call to
// resize the PTY device "which is just a FD in unix systems".
// see: https://github.com/gliderlabs/ssh/blob/master/_examples/ssh-pty/pty.go
func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}
