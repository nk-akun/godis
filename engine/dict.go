package godis

const (
	DICT_STEP_HASH_SIZE = 1
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
	size     uint64
	sizeMask uint64
	used     uint64
}

// DictFunc stores functions of dict
type DictFunc struct {
	calHash    func(key *Object) uint64
	keyCompare func(key1 *Object, key2 *Object) int
}

// Dict stores
type Dict struct {
	ht          [2]*DictHT
	funcs       *DictFunc
	rehashIndex int64
	iterators   uint64
}

// NewDict return a new dict pointer
func NewDict(funcs *DictFunc) *Dict {
	d := new(Dict)
	d.ht[0].init()
	d.ht[1].init()
	d.funcs = funcs
	d.rehashIndex = -1
	d.iterators = 0
	return d
}

func (ht *DictHT) init() {
	ht = &DictHT{
		table:    make([]*DictNode, 0), // init 2
		size:     0,
		sizeMask: 0,
		used:     0,
	}
}

// Add add key value to dict d
func (d *Dict) Add(key *Object, value *Object) {
	d.addRaw(key, value)
}

func (d *Dict) addRaw(key *Object, value *Object) {
	// incorporate rehash into add action
	if d.isRehashing() {
		d.rehashStep(DICT_STEP_HASH_SIZE)
	}

	var index uint64
	var ht *DictHT
	if i := d.getIndexKey(key, d.funcs.calHash(key)); i == -1 {
		return
	index = uint64(i)

	if d.isRehashing() {
		ht = d.ht[1]
	} else {
		ht = d.ht[0]
	}	
	
}

func (d *Dict) isRehashing() bool {
	return d.rehashIndex != -1
}

func (d *Dict) rehashStep(num int) {
	// if there are visitors traversing this dict
	if d.iterators >= 0 {
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
		d.ht[1].init()
	}
}

func (d *Dict) getIndexKey(key *Object, hashValue uint64) int64 {
	var index uint64
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

// Search find dict node with the key
func (d *Dict) Search(key *Object) *DictNode {
	return nil
}
