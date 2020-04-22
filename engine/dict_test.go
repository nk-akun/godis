package godis

import (
	"fmt"
	"math/rand"
	"time"
)

// TestDict ...
func TestDict() {
	dtf := &DictFunc{}
	dtf.calHash = CalHashCommon
	dtf.keyCompare = CompareValueCommon

	rand.Seed(time.Now().Unix())
	keys := make([]string, 0, 100)
	dt := NewDict(dtf)
	for i := 0; i < 100; i++ {
		s := ""
		for j := 0; j < 7; j++ {
			s += string('a' + rand.Intn(26))
		}
		dt.Add(NewObject(OBJString, s), NewObject(OBJInt, rand.Intn(100007)))
		keys = append(keys, s)
	}

	outputDict(dt)

	for i := 0; i < 50; i++ {
		fmt.Println(dt.Delete(NewObject(OBJString, keys[i])))
	}

	outputDict(dt)
}

func outputDict(dt *Dict) {
	iter := NewDictIterator(dt)
	sum := 0
	for node := iter.Next(); node != nil; node = iter.Next() {
		value := dt.Get(node.key)
		num, _ := value.Ptr.(int)
		fmt.Printf("%d ", num)
		sum++
	}
	fmt.Println("count:", sum)
}
