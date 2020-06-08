package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	winio "github.com/Microsoft/go-winio"
)

func main() {
	pipeName := flag.String("p", "", "pipe name")
	flag.Parse()
	pipeServer(*pipeName)
}

func pipeServer(pipeName string) {
	fullPipeName := fmt.Sprintf(`\\.\pipe\%s`, pipeName)
	log.Println(fullPipeName)
	listener, err := winio.ListenPipe(fullPipeName, nil)
	if err != nil {
		log.Println(err)
		return
	}
	closeFlag := false
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(c)
		defer listener.Close()
		<-c
		log.Println("shutdown")
		closeFlag = true
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if closeFlag {
				break
			} else {
				log.Println("listen with error:", err)
				continue
			}
		}
		go func(cn *net.Conn) {
			log.Println("got connection")
			for scanner := bufio.NewScanner(*cn); scanner.Scan(); {
				fmt.Println(scanner.Text())
			}
			defer func() {
				log.Println("close connection")
				(*cn).Close()
			}()
		}(&conn)
	}
}
