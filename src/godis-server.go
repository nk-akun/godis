package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var room *groom = newRoom()

type groom struct {
	users map[string]net.Conn
}

func newRoom() *groom {
	return &groom{
		users: make(map[string]net.Conn),
	}
}

func (room *groom) leave(user string) {
	conn, ok := room.users[user]
	if !ok {
		return
	}

	conn.Close()
	delete(room.users, user)
	fmt.Printf("%s leave room\n", user)
}

func (room *groom) join(user string, conn net.Conn) {
	if _, ok := room.users[user]; ok {
		room.leave(user)
	}
	room.users[user] = conn
	fmt.Printf("%s join the room\n", user)
	conn.Write([]byte(fmt.Sprintf("%s join the room\n", user)))
}

func (room *groom) broadcast(user string, msg string) {

	timeInfo := time.Now().Format("2006年01月02日 15:04:05")
	allMsg := fmt.Sprintf("%v %s:%s\n", timeInfo, user, msg)
	for username, conn := range room.users {
		if username == user {
			continue
		}
		conn.Write([]byte(allMsg))
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	line, err := r.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		return
	}

	line = strings.TrimSpace(line)
	fileds := strings.Fields(line)

	if len(fileds) != 2 {
		fmt.Println("Username or Password is wrong!")
		return
	}

	room.join(fileds[0], conn)
	room.broadcast(fileds[0], fmt.Sprintf("welcome %s join the room\n", fileds[0]))

	for {
		conn.Write([]byte("Please input message >>> \n"))
		content, err := r.ReadString('\n')
		if err != nil {
			break
		}
		content = strings.TrimSpace(content)
		room.broadcast(fileds[0], content)
	}
	room.broadcast("system", fmt.Sprintf("%s leave the room\n", fileds[0]))
	room.leave(fileds[0])
}

func main() {
	addr := "127.0.0.1:8080"
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handle(conn)
	}
}
