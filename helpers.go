package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/gliderlabs/ssh"
)

func checkPassword(ctx ssh.Context, pass string) bool {
	return redisConn.Get(fmt.Sprintf("%s:%s:pass", redisDocksshPrefix, ctx.User())).Val() == pass
}

func checkPublicKey(ctx ssh.Context, key ssh.PublicKey) bool {
	data := redisConn.Get(fmt.Sprintf("%s:%s:key", redisDocksshPrefix, ctx.User())).Val()
	allowed, _, _, _, _ := ssh.ParseAuthorizedKey([]byte(data))
	return ssh.KeysEqual(key, allowed)
}

// this function uses low-level system call to
// resize the PTY device "which is just a FD in unix systems".
// see: https://github.com/gliderlabs/ssh/blob/master/_examples/ssh-pty/pty.go
func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}
