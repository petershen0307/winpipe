package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

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
	conn, err := winio.DialPipe(fullPipeName, nil)
	if err == nil {
		defer conn.Close()
	}
	for {
		for scanner := bufio.NewScanner(os.Stdin); scanner.Scan(); {
			fmt.Fprintln(conn, scanner.Text())
		}
	}
}
