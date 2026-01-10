package signal

import (
	"fmt"
	"strings"
	"syscall"
)

func Parse(s string) (syscall.Signal, error) {
	s = strings.ToUpper(s)
	s = strings.TrimPrefix(s, "SIG")

	switch s {
	case "KILL", "9":
		return syscall.SIGKILL, nil
	case "TERM", "15":
		return syscall.SIGTERM, nil
	case "INT", "2":
		return syscall.SIGINT, nil
	case "HUP", "1":
		return syscall.SIGHUP, nil
	case "QUIT", "3":
		return syscall.SIGQUIT, nil
	case "USR1", "10":
		return syscall.SIGUSR1, nil
	case "USR2", "12":
		return syscall.SIGUSR2, nil
	case "STOP", "19":
		return syscall.SIGSTOP, nil
	case "CONT", "18":
		return syscall.SIGCONT, nil
	default:
		return 0, fmt.Errorf("unknown signal: %s", s)
	}
}

func Kill(pid int, sig syscall.Signal) error {
	return syscall.Kill(pid, sig)
}

func Name(sig syscall.Signal) string {
	switch sig {
	case syscall.SIGKILL:
		return "SIGKILL"
	case syscall.SIGTERM:
		return "SIGTERM"
	case syscall.SIGINT:
		return "SIGINT"
	case syscall.SIGHUP:
		return "SIGHUP"
	case syscall.SIGQUIT:
		return "SIGQUIT"
	default:
		return "UNKNOWN"
	}
}
