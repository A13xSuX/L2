package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Host    string
	Port    int
	Timeout time.Duration
}

// to stdout from socket
func readFromSocket(conn net.Conn, done chan struct{}) {
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			if err != io.EOF {
				fmt.Println("readFromSocket err:", err)
			}
			done <- struct{}{}
			return
		}

		fmt.Print(line)
	}
}

func writeToSocket(conn net.Conn, done chan struct{}) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text() + "\n"
		_, err := conn.Write([]byte(text))

		if err != nil {
			fmt.Println("writeToSocket err:", err)
			done <- struct{}{}
			return
		}
	}

	if scanner.Err() != nil { //это ctrl+d
		fmt.Println("writeToSocket err:", scanner.Err())
	}
	done <- struct{}{}
}

func main() {
	var config Config

	flag.StringVar(&config.Host, "host", "", "Host")
	flag.IntVar(&config.Port, "port", 0, "Port")
	timeout := flag.Int("timeout", 10, "Timeout")

	flag.Parse()

	config.Timeout = time.Duration(*timeout) * time.Second // получше разобраться

	//валидация
	if config.Host == "" {
		fmt.Fprintln(os.Stderr, "Host is required")
		os.Exit(1)
	}
	if config.Port == 0 {
		fmt.Fprintln(os.Stderr, "Port is required")
		os.Exit(1)
	}

	// установка соединения
	address := config.Host + ":" + strconv.Itoa(config.Port)
	connect, err := net.DialTimeout("tcp", address, config.Timeout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Connection error:", err)
		os.Exit(1)
	}
	defer connect.Close()

	//graceful shutdown
	done := make(chan struct{})

	go readFromSocket(connect, done)
	go writeToSocket(connect, done)

	<-done

}
