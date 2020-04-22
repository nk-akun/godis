package godis

import (
	"fmt"
	"strconv"

	"github.com/nk-akun/godis/engine/util"
)

const (
	DICT_STEP_HASH_SIZE  = 1
	DICT_HT_INITIAL_SIZE = 4
	DICT_RESIZE_RATIO    = 5
)

const (
	DICT_ERROR = 1
	DICT_OK    = 0
)

// DictNode stores key and value
type DictNode struct {
	key   *Object
	value *Object
	next  *DictNode
}

// DictHT stores a hash table
type DictHT struct {
	table    []*DictNode
	size     uint32
	sizeMask uint32
	used     uint32
}

// DictFunc stores functions of dict
type DictFunc struct {
	calHash    func(*Object) uint32
	keyCompare func(*Object, *Object) int
}

// Dict stores
type Dict struct {
	ht          [2]*DictHT
	funcs       *DictFunc
	rehashIndex int64
	iterators   uint32
}

// DictIterator is used to traverse dict
type DictIterator struct {
	dt         *Dict
	index      int64
	table      int
	safe       bool
	node       *DictNode
	nextNode   *DictNode
	stateLable uint32
}

// NewDict return a new dict pointer
func NewDict(funcs *DictFunc) *Dict {
	d := new(Dict)
	d.ht[0] = NewDictHT(0)
	d.ht[1] = NewDictHT(0)
	d.funcs = funcs
	d.rehashIndex = -1
	d.iterators = 0
	return d
}

// NewDictNode return a new dict node
func NewDictNode(key *Object, value *Object) *DictNode {
	return &DictNode{
		key:   key,
		value: value,
	}
}

// NewDictHT ...
func NewDictHT(size uint32) *DictHT {
	if size == 0 {
		return &DictHT{
			table:    nil,
			size:     0,
			sizeMask: 0,
			used:     0,
		}
	}
	ht := &DictHT{
		table:    make([]*DictNode, size),
		size:     size,
		sizeMask: size - 1,
		used:     0,
	}
	return ht
}

func (ht *DictHT) reset() {
	ht.table = nil
	ht.size = 0
	ht.sizeMask = 0
	ht.used = 0
}

// Add add key value to dict d
func (d *Dict) Add(key *Object, value *Object) int {
	node := d.addRaw(key, value)
	if node == nil {
		return DICT_ERROR
	}
	return DICT_OK
}

// Get return value by key or nil if can't find the node
func (d *Dict) Get(key *Object) *Object {
	node := d.Search(key)
	return node.value
}

// Delete delete node with given key
func (d *Dict) Delete(key *Object) int {
	return d.deleteNode(key)
}

func (d *Dict) addRaw(key *Object, value *Object) *DictNode {
	// incorporate rehash into add action
	if d.isRehashing() {
		d.rehashStep(DICT_STEP_HASH_SIZE)
	}

	d.expandIfFull()

	var index uint32
	var ht *DictHT
	if i := d.getIndexKey(key, d.funcs.calHash(key)); i == -1 {
		return nil
	} else {
		index = uint32(i)
	}

	if d.isRehashing() {
		ht = d.ht[1]
	} else {
		ht = d.ht[0]
	}

	node := NewDictNode(key, value)
	node.next = ht.table[index]
	ht.table[index] = node
	ht.used++
	return node
}

func (d *Dict) deleteNode(key *Object) int {
	if d.ht[0].used+d.ht[1].used == 0 { // empty dict
		return DICT_ERROR
	}
	if d.isRehashing() {
		d.rehashStep(DICT_STEP_HASH_SIZE)
	}

	hashValue := d.funcs.calHash(key)
	for i := 0; i < 2; i++ {
		index := hashValue & d.ht[i].sizeMask
		node := d.ht[i].table[index]
		var preNode *DictNode = nil

		for node != nil {
			if node.key == key || d.funcs.keyCompare(key, node.key) == 0 {
				if preNode == nil { // head
					d.ht[i].table[index] = node.next
				} else { // not head
					preNode.next = node.next
				}
				ReleaseDictNode(node)
				d.ht[i].used--
				return DICT_OK
			}
			preNode = node
			node = node.next
		}
	}
	return DICT_ERROR
}

func (d *Dict) isRehashing() bool {
	return d.rehashIndex != -1
}

func (d *Dict) rehashStep(num int) {
	// if there are visitors traversing this dict
	if d.iterators > 0 {
		return
	}

	for ; num > 0 && d.ht[0].used > 0; num-- {
		for d.ht[0].table[d.rehashIndex] == nil {
			d.rehashIndex++
		}

		node := d.ht[0].table[d.rehashIndex]
		for node != nil {
			temp := node.next
			// get hash value
			id := d.funcs.calHash(node.key) & d.ht[1].sizeMask

			// insert into corect position in ht[1]
			node.next = d.ht[1].table[id]
			d.ht[1].table[id] = node
			d.ht[0].used--
			d.ht[1].used++
			node = temp
		}
		d.ht[0].table[d.rehashIndex] = nil
		d.rehashIndex++
	}

	// if rehash is finished
	if d.ht[0].used == 0 {
		d.ht[0] = d.ht[1]
		d.rehashIndex = -1
		d.ht[1] = NewDictHT(0)
	}
}

