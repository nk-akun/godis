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

// CalHashCommon is the common function to calculate the hash value for the string type
func CalHashCommon(v *Object) uint32 {
	switch v.ObjectType {
	case OBJSDS:
		key, _ := v.Ptr.(*Sdshdr)
		return GetHashSds(key)
	case OBJString:
		key, _ := v.Ptr.(string)
		return GetHashString(key)
	}
	return 0
}

// CompareValueCommon compare value if the type if int, string or sds
func CompareValueCommon(v1 *Object, v2 *Object) int {
	switch v1.ObjectType {
	case OBJSDS:
		value1, _ := v1.Ptr.(*Sdshdr)
		value2, _ := v2.Ptr.(*Sdshdr)
		return SdsCmp(value1, value2)
	case OBJString:
		value1, _ := v1.Ptr.(string)
		value2, _ := v2.Ptr.(string)
		if value1 < value2 {
			return -1
		} else if value1 == value2 {
			return 0
		}
		return 1
	}
	return -1
}
