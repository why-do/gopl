package intset

import (
	"bytes"
	"fmt"
)

// 这是一个包含非负整数的集合
// 零值代表空的集合
type IntSet struct {
	words []uint
}

const bitCounts = 32 << (^uint(0) >> 63) // 32位平台这个值就是32，64位平台这个值就是64

// 集合中是否存在非负整数x
func (s *IntSet) Has(x int) bool {
	word, bit := x/bitCounts, uint(x%bitCounts)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// 添加一个数x到集合中
func (s *IntSet) Add(x int) {
	word, bit := x/bitCounts, uint(x%bitCounts)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// 求并集，并保存到s中
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// 以字符串"{1 2 3}"的形式返回集合
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < bitCounts; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", bitCounts*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

// 返回元素个数，查表法
func (s *IntSet) Len() int {
	var pc [256]byte
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}

	var counts int
	for _, word := range s.words {
		counts += int(pc[byte(word>>(0*8))])
		counts += int(pc[byte(word>>(1*8))])
		counts += int(pc[byte(word>>(2*8))])
		counts += int(pc[byte(word>>(3*8))])
		counts += int(pc[byte(word>>(4*8))])
		counts += int(pc[byte(word>>(5*8))])
		counts += int(pc[byte(word>>(6*8))])
		counts += int(pc[byte(word>>(7*8))])
	}
	return counts
}

// 返回元素个数，右移循环算法
func (s *IntSet) Len2() int {
	var count int
	for _, x := range s.words {
		for x != 0 {
			if x & 1 == 1 {
				count++
			}
			x >>= 1
		}
	}
	return count
}

// 返回元素个数，快速法
func (s *IntSet) Len3() int {
	var count int
	for _, x := range s.words {
		for x != 0 {
			x = x & (x - 1)
			count++
		}
	}
	return count
}

// 一次添加多个元素
func (s *IntSet) AddAll(nums ...int) {
	for _, x := range nums {
		s.Add(x)
	}
}

// 移除元素，无论是否在集合中，都把该位置置0
func (s *IntSet) Remove(x int) {
	word, bit := x/bitCounts, uint(x%bitCounts)
	if word < len(s.words) {
		s.words[word] &^= 1 << bit
	}
	// 移除高位全零的元素
	for i := len(s.words)-1; i >=0; i-- {
		if s.words[i] == 0 {
			s.words = s.words[:i]
		} else {
			break
		}
	}
}

// 删除所有元素
func (s *IntSet) Clear() {
	*s = IntSet{}
}

// 返回集合的副本
func (s *IntSet) Copy() *IntSet {
	x := IntSet{words: make([]uint, len(s.words))}
	copy(x.words, s.words)
	return &x
}

// 交集 Intersection
func (s *IntSet) IntersectionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &= tword
		}
	}
}

// 差集 Difference
func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &^= tword
		}
	}
}

// 对称差 SymmetricDifference
func (s *IntSet) SymmetricDifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] ^= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// 返回包含集合元素的 slice，这适合在 range 循环中使用
func (s *IntSet) Elems() []int {
	var ret []int
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < bitCounts; j++ {
			if word&(1<<uint(j)) != 0 {
				ret = append(ret, bitCounts*i+j)
			}
		}
	}
	return ret
}