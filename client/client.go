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
	"time"

	"github.com/Microsoft/go-winio"
)

func main() {
	pipeName := flag.String("p", "", "pipe name")
	flag.Parse()
	pipeClient(*pipeName)
}

func pipeClient(pipeName string) {
	fullPipeName := fmt.Sprintf(`\\.\pipe\%s`, pipeName)
	log.Println(fullPipeName)
	var conn net.Conn
	var err error
	go func() {
		for {
			conn, err = winio.DialPipe(fullPipeName, nil)
			if err != nil {
				log.Println("dial with error:", err, "sleep 1 second")
				time.Sleep(1 * time.Second)
				continue
			} else {
				break
			}
		}
		log.Println("start connection")
		for scanner := bufio.NewScanner(os.Stdin); scanner.Scan(); {
			fmt.Fprintln(conn, scanner.Text())
		}
	}()
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(c)
	<-c
	log.Println("shutdown")
	if conn != nil {
		log.Println("close connection")
		defer conn.Close()
	}
}
