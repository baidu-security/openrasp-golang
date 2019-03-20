package common

const (
	initialSize = 8192
)

type DoubleArrayTrie struct {
	check        []int
	base         []int
	used         []bool
	size         int
	allocSize    int
	key          []string
	keySize      int
	length       []int
	value        []int
	progress     int
	nextCheckPos int
	error_       int
}

type Node struct {
	code  int
	depth int
	left  int
	right int
}

func NewDoubleArrayTrie() *DoubleArrayTrie {
	dat := &DoubleArrayTrie{}
	return dat
}

func (dat *DoubleArrayTrie) Clear() {
	dat.check = nil
	dat.base = nil
	dat.used = nil
	dat.allocSize = 0
	dat.size = 0
}

func (dat *DoubleArrayTrie) Build(_key []string, _length, _value []int, _keySize int) int {
	if _keySize > len(_key) || _key == nil {
		return 0
	}
	dat.key = _key
	dat.length = _length
	dat.keySize = _keySize
	dat.value = _value
	dat.progress = 0

	dat.resize(initialSize)

	dat.base[0] = 1
	dat.nextCheckPos = 0
	root_node := &Node{
		left:  0,
		right: dat.keySize,
		depth: 0,
	}
	siblingsSize, siblings := dat.fetch(root_node)
	if siblingsSize != 0 {
		dat.insert(siblings)
	}
	dat.used = nil
	dat.key = nil
	return dat.error_
}

func (dat *DoubleArrayTrie) insert(siblings []*Node) int {
	if dat.error_ < 0 {
		return 0
	}
	begin := 0
	var pos int
	if siblings[0].code+1 > dat.nextCheckPos {
		pos = siblings[0].code
	} else {
		pos = dat.nextCheckPos - 1
	}
	nonzeor_num := 0
	first := 0
	if dat.allocSize <= pos {
		dat.resize(pos + 1)
	}
OUTER:
	for true {
		pos++
		if dat.allocSize <= pos {
			dat.resize(pos + 1)
		}
		if dat.check[pos] != 0 {
			nonzeor_num++
			continue
		} else if first == 0 {
			dat.nextCheckPos = pos
			first = 1
		}
		begin = pos - siblings[0].code
		if dat.allocSize <= (begin + siblings[len(siblings)-1].code) {
			l := float64(1.0) * float64(dat.keySize) / float64(dat.progress+1)
			if 1.05 > l {
				l = 1.05
			}
			dat.resize(int(float64(dat.allocSize) * l))
		}
		if dat.used[begin] {
			continue
		}
		for i := 1; i < len(siblings); i++ {
			if dat.check[begin+siblings[i].code] != 0 {
				continue OUTER
			}
		}
		break
	}
	if float64(1.0)*float64(nonzeor_num)/float64(pos-dat.nextCheckPos+1) >= 0.95 {
		dat.nextCheckPos = pos
	}
	dat.used[begin] = true
	if tmpSize := begin + siblings[len(siblings)-1].code + 1; tmpSize > dat.size {
		dat.size = tmpSize
	}
	for i := 0; i < len(siblings); i++ {
		dat.check[begin+siblings[i].code] = begin
	}
	for i := 0; i < len(siblings); i++ {
		if newSiblingsSize, newSiblings := dat.fetch(siblings[i]); newSiblingsSize == 0 {
			if dat.value != nil {
				dat.base[begin+siblings[i].code] = -dat.value[siblings[i].left] - 1
			} else {
				dat.base[begin+siblings[i].code] = -siblings[i].left - 1
			}
			if dat.value != nil && -dat.value[siblings[i].left] >= 0 {
				dat.error_ = -2
				return 0
			}
			dat.progress++
		} else {
			h := dat.insert(newSiblings)
			dat.base[begin+siblings[i].code] = h
		}
	}
	return begin
}

func (dat *DoubleArrayTrie) fetch(parent *Node) (int, []*Node) {
	if dat.error_ < 0 {
		return 0, nil
	}
	var siblings []*Node
	prev := 0
	for i := parent.left; i < parent.right; i++ {
		if dat.length != nil {
			if dat.length[i] < parent.depth {
				continue
			}
		} else {
			if len(dat.key[i]) < parent.depth {
				continue
			}
		}
		tmp := dat.key[i]
		cur := 0
		if dat.length != nil {
			if dat.length[i] != parent.depth {
				cur = int(tmp[parent.depth]) + 1
			}
		} else {
			if len(tmp) != parent.depth {
				cur = int(tmp[parent.depth]) + 1
			}
		}
		if prev > cur {
			dat.error_ = -3
			return 0, nil
		}
		if cur != prev || len(siblings) == 0 {
			tmp_node := &Node{
				depth: parent.depth + 1,
				code:  cur,
				left:  i,
			}
			if len(siblings) != 0 {
				siblings[len(siblings)-1].right = i
			}
			siblings = append(siblings, tmp_node)
		}
		prev = cur
	}
	if len(siblings) != 0 {
		siblings[len(siblings)-1].right = parent.right
	}
	return len(siblings), siblings
}

func (dat *DoubleArrayTrie) resize(newSize int) int {
	base2 := make([]int, newSize)
	check2 := make([]int, newSize)
	used2 := make([]bool, newSize)
	if dat.allocSize > 0 {
		copy(base2, dat.base)
		copy(check2, dat.check)
		copy(used2, dat.used)
	}
	dat.base = base2
	dat.check = check2
	dat.used = used2
	dat.allocSize = newSize
	return newSize
}

func (dat *DoubleArrayTrie) CommonPrefixSearch(key string) []int {
	return dat.commonPrefixSearch(key, 0, 0, 0)
}

func (dat *DoubleArrayTrie) commonPrefixSearch(key string, pos, length, nodePos int) []int {
	if length <= 0 {
		length = len(key)
	}
	if nodePos <= 0 {
		nodePos = 0
	}
	var result []int
	if dat.base == nil || dat.check == nil {
		return result
	}
	b := dat.base[nodePos]
	var n int
	var p int
	for i := pos; i < length; i++ {
		p = b
		n = dat.base[p]
		if b == dat.check[p] && n < 0 {
			result = append(result, -n-1)
		}
		p = b + int(key[i]) + 1
		if b == dat.check[p] {
			b = dat.base[p]
		} else {
			return result
		}
	}
	p = b
	n = dat.base[p]
	if b == dat.check[p] && n < 0 {
		result = append(result, -n-1)
	}
	return result
}
