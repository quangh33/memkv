package data_structure

const ZAddInNX = 1 << 1 /* Only add new elements. Don't update already existing elements. */
const ZAddInXX = 1 << 2 /* Only update elements that already exist. Don't add new elements. */

const ZAddOutNop = 1 << 0     /* Operation not performed because of conditionals.*/
const ZAddOutAdded = 1 << 1   /* The element was new and was added. */
const ZAddOutUpdated = 1 << 2 /* The element already existed, score updated. */

type ZSet struct {
	zskiplist *Skiplist
	// map from ele to score
	dict map[string]float64
}

func (zs *ZSet) Add(score float64, ele string, flag int) (int, int) {
	nx := flag & ZAddInNX
	xx := flag & ZAddInXX

	if len(ele) == 0 {
		return 0, ZAddOutNop
	}
	if curScore, exist := zs.dict[ele]; exist {
		if nx != 0 {
			return 1, ZAddOutNop
		}
		if curScore != score {
			znode := zs.zskiplist.UpdateScore(curScore, ele, score)
			zs.dict[ele] = znode.score
			return 1, ZAddOutUpdated
		}
		return 1, ZAddOutNop
	}

	// not exist
	if xx != 0 {
		return 1, ZAddOutNop
	}
	znode := zs.zskiplist.Insert(score, ele)
	zs.dict[ele] = znode.score
	return 1, ZAddOutAdded
}

/*
Return 1 if element existed and was deleted, 0 otherwise
*/
func (zs *ZSet) Del(ele string) int {
	score, exist := zs.dict[ele]
	if !exist {
		return 0
	}
	delete(zs.dict, ele)
	zs.zskiplist.Delete(score, ele)
	return 1
}

func CreateZSet() *ZSet {
	zs := ZSet{
		zskiplist: CreateSkiplist(),
		dict:      map[string]float64{},
	}
	return &zs
}
