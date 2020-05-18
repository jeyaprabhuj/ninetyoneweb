package handler

import (
	"github.com/jeyaprabhuj/forty/structures/trie"
)

var route *trie.Trie

func init() {
	route = trie.CreateTrie()
}
