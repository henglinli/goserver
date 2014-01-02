package utils

import (
	"os"
	"os/signal"
	"syscall"
	"log"
)

func Wait() {
	exit_chan := make(chan int)
	signal_chan := make(chan os.Signal, 1)
	go func() {
		<-signal_chan
		log.Println("Caught signal, exiting...")
		exit_chan <- 1
	}()
	signal.Notify(signal_chan, syscall.SIGINT, syscall.SIGTERM)
	<- exit_chan	
}
