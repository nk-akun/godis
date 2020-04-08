package godis

const (
	LIST_START_HEAD = 0
	LIST_START_TAIL = 1
)

// ListNode stores data of node in list
type ListNode struct {
	prev  *ListNode
	next  *ListNode
	value *Object
}

// List is made by joining ListNode
type List struct {
	head   *ListNode
	tail   *ListNode
	length int64
	match  func(*Object, *Object) bool
	// TODO:dup free functions
}

// ListIter is iterater of list
type ListIter struct {
	next      *ListNode
	direction int
}

// Length return the length of list l
func (l *List) Length() int64 {
	return l.length
}

// First return the first node of list l
func (l *List) First() *ListNode {
	return l.head
}

// Last return the last node of list l
func (l *List) Last() *ListNode {
	return l.tail
}

// NewList create a new list,return list's pointer
func NewList() *List {
	return &List{
		head:   nil,
		tail:   nil,
		length: 0,
	}
}

// AddNodeHead add a node at head of list l
func (l *List) AddNodeHead(v *Object) {
	node := NewNode(v)
	if l.First() == nil {
		l.head, l.tail = node, node
	} else {
		l.head.prev = node
		node.next = l.head
		l.head = node
	}
	l.length++
}

// AddNodeTail add a node at end of list l
func (l *List) AddNodeTail(v *Object) {
	node := NewNode(v)
	if l.Last() == nil {
		l.head, l.tail = node, node
	} else {
		l.tail.next = node
		node.prev = l.tail
		l.tail = node
	}
	l.length++
}

// InsertNode insert node that value is v after oldNode if after is true else before oldNode
func (l *List) InsertNode(oldNode *ListNode, v *Object, after bool) {
	node := NewNode(v)
	if after {
		node.prev = oldNode
		node.next = oldNode.next
		if oldNode == l.tail {
			l.tail = node
		}
	} else {
		node.next = oldNode
		node.prev = oldNode.prev
		if oldNode == l.head {
			l.head = node
		}
	}
	if node.prev != nil {
		node.prev.next = node
	}
	if node.next != nil {
		node.next.prev = node
	}
}

// Index return the element index where 0 is the head,1 is the element next
// to head, so on.-1 is the tail,-2 is the previous element, so on.
func (l *List) Index(index int64) (node *ListNode) {
	if index >= 0 {
		node = l.head
		for ; node != nil && index > 0; node = node.next {
			index--
		}
	} else {
		index = -index - 1
		node = l.tail
		for ; node != nil && index > 0; node = node.prev {
			index--
		}
	}
	return node
}

// DelNode delete node in list
func (l *List) DelNode(node *ListNode) {
	// very good
	if node.prev == nil {
		l.head = node.next
	} else {
		node.prev.next = node.next
	}
	if node.next == nil {
		l.tail = node.prev
	} else {
		node.next.prev = node.prev
	}
}

// Rotate take tail element to head
func (l *List) Rotate() {
	if l.tail == nil || l.head == l.tail {
		return
	}

	node := l.tail

	l.tail = node.prev
	node.prev.next = nil
	node.prev = nil

	node.next = l.head
	l.head.prev = node
	l.head = node
}

// SearchKey search the ListNode
func (l *List) SearchKey(key *Object) *ListNode {
	iter := l.RewindHead()
	var node *ListNode
	for {
		node = iter.NextNode()
		if node == nil {
			break
		}
		if l.match != nil {
			if l.match(node.value, key) {
				return node
			}
		} else {
			if node.value == key {
				return node
			}
		}
	}
	return nil
}

// RewindHead return iterater that rewind from head to tail
func (l *List) RewindHead() *ListIter {
	return &ListIter{
		next:      l.head,
		direction: LIST_START_HEAD,
	}
}

// RewindTail return iterater that rewind from tail to head
func (l *List) RewindTail() *ListIter {
	return &ListIter{
		next:      l.tail,
		direction: LIST_START_TAIL,
	}
}

// ListNode

// NewNode return a new node whose node is value
func NewNode(value *Object) *ListNode {
	return &ListNode{
		prev:  nil,
		next:  nil,
		value: value,
	}
}

// PrevNode return the previous node of node n
func (node *ListNode) PrevNode() *ListNode {
	return node.prev
}

// NextNode return the next node of node n
func (node *ListNode) NextNode() *ListNode {
	return node.next
}

// Value return value of node n
func (node *ListNode) Value() *Object {
	return node.value
}

// NextNode return next node stored in iter
func (iter *ListIter) NextNode() *ListNode {
	node := iter.next
	if node != nil {
		if iter.direction == LIST_START_HEAD {
			iter.next = node.next
		} else {
			iter.next = node.prev
		}
	}
	return node
}
