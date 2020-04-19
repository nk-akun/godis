package godis

import (
	"math/rand"
	"time"
)

const (
	ZSL_MAX_LEVEL = (1 << 5)
)

// ZskipList ...
type ZskipList struct {
	header *ZskipListNode
	tail   *ZskipListNode
	length uint32
	level  int
}

// ZskipListNode ...
type ZskipListNode struct {
	value    *Sdshdr
	score    float64
	backward *ZskipListNode
	level    []zskipListLevel
}

type zskipListLevel struct {
	forward *ZskipListNode
	span    uint32
}

// NewZsl create a new zskipList
func NewZsl() *ZskipList {
	zsl := &ZskipList{
		tail:   nil,
		length: 0,
		level:  ZSL_MAX_LEVEL,
	}
	zsl.header = NewZslNode(ZSL_MAX_LEVEL, nil, 0)
	zsl.header.backward = nil
	for i := 0; i < ZSL_MAX_LEVEL; i++ {
		zsl.header.level[i].forward = nil
		zsl.header.level[i].span = 0
	}
	return zsl
}

// NewZslNode create a new skipList node
func NewZslNode(lv int, value *Sdshdr, score float64) *ZskipListNode {
	return &ZskipListNode{
		level: make([]zskipListLevel, lv),
		value: value,
		score: score,
	}
}

func zslRandomLevel() int {
	rand.Seed(time.Now().Unix()) //TODO: change setting seed mod to others
	return rand.Intn(ZSL_MAX_LEVEL) + 1
}

// Insert insert node into skiplist
func (zsl *ZskipList) Insert(score float64, value *Sdshdr) *ZskipListNode {
	// borders[i] means we need insert node after borders[i] when level is i
	// dis[i] is distence from head to borders[i]
	borders := make([]*ZskipListNode, ZSL_MAX_LEVEL)
	dis := make([]uint32, ZSL_MAX_LEVEL)

	p := zsl.header
	for i := zsl.level - 1; i >= 0; i-- {
		if i == zsl.level-1 {
			dis[i] = 0
		} else {
			dis[i] = dis[i+1]
		}
		for p.level[i].forward != nil && (p.level[i].forward.score < score || (p.level[i].forward.score == score && SdsCmp(p.level[i].forward.value, value) == -1)) {
			dis[i] += p.level[i].span
			p = p.level[i].forward
		}
		borders[i] = p.level[i].forward
	}

	level := zslRandomLevel()
	node := NewZslNode(level, value, score)
	for i := 0; i < level; i++ {
		//update forward
		node.level[i].forward = borders[i].level[i].forward
		borders[i].level[i].forward = node

		//update span
		node.level[i].span = borders[i].level[i].span - (dis[0] - dis[i])
		borders[i].level[i].span = dis[0] - dis[i] + 1
	}

	// update span of borders for levels larger than new node's level
	for i := level; i < zsl.level; i++ {
		borders[i].level[i].span++
	}

	// update backward
	if borders[0] == zsl.header {
		node.backward = nil
	} else {
		node.backward = borders[0]
	}

	if node.level[0].forward == nil {
		zsl.tail = node
	} else {
		node.level[0].forward.backward = node
	}
	zsl.length++
	return node
}

// Delete remove the node with score and value
func (zsl *ZskipList) Delete(score float64, value *Sdshdr) {
	borders := make([]*ZskipListNode, ZSL_MAX_LEVEL)
	p := zsl.header

	for i := zsl.level; i >= 0; i-- {
		for p.level[i].forward != nil && (p.level[i].forward.score < score || (p.level[i].forward.score == score && SdsCmp(p.level[i].forward.value, value) == -1)) {
			p = p.level[i].forward
		}
		borders[i] = p
	}

	node := p.level[0].forward
	if node != nil && node.score == score && SdsCmp(node.value, value) == 0 {
		zsl.DeleteNode(node, borders)
	}
}

// DeleteNode remove node, connect the forward node and backward node
func (zsl *ZskipList) DeleteNode(node *ZskipListNode, borders []*ZskipListNode) {
	// update forward
	for i := 0; i < zsl.level; i++ {
		if borders[i].level[i].forward == node {
			borders[i].level[i].forward = node.level[i].forward
			borders[i].level[i].span += node.level[i].span - 1
		} else {
			borders[i].level[i].span--
		}
	}

	// update backward
	if node.level[0].forward != nil {
		node.level[0].forward.backward = node.backward
	} else {
		zsl.tail = node.backward
	}
	zsl.length--
}

// Update search the real node with curScore and newScore,update its score or delete it and recreate a new node
func (zsl *ZskipList) Update(value *Sdshdr, curScore float64, newScore float64) *ZskipListNode {
	borders := make([]*ZskipListNode, ZSL_MAX_LEVEL)
	p := zsl.header

	for i := zsl.level; i >= 0; i-- {
		for p.level[i].forward != nil && (p.level[i].forward.score < curScore || (p.level[i].forward.score == curScore && SdsCmp(p.level[i].forward.value, value) == -1)) {
			p = p.level[i].forward
		}
		borders[i] = p
	}

	node := p.level[0].forward

	// if the forward is larger or NULL and backward is smaller or NULL
	if (node.level[0].forward == nil || node.level[0].forward.score > newScore) && (node.backward == nil || node.backward.score < newScore) {
		node.score = newScore
		return node
	}

	// it's not sure of the node's position,delete it and recreate
	zsl.DeleteNode(node, borders)
	newNode := zsl.Insert(newScore, node.value)
	return newNode
}
