package service

import (
	"os"
	"syscall"
)

var signalWindows = []os.Signal{
	os.Interrupt,
	os.Kill,

	syscall.SIGKILL,
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGQUIT,
	syscall.SIGILL,
	syscall.SIGTRAP,
	syscall.SIGABRT,
	syscall.SIGTERM,
}
