package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	addr := "127.0.0.1:8080"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	// buf := make([]byte, 1024)
	// len, err := conn.Read(buf)
	// if err != nil && err != io.EOF {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(buf[:len]))

	output := bufio.NewReader(conn)

	go send(conn)
	for {
		// fmt.Printf("Please input >>> ")

		content, err := output.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Println(content)
	}

	conn.Close()
}

func send(conn net.Conn) {
	for {

		input := bufio.NewReader(os.Stdin)
		line, _ := input.ReadString('\n')

		if line == "byebye123" {
			break
		}

		// fmt.Printf(";;;;;;;;;;;;;;%s", line)

		_, err := conn.Write([]byte(line))
		if err != nil {
			break
		}
	}
}
