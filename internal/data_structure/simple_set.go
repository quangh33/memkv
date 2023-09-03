package data_structure

import "math/rand"

// simpleSet simply use native hashmap to store keys
type simpleSet struct {
	key  string
	dict map[string]struct{}
}

func newSimpleSet(key string) Set {
	return &simpleSet{
		key:  key,
		dict: make(map[string]struct{}),
	}
}

func (s *simpleSet) Add(members ...string) int {
	added := 0
	for _, m := range members {
		if _, exist := s.dict[m]; !exist {
			s.dict[m] = struct{}{}
			added++
		}
	}
	return added
}

func (s *simpleSet) Rem(members ...string) int {
	removed := 0
	for _, m := range members {
		if _, exist := s.dict[m]; exist {
			delete(s.dict, m)
			removed++
		}
	}
	return removed
}

func (s *simpleSet) Size() int {
	return len(s.dict)
}

func (s *simpleSet) IsMember(member string) int {
	_, exist := s.dict[member]
	if exist {
		return 1
	}
	return 0
}

func (s *simpleSet) MIsMember(members ...string) []int {
	res := make([]int, len(members))
	for i, m := range members {
		res[i] = s.IsMember(m)
	}
	return res
}

func (s *simpleSet) Members() []string {
	m := make([]string, 0, len(s.dict))
	for k, _ := range s.dict {
		m = append(m, k)
	}
	return m
}

func (s *simpleSet) Pop(count int) []string {
	randKeys := s.Rand(count)
	for _, k := range randKeys {
		delete(s.dict, k)
	}
	return randKeys
}

// TODO: optimize
func (s *simpleSet) Rand(count int) []string {
	temp := make([]string, 0, s.Size())
	for k := range s.dict {
		temp = append(temp, k)
	}

	res := make([]string, count)
	r := make(map[int]struct{})
	for i := 0; i < count; i++ {
		for {
			picked := rand.Intn(s.Size())
			if _, ok := r[picked]; !ok {
				res[i] = temp[picked]
				r[picked] = struct{}{}
				break
			}
		}
	}
	return res
}
