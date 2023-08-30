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

/*
Returns the 0-based rank of the object or -1 if the object does not exist.
If reverse is false, rank is computed considering as first element the one
with the lowest score. If reverse is true, rank is computed considering as element with rank 0 the
one with the highest score.
*/
func (zs *ZSet) GetRank(ele string, reverse bool) (rank int64, score float64) {
	setSize := zs.zskiplist.length
	score, exist := zs.dict[ele]
	if !exist {
		return -1, 0
	}
	rank = int64(zs.zskiplist.GetRank(score, ele))
	if reverse {
		rank = int64(setSize) - rank
	} else {
		rank--
	}
	return rank, score
}

func (zs *ZSet) GetScore(ele string) (int, float64) {
	score, exist := zs.dict[ele]
	if !exist {
		return -1, 0
	}
	return 0, score
}

func (zs *ZSet) Len() int {
	return len(zs.dict)
}

func CreateZSet() *ZSet {
	zs := ZSet{
		zskiplist: CreateSkiplist(),
		dict:      map[string]float64{},
	}
	return &zs
}
