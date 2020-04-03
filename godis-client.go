package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	godis "github.com/nk-akun/godis/engine"
)

func main() {
	addr := "127.0.0.1:10010"

	reader := bufio.NewReader(os.Stdin)

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println("error ", err)
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("error ", err)
		os.Exit(1)
	}

	for {
		fmt.Printf("%s>", addr)
		content, _ := reader.ReadString('\n')
		content = strings.Replace(content, "\n", "", -1)

		sendToServer(conn, content)

		buff := make([]byte, 1024)
		conn.Read(buff)
		fmt.Println(string(buff))
	}
}

func sendToServer(conn *net.TCPConn, content string) (n int, err error) {
	b, err := godis.EncodeCmd(content)
	if err != nil {
		return 0, err
	}

	return conn.Write(b)
}
