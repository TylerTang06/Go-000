package week03

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出

// Signal ...
func Signal() {
	// Go signal notification works by sending `os.Signal`
	// values on a channel. We'll create a channel to
	// receive these notifications (we'll also make one to
	// notify us when the program can exit).
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	// This goroutine executes a blocking receive for
	// signals. When it gets one it'll print it out
	// and then notify the program that it can finish.
	go func() {
		for sig := range sigs {
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("Program Exit...", sig)
			case syscall.SIGUSR1:
				fmt.Println("usr1 signal", sig)
			case syscall.SIGUSR2:
				fmt.Println("usr2 signal", sig)
			default:
				fmt.Println("other signal", sig)
			}
		}

		done <- true
	}()

	// The program will wait here until it gets the
	// expected signal (as indicated by the goroutine
	// above sending a value on `done`) and then exit.
	fmt.Println("awaiting signal")
	<-done
	fmt.Println("exiting")
}
