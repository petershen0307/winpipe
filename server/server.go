package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"

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
		fmt.Println(err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("exit with error:", err)
			return
		}

		for scanner := bufio.NewScanner(conn); scanner.Scan(); {
			fmt.Println(scanner.Text())
		}
		conn.Close()
	}
}
