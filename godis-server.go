package main

import (
	"fmt"
	"net"
	"os"

	engine "github.com/nk-akun/godis/engine"
)

var server *engine.Server

func main() {
	engine.InitLogger("/Users/marathon/Work/mygo/src/github.com/nk-akun/godis/log", 1, 1, 2, false)
	listener, err := net.Listen("tcp", "127.0.0.1:10010")
	if err != nil {
		fmt.Println("listen :%v", err)
		os.Exit(-1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	client := server.CreateClient()
	for {
		client.ReadClientContent(conn)
		client.TransClientContent()
		responseClient(conn, client)
	}
}

func responseClient(conn net.Conn, c *engine.Client) {
	ans := ""
	for _, t := range c.Argv {
		ans += t.Ptr.(string)
	}
	conn.Write([]byte(ans))
}
