package data_structure

import (
	"math/rand"
	"strings"
)

const SkiplistMaxLevel = 32

/*

	/level 2: span=2 | forward\ --------------------------------------> /span=0 | forward\ ----> NULL
	|level 1: span=1 | forward| --------> /span=1 | forward\ ---------> |span=0 | forward| ----> NULL
	|ele                      |           |ele             |            |ele             |
	|score                    |           |score           |            |score           |
	|backward                 | <-------- |backward        | <--------- |backward        |
	\node1                    /           \node2           /            \node3           /
*/

type SkiplistLevel struct {
	forward *SkiplistNode
	// span is number of nodes between current node and node->forward at current level
	span uint32
}

type SkiplistNode struct {
	ele      string
	score    float64
	backward *SkiplistNode
	levels   []SkiplistLevel
}

type Skiplist struct {
	head   *SkiplistNode
	tail   *SkiplistNode
	length uint32
	level  int
}

func (sl *Skiplist) randomLevel() int {
	level := 1
	for rand.Intn(2) == 1 {
		level++
	}
	if level > SkiplistMaxLevel {
		return SkiplistMaxLevel
	}
	return level
}

/*
               /level 31: span=0 | forward\ ----> NULL
               |....                      |
		       |level 2: span=0 | forward | ----> NULL
               |level 1: span=0 | forward | ----> NULL
		       |ele                       |
		       |score                     |
	 NULL <--- |backward                  |
               \head                      /
*/

func (sl *Skiplist) CreateNode(level int, score float64, ele string) *SkiplistNode {
	res := &SkiplistNode{
		ele:      ele,
		score:    score,
		backward: nil,
	}
	res.levels = make([]SkiplistLevel, level)
	return res
}

func CreateSkiplist() *Skiplist {
	sl := Skiplist{
		length: 0,
		level:  1,
	}
	sl.head = sl.CreateNode(SkiplistMaxLevel, 0, "")
	sl.head.backward = nil
	sl.tail = nil
	return &sl
}

/*
Insert a new element to the Skiplist, we allow duplicated scores.
Caller should check if ele is already inserted or not
*/
func (sl *Skiplist) Insert(score float64, ele string) *SkiplistNode {
	// update stores nodes we have to cross to reach the insert position.
	// rank scores the corresponding "rank" of each node in update. Skiplist head's rank == 0.
	update := [SkiplistMaxLevel]*SkiplistNode{}
	rank := [SkiplistMaxLevel]uint32{}
	x := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		if i == sl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		for x.levels[i].forward != nil && (x.levels[i].forward.score < score ||
			(x.levels[i].forward.score == score && strings.Compare(x.levels[i].forward.ele, ele) == -1)) {
			rank[i] += x.levels[i].span
			x = x.levels[i].forward
		}
		update[i] = x
	}

	level := sl.randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.head
			update[i].levels[i].span = sl.length
		}
		sl.level = level
	}

	// create new node and insert
	x = sl.CreateNode(level, score, ele)
	for i := 0; i < level; i++ {
		x.levels[i].forward = update[i].levels[i].forward
		update[i].levels[i].forward = x
		x.levels[i].span = update[i].levels[i].span - (rank[0] - rank[i])
		update[i].levels[i].span = rank[0] - rank[i] + 1
	}

	// increase span for untouched level because we have a new node
	for i := level; i < sl.level; i++ {
		update[i].levels[i].span++
	}

	if update[0] == sl.head {
		x.backward = nil
	} else {
		x.backward = update[0]
	}

	if x.levels[0].forward != nil {
		x.levels[0].forward.backward = x
	} else {
		sl.tail = x
	}

	sl.length++
	return x
}

func (sl *Skiplist) DeleteNode(x *SkiplistNode, update [SkiplistMaxLevel]*SkiplistNode) {
	for i := 0; i < sl.level; i++ {
		if update[i].levels[i].forward == x {
			update[i].levels[i].span += x.levels[i].span - 1
			update[i].levels[i].forward = x.levels[i].forward
		} else {
			update[i].levels[i].span--
		}
	}
	if x.levels[0].forward != nil {
		x.levels[0].forward.backward = x.backward
	} else {
		// x is tail
		sl.tail = x.backward
	}
	for sl.level > 1 && sl.head.levels[sl.level-1].forward == nil {
		sl.level--
	}
	sl.length--
}

func (sl *Skiplist) Delete(score float64, ele string) int {
	update := [SkiplistMaxLevel]*SkiplistNode{}
	x := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil && (x.levels[i].forward.score < score ||
			(x.levels[i].forward.score == score &&
				strings.Compare(x.levels[i].forward.ele, ele) == -1)) {
			x = x.levels[i].forward
		}
		update[i] = x
	}
	x = x.levels[0].forward
	if x != nil && x.score == score && strings.Compare(x.ele, ele) == 0 {
		sl.DeleteNode(x, update)
		return 1
	}
	return 0
}

func (sl *Skiplist) FindFirstInRange(min float64, max float64) *SkiplistNode {
	// TODO: implement detail
	return nil
}
