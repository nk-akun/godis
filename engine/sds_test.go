package godis

import "fmt"

// TestSds ...
func TestSds() {

	str := "abcdefse"
	s1 := SdsNewString(&str)
	str = "abcdefs"
	s2 := SdsNewString(&str)

	//cmp
	fmt.Println(SdsCmp(s1, s2))

	// len
	fmt.Println(s1.SdsLen(), s2.SdsLen())

	//get string
	fmt.Println(*s1.SdsGetString())

	//copy
	str = "poiuytersds"
	s1.SdsCopy(&str)

	fmt.Println(s1.SdsLen(), *s1.SdsGetString())

	//cat
	str = "123213"
	s1.SdsCat(&str)
	fmt.Println(s1.SdsLen(), *s1.SdsGetString())

	s1.SdsCatSds(s2)
	fmt.Println(s1.SdsLen(), *s1.SdsGetString())
}
