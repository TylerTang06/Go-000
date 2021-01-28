package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	listen, err := net.Listen("tcp", ":18000")
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	}
	defer listen.Close()

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	for sig := range sigs {
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			fmt.Println("Program Exit...", sig)
			cancel()
			return
		default:
			fmt.Println("other signal", sig)
			return
		}
	}

	for {
		fmt.Println("here")
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("accept err: %v\n", err)
			continue
		}

		go handleConn(ctx, conn)
	}
}

func handleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	msgCh := make(chan string, 10)
	go sendMsg(ctx, conn, msgCh)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		msgCh <- input.Text()
	}

	if err := input.Err(); err != nil {
		log.Printf("read msg meet err: %v\n", err)
	}
}

func sendMsg(ctx context.Context, conn net.Conn, ch <-chan string) {
	wr := bufio.NewWriter(conn)
	count := 0
	for msg := range ch {
		wr.WriteString(strconv.Itoa(count) + ":" + msg)
		wr.Flush()
	}
}
