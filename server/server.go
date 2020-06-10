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

func pipeServer(cxt context.Context, pipeName string) {
	fullPipeName := fmt.Sprintf(`\\.\pipe\%s`, pipeName)
	log.Println(fullPipeName)
	listener, err := winio.ListenPipe(fullPipeName, nil)
	if err != nil {
		log.Println(err)
		return
	}
	isListenerClose := false
	defer func() {
		log.Println("close listener")
		isListenerClose = true
		listener.Close()
	}()

	go func() {
		for {
			// don't block waiting stop event
			conn, err := listener.Accept()
			if err != nil {
				if isListenerClose {
					break
				} else {
					log.Println("listen with error:", err)
					continue
				}
			}
			go func(cn net.Conn) {
				log.Println("got connection")
				defer func() {
					log.Println("close connection")
					(cn).Close()
				}()
				for scanner := bufio.NewScanner(cn); scanner.Scan(); {
					fmt.Println(scanner.Text())
				}

				// wait stop event
				for {
					select {
					case <-cxt.Done():
						log.Println("connection receive stop event")
						return
					default:
					}
				}
			}(conn)
		}
	}()
	for {
		// wait stop event
		select {
		case <-cxt.Done():
			log.Println("listener receive stop event")
			return
		default:
		}
	}
}
