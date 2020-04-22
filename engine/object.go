package godis

// Object stores data whose type is Object.Type
type Object struct {
	ObjectType int
	Ptr        interface{}
}

const OBJString = 0
const OBJSDS = 1
const OBJInt = 2
const OBJSet = 3
const OBJZset = 4
const OBJHash = 5
const OBJList = 6
const OBJCommand = 7

// NewObject return a new Object
func NewObject(tp int, ptr interface{}) *Object {
	o := new(Object)
	o.ObjectType = tp
	o.Ptr = ptr
	return o
}
