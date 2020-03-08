package main

import "fmt"

func main() {

	buf := make([]byte, 3)
	s := "abcdef"

	n := copy(buf[2:], s)

	fmt.Println(s)
	fmt.Println(n)
	fmt.Println(len(buf))
	fmt.Println(buf[:])
}
