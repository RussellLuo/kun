package trie

import (
	"strings"
)

// Trie is a trie of string keys and interface{} values. Internal nodes
// have nil values so stored nil values cannot be distinguished and are
// excluded from walks. Trie will split keys by slashes
// (e.g. "a/b/c" -> "a", "b", "c").
type Trie struct {
	value    interface{}
	children map[string]*Trie
	ordered  []string
}

// New creates a new *Trie.
func New() *Trie {
	return &Trie{}
}

// Insert inserts the value into the trie at the given key, replacing any
// existing items. It returns true if the put adds a new value, false
// if it replaces an existing value.
//
// Note that internal nodes have nil values so a stored nil value will not
// be distinguishable and will not be included in Walks.
func (trie *Trie) Insert(key string, value interface{}) bool {
	node := trie
	for _, part := range strings.Split(key, "/") {
		child := node.children[part]
		if child == nil {
			if node.children == nil {
				node.children = map[string]*Trie{}
			}
			child = New()
			node.children[part] = child
			node.ordered = append(node.ordered, part)
		}
		node = child
	}
	// does node have an existing value?
	isNewVal := node.value == nil
	node.value = value
	return isNewVal
}

// Get returns the value stored at the given key. Returns nil for internal
// nodes or for nodes with a value of nil.
func (trie *Trie) Get(key string) interface{} {
	node := trie
	for _, part := range strings.Split(key, "/") {
		node = node.children[part]
		if node == nil {
			return nil
		}
	}
	return node.value
}

// Node represents a node of the trie tree.
type Node struct {
	Key   string
	Value interface{}
}

// WalkFunc is used when walking the trie tree. Takes a
// parent and current node, returning if iteration should
// be terminated.
type WalkFunc func(parent *Node, n *Node) bool

// recursiveWalk is used to do a post-order walk of a node
// recursively. Returns true if the walk should be aborted.
func (trie *Trie) recursiveWalk(parent *Node, key string, walk WalkFunc) bool {
	n := &Node{Key: key, Value: trie.value}
	for _, part := range trie.ordered {
		newKey := part
		if key != "" {
			newKey = key + "/" + part
		}

		child := trie.children[part]
		if child.recursiveWalk(n, newKey, walk) {
			return true
		}
	}

	if trie.value != nil {
		if walk(parent, n) {
			return true
		}
	}
	return false
}

// Walk is used to walk the trie tree.
func (trie *Trie) Walk(walk WalkFunc) bool {
	return trie.recursiveWalk(nil, "", walk)
}
