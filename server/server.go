package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unicode/utf8"

	winio "github.com/Microsoft/go-winio"
)

func main() {
	pipeName := flag.String("p", "", "pipe name")
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	go pipeServer(ctx, *pipeName)

	// blocking main
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(c)
	<-c
	log.Println("shutdown")
	cancel()
	time.Sleep(2 * time.Second)
}

func pipeServer(ctx context.Context, pipeName string) {
	fullPipeName := fmt.Sprintf(`\\.\pipe\%s`, pipeName)
	log.Println(fullPipeName)
	listener, err := winio.ListenPipe(fullPipeName, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		log.Println("close listener")
		listener.Close()
	}()

	go pipeServerHandleListener(ctx, listener)
	for {
		// wait stop event
		select {
		case <-ctx.Done():
			log.Println("listener receive stop event")
			return
		default:
		}
	}
}

func pipeServerHandleListener(ctx context.Context, listener net.Listener) {
	for {
		// don't block waiting stop event
		conn, err := listener.Accept()
		if err != nil {
			if err == winio.ErrPipeListenerClosed {
				break
			} else {
				log.Println("listen with error:", err)
				continue
			}
		}
		go pipeServerHandleConnection(conn)
	}
}

func pipeServerHandleConnection(cn net.Conn) {
	eof, _ := utf8.DecodeRune([]byte{26})
	log.Println("got connection")
	defer func() {
		log.Println("close connection")
		cn.Close()
	}()
	for scanner := bufio.NewScanner(cn); scanner.Scan(); {
		if scanner.Text() == string(eof) {
			return
		}
		fmt.Printf("%s\n", scanner.Text())
	}
}
