package godis

import (
	"bytes"
	"fmt"
)

func Test_ParseConf() {
	ParseConf()
	fmt.Println(godisConf)
}

func Test_Bufio2() {
	buf := make([]byte, 20)
	b := bytes.Buffer{}

	copy(buf, "asdf")
	b.Write(buf)

	fmt.Println(b.String())
}
