package main

import (
	"sort"
)

const (
	M = 3
	N = ((M + 1) / 2) - 1
)

type Node struct {
	keys     []int
	parent   *Node
	x        int
	children []*Node
}

func NewNode(key int) *Node {
	return &Node{}
}

func (b Node) New(key int) *Node {
	return &Node{keys: []int{key}}
}

func (n *Node) FindLeaf(key int) (*Node, int) {
	p := 0
	for {
		i := sort.Search(len(n.keys), func(k int) bool {
			return key < n.keys[k]
		})
		if len(n.children) == 0 { // leaf
			return n, p
		}
		n = n.children[i]
		p = i
	}
}

func (n *Node) Query(key int) (int, *Node) {
	for {
		i := sort.Search(len(n.keys), func(k int) bool {
			return n.keys[k] >= key
		})
		if i < len(n.keys) && n.keys[i] == key {
			return i, n
		}
		if len(n.children) == 0 {
			return -1, n
		}
		n = n.children[i]
	}
}

func (n *Node) Predecessor() *Node {
	for len(n.children) > 0 {
		n = n.children[len(n.children)-1]
	}
	return n
}

func (n *Node) Successor() *Node {
	for len(n.children) > 0 {
		n = n.children[0]
	}
	return n
}

func (n *Node) Swap(i, j int, s *Node) {
	t := n.keys[i]
	n.keys[i] = s.keys[j]
	s.keys[j] = t
}

type BTree struct {
	root *Node
}

func NewBTree() *BTree {
	return &BTree{root: &Node{}}
}

func (t *BTree) Insert(key int) {
	node, pos := t.root.FindLeaf(key)
	node.keys = append(node.keys, key)
	sort.Slice(node.keys, func(i, j int) bool { return node.keys[i] < node.keys[j] })
	for len(node.keys) == M {
		x := int(M / 2) // floor(M/2)
		l := &Node{keys: node.keys[:x]}
		r := &Node{keys: node.keys[x+1:]}
		if len(node.children) > 0 {
			l.children = node.children[:x+1]
			for i, c := range l.children {
				c.parent = l
				c.x = i
			}
			r.children = node.children[x+1:]
			for i, c := range r.children {
				c.parent = r
				c.x = i
			}
		}
		p := node.parent
		if p == nil {
			median := &Node{keys: []int{node.keys[x]}}
			l.parent = median
			l.x = 0
			r.parent = median
			r.x = 1
			median.children = []*Node{l, r}
			t.root = median
			return
		}
		p.keys = append(p.keys, node.keys[x])
		sort.Slice(p.keys, func(i, j int) bool { return p.keys[i] < p.keys[j] })
		p.children = append(p.children[:pos], p.children[pos+1:]...)
		l.parent = p
		r.parent = p
		z := make([]*Node, len(p.children)+2)
		switch {
		case pos == 0:
			z[0] = l
			z[0].x = 0
			z[1] = r
			z[1].x = 1
			for i, j := 2, 0; j < len(p.children); j++ {
				z[i] = p.children[j]
				z[i].x = i
				i++
			}
		case pos == len(p.children)-1:
			i := 0
			for ; i < len(p.children); i++ {
				z[i] = p.children[i]
				z[i].x = i
			}
			z[i] = l
			z[i].x = i
			i++
			z[i] = r
			z[i].x = i
		default:
			i := 0
			if len(p.children) > 0 {
				for ; i < pos; i++ {
					z[i] = p.children[i]
					z[i].x = i
				}
			}
			z[i] = l
			z[i].x = i
			i++
			z[i] = r
			z[i].x = i
			i++
			for j := pos; j < len(p.children); j++ {
				z[i] = p.children[j]
				z[i].x = i
				i++
			}
		}
		p.children = z
		if len(p.keys) < M {
			return
		}
		node = p
		pos = node.x
	}
}

func (t *BTree) Delete(key int) bool {
	o, n := t.root.Query(key)
	if o < 0 {
		return false
	}
	if len(n.children) > 0 { // internal node
		d := n.children[o].Predecessor()
		z := len(d.keys) - 1
		n.Swap(o, z, d)
		n, o = d, z
	}
	var p, l, r *Node
	n.keys = append(n.keys[:o], n.keys[o+1:]...)
	for len(n.keys) < N {
		if p = n.parent; p == nil {
			t.root = n
			return true
		}
		x := n.x - 1
		if x >= 0 {
			l = p.children[x]
			if s := len(l.keys); s > N {
				s--
				n.keys = append([]int{p.keys[x]}, n.keys...)
				p.keys[x] = l.keys[s]
				l.keys = l.keys[:s]
				return true
			}
		}
		x = n.x + 1
		if x < len(p.children) {
			r = p.children[x]
			if len(r.keys) > N {
				x--
				n.keys = append(n.keys, p.keys[x])
				p.keys[x] = r.keys[0]
				r.keys = r.keys[1:]
				return true
			}
		}
		pc := make([]*Node, len(p.children)-1)
		switch n.x {
		case 0:
			r = p.children[1]
			n.keys = append(n.keys, p.keys[0])
			n.keys = append(n.keys, r.keys...)
			p.keys = p.keys[1:]
			c := n.children[:]
			i := len(c)
			for j := 0; j < len(r.children); j++ {
				r.children[j].parent = n
				r.children[j].x = i
				c = append(c, r.children[j])
				i++
			}
			n.children = c
			pc[0] = p.children[0]
			pc[0].x = 0
			for i = 1; i < len(p.children)-1; i++ {
				pc[i] = p.children[i+1]
				pc[i].x = i
			}
			l = n
		default:
			l = p.children[n.x-1]
			l.keys = append(l.keys, p.keys[n.x-1])
			l.keys = append(l.keys, n.keys...)
			p.keys = append(p.keys[:n.x-1], p.keys[n.x:]...)
			c := l.children[:] // merge children and parent's list
			i := len(c)
			for j := 0; j < len(n.children); j++ {
				n.children[j].parent = l
				n.children[j].x = i
				c = append(c, n.children[j])
				i++
			}
			l.children = c
			for i = 0; i < n.x; i++ {
				pc[i] = p.children[i]
				pc[i].x = i
			}
			for j := n.x + 1; j < len(p.children); j++ {
				pc[i] = p.children[j]
				pc[i].x = i
				i++
			}
		}
		p.children = pc
		if p.parent == nil {
			// Exception condition,
			// the root node does not have to satisfy the minimum numer of keys,
			// but it must have at least one key
			if len(p.keys) == 0 {
				l.parent = nil
				t.root = l
				return true
			}
		}
		n = p
	}
	return true
}

func (t *BTree) Query(key int) (int, *Node) {
	return t.root.Query(key)
}
