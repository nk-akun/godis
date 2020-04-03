package godis

// Object stores data whose type is Object.Type
type Object struct {
	ObjectType int
	Ptr        interface{}
}

const OBJString = 0
const OBJInt = 1
const OBJSet = 2
const OBJZset = 3
const OBJHash = 4
const OBJList = 5

// NewObject return a new Object
func NewObject(tp int, ptr interface{}) *Object {
	o := new(Object)
	o.ObjectType = tp
	o.Ptr = ptr
	return o
}