func (d *Dict) getIndexKey(key *Object, hashValue uint32) int64 {
	var index uint32
	for i := 0; i < 2; i++ {
		index = hashValue & d.ht[i].sizeMask
		node := d.ht[i].table[index]
		for node != nil {
			if node.key == key || d.funcs.keyCompare(key, d.ht[i].table[index].key) == 0 {
				return -1
			}
			node = node.next
		}
		if !d.isRehashing() {
			break
		}
	}
	return int64(index)
}

func (d *Dict) expandIfFull() {
	if d.isRehashing() {
		return
	}
	if d.ht[0].size == 0 {
		d.expand(DICT_HT_INITIAL_SIZE)
	} else if d.ht[0].used/d.ht[0].size >= DICT_RESIZE_RATIO {
		d.expand(d.ht[0].size << 1)
	}
}

func (d *Dict) expand(size uint32) {
	size = util.NearLargeUnsignedBinary(size)
	if size <= d.ht[0].size {
		return
	}

	ht := NewDictHT(size)

	// if dict has not been used,the initialization is necessary
	if d.ht[0].table == nil {
		d.ht[0] = ht
		return
	}

	d.ht[1] = ht
	d.rehashIndex = 0
}

// Search find dict node with the key
func (d *Dict) Search(key *Object) *DictNode {
	if d.ht[0].used+d.ht[1].used == 0 {
		return nil
	}
	if d.isRehashing() {
		d.rehashStep(DICT_STEP_HASH_SIZE)
	}
	hashValue := d.funcs.calHash(key)
	for i := 0; i < 2; i++ {
		index := hashValue & d.ht[i].sizeMask
		node := d.ht[i].table[index]
		for node != nil {
			if node.key == key || d.funcs.keyCompare(node.key, key) == 0 {
				return node
			}
			node = node.next
		}
		if !d.isRehashing() {
			break
		}
	}
	return nil
}

// NewDictIterator return a new iterator
func NewDictIterator(d *Dict) *DictIterator {
	return &DictIterator{
		dt:       d,
		table:    0,
		index:    -1,
		safe:     false,
		node:     nil,
		nextNode: nil,
	}
}

// NewSafeDictIterator return a new iterator which is safe
func NewSafeDictIterator(d *Dict) *DictIterator {
	iter := NewDictIterator(d)
	iter.safe = true
	return iter
}

// ReleaseIterator release interator so dict could rehash
// or check whether dict has changed or not
func ReleaseIterator(iter *DictIterator) {
	if !(iter.index == -1 && iter.table == 0) {
		if iter.safe {
			iter.dt.iterators--
		} else if iter.stateLable != iter.dt.stateLable() {
			log.Errorf("dict stateTable has changed when iterator traverse")
		}
	}
	iter = nil
}

// ReleaseDictNode release dict node
func ReleaseDictNode(node *DictNode) {
	node.key = nil
	node.value = nil
	node.next = nil
	node = nil
}

// Next return the next node of node
func (iter *DictIterator) Next() *DictNode {
	for {
		if iter.node == nil {
			ht := iter.dt.ht[iter.table]
			if iter.table == 0 && iter.index == -1 { // using iter for the first time
				if iter.safe {
					iter.dt.iterators++
				} else {
					iter.stateLable = iter.dt.stateLable()
				}
			}

			iter.index++
			if iter.index >= int64(ht.size) {
				if iter.table == 0 && iter.dt.isRehashing() { // next table
					iter.table++
					iter.index = 0
					ht = iter.dt.ht[iter.table]
				} else {
					break
				}
			}
			iter.node = ht.table[iter.index]
		} else {
			iter.node = iter.nextNode
		}
		if iter.node != nil {
			iter.nextNode = iter.node.next
			return iter.node
		}
	}
	return nil
}

func (d *Dict) stateLable() uint32 {
	nums := make([]uint32, 6)
	addr := fmt.Sprintf("%p", d.ht[0])
	tmp, _ := strconv.ParseUint(addr[2:], 16, 32)
	nums[0] = uint32(tmp)
	nums[1] = d.ht[0].size
	nums[2] = d.ht[0].used
	addr = fmt.Sprintf("%p", d.ht[1])
	tmp, _ = strconv.ParseUint(addr[2:], 16, 32)
	nums[3] = uint32(tmp)
	nums[4] = d.ht[1].size
	nums[5] = d.ht[1].used

	var hash uint32
	for i := 0; i < 6; i++ {
		hash += nums[i]
		hash = (^hash) + (hash << 21)
		hash = hash ^ (hash >> 24)
		hash = (hash + (hash << 3)) + (hash << 8)
		hash = hash ^ (hash >> 14)
		hash = (hash + (hash << 2)) + (hash << 4)
		hash = hash ^ (hash >> 28)
		hash = hash + (hash << 31)
	}
	return hash
}
