package main

import (
	"fmt"

	godis "github.com/nk-akun/godis/engine"
)

func main() {

	// godis.TestSds()
	// godis.TestList()
	// godis.TestDict()

	str := "I'm a cool boy!"
	fmt.Println(godis.GetHashString(str))
	str = "I'm a cool boz!"
	fmt.Println(godis.GetHashString(str))
}
