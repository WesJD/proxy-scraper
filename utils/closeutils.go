package utils

import (
	"os"
	"os/signal"
	"syscall"
)

func WatchForKill() (signals chan os.Signal) {
	signals = make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	return
}
