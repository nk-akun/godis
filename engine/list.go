package godis

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

	// TODO:dup free match functions
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

// PrevNode return the previous node of node n
func (n *ListNode) PrevNode() *ListNode {
	return n.prev
}

// NextNode return the next node of node n
func (n *ListNode) NextNode() *ListNode {
	return n.next
}

// Value return value of node n
func (n *ListNode) Value() *Object {
	return n.value
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
func (l *List) AddNodeHead(n *ListNode) {
	if l.First() == nil {
		l.head, l.tail = n, n
	} else {
		l.head.prev = n
		n.next = l.head
		l.head = n
	}
	l.length++
}

// AddNodeTail add a node at end of list l
func (l *List) AddNodeTail(n *ListNode) {
	if l.Last() == nil {
		l.head, l.tail = n, n
	} else {
		l.tail.next = n
		n.prev = l.tail
		l.tail = n
	}
	l.length++
}
