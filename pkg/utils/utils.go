package utils

import (
	"os"
	"os/signal"
	"syscall"
)

var signalInterruptChan chan os.Signal

// TriggerSignalInterrupt sends os interrupt signal to the listening channel
func TriggerSignalInterrupt() {
	if signalInterruptChan != nil {
		signalInterruptChan <- os.Interrupt
	}
}

// WaitForSignalInterrupt creates channel that listens for os signal interrupt event
func WaitForSignalInterrupt() {
	signalInterruptChan = make(chan os.Signal, 1)
	signal.Notify(signalInterruptChan, os.Interrupt, syscall.SIGTERM)
	<-signalInterruptChan
}
