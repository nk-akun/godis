package main

import (
	"fmt"
	"net"
	"os"

	engine "github.com/nk-akun/godis/engine"
)

const (
	// DefaultLogFile log file
	DefaultLogFile = "/Users/marathon/Work/mygo/src/github.com/nk-akun/godis/log"
)

var server *engine.Server

func main() {
	engine.InitLogger(DefaultLogFile, 1, 1, 2, false)

	server = engine.InitServer()

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
		err := client.ReadClientContent(conn)
		if err != nil {
			engine.GetGodisLogger().Errorf("read query content error:%+v", err)
			return
		}
		err = client.TransClientContent()
		if err != nil {
			engine.GetGodisLogger().Errorf("process input buffer error:%+v", err)
		}
		server.ProcessCommand(client)
		responseClient(conn, client)
	}
}

func responseClient(conn net.Conn, c *engine.Client) {
	conn.Write([]byte(*c.Buf.SdsGetString()))
}
