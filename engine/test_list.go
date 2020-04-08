package godis

import "fmt"

// TestList ...
func TestList() {
	match := func(v1 *Object, v2 *Object) bool {
		if v1.ObjectType != v2.ObjectType {
			return false
		}
		num1, ok := v1.Ptr.(int)
		if !ok {
			return false
		}
		num2, ok := v2.Ptr.(int)
		if !ok {
			return false
		}

		if num1 == num2 {
			return true
		}
		return false

	}

	list := NewList()
	list.match = match

	//AddNodeTail
	fmt.Println("add node tail 0~10")
	for i := 1; i <= 5; i++ {
		list.AddNodeTail(NewObject(OBJInt, i))
	}
	testListOutputList(list)

	//AddNodeHead
	fmt.Println("add node head 11~20")
	for i := 6; i <= 10; i++ {
		list.AddNodeHead(NewObject(OBJInt, i))
	}
	testListOutputList(list)

	//InsertNode

	// node := list.head
	// for i := 11; node != nil && i <= 20; node = node.next {
	// 	list.InsertNode(node, NewObject(OBJInt, i), true)
	// 	node = node.next
	// 	i++
	// }
	// testListOutputList(list)

	node := list.tail
	for i := 11; node != nil && i <= 20; node = node.prev {
		list.InsertNode(node, NewObject(OBJInt, i), false)
		node = node.prev
		i++
	}
	testListOutputList(list)

	// Index
	testListOutputNode(list.Index(-3))
	testListOutputNode(list.Index(4))

	// DelNode
	list.DelNode(list.Index(4))
	testListOutputNode(list.Index(4))
	list.DelNode(list.Index(-3))
	testListOutputNode(list.Index(-3))

	testListOutputList(list)
	for i := 0; i < 5; i++ {
		list.Rotate()
	}
	testListOutputList(list)

	//SearchKey
	node = list.SearchKey(NewObject(OBJInt, 100))
	if node != nil {
		testListOutputNode(node)
	}
}

func testListOutputNode(node *ListNode) {
	value := node.Value()
	num, ok := value.Ptr.(int)
	if !ok {
		return
	}
	fmt.Printf("%d\n", num)
}

func testListOutputList(l *List) {
	fmt.Println(l.length)
	node := l.head
	for ; node != nil; node = node.next {
		value := node.Value()
		num, ok := value.Ptr.(int)
		if !ok {
			return
		}
		fmt.Printf("%d ", num)
	}
	fmt.Println()
}
