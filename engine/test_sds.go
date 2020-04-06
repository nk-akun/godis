package godis

import "fmt"

// TestSds ...
func TestSds() {

	str := "abcdef"
	s1 := SdsNew(&str)
	str = "qwertt"
	s2 := SdsNew(&str)

	fmt.Println(s1.SdsLen(), s2.SdsLen())
}
