package utils

import (
	"syscall"
	"os"
	"os/signal"
	"fmt"
)

func Pause() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(interrupt)
	sig := <-interrupt
	fmt.Println(" program received exit signal:", sig)
}
