package godis

import "github.com/twmb/murmur3"

// GetHashSds return the hash value of given string
func GetHashSds(sds *Sdshdr) uint32 {
	str := sds.SdsGetString()
	return murmur3.StringSum32(*str)
}

// GetHashString ...
func GetHashString(str string) uint32 {
	return murmur3.StringSum32(str)
}
