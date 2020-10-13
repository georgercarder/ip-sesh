package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	. "github.com/georgercarder/ip-sesh/common"
	sh "github.com/georgercarder/ip-sesh/shell"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		// TODO LOG, AND GRACEFUL
		panic(err)
	}
	fmt.Println("Press ENTER for demo.")
	getCharReader := bufio.NewReader(os.Stdin)
	domain, err := getCharReader.ReadString('\n')
	if err != nil {
		fmt.Println("debug error", err)
	}
	fmt.Println("debug domain", domain)
	conn.Write([]byte(domain))
	//nd.StartHandshake("test.domain.com")
	b := make([]byte, 1024)
	_, err = conn.Read(b)
	if Trim(string(b)) == sh.StartShellSession {
		fmt.Println("debug START SHELL SESSION HERE")
		sh.Client(conn)
	}

	select {}
}
