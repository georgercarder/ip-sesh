package main

import (
	"fmt"
	"net"
	"bufio"

	"os"
)

func main() {
	fmt.Println("Press ENTER for demo.")
	getCharReader := bufio.NewReader(os.Stdin)
	domain, err := getCharReader.ReadString('\n')
	if err != nil {
		fmt.Println("debug error", err)
	}
	fmt.Println("debug domain", domain)
	//nd.StartHandshake("test.domain.com")
	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		// TODO LOG, AND GRACEFUL
		panic(err)
	}
	go func() {
		b := make([]byte, 1024)
		n, err := conn.Read(b)
		fmt.Println("debug", n, string(b), err)
	}()
	conn.Write([]byte(domain))
	fmt.Println("debug conn", conn)
	select {}
}
