package handler

import (
	"fmt"
	"github.com/jeyaprabhuj/forty/structures/trie"
	"net/http"
	"strings"
)

type Mux struct {
	root string
}

func NewMux(root string) *Mux { return &Mux{root: root} }

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f := fetchRoute(r.URL.Path, r.Method)
	if f != nil {
		f(w, r)
	} else {
		http.Error(w, "Resource not found", http.StatusNotFound)
	}
}

func (m *Mux) GET(path string, handler http.HandlerFunc) {
	addRoute(m.GetRootRoute()+path, http.MethodGet, handler)
}

func (m *Mux) POST(path string, handler http.HandlerFunc) {
	addRoute(m.GetRootRoute()+path, http.MethodPost, handler)
}

func (m *Mux) Handle(path string, method string, handler http.HandlerFunc) {
	addRoute(m.GetRootRoute()+path, strings.ToUpper(method), handler)
}

func (m *Mux) GetRootRoute() string {
	return fmt.Sprintf("/%s/", m.root)
}

func addRoute(path string, method string, handler http.HandlerFunc) {
	mapkeys := extractKeys(path)

	traverseTrie := route

	var i int

	var node *trie.TrieNode
	if len(mapkeys) == 1 {
		node = route.Insert(mapkeys[0])
		node.AddAttribute(method, handler)
	} else {
		for i = 0; i < len(mapkeys); i++ {
			node = traverseTrie.GetNode(mapkeys[i])
			if node == nil {
				node = traverseTrie.Insert(mapkeys[i])
			}

			if i < len(mapkeys)-1 {
				attribute := node.GetAttributeValue(mapkeys[i])
				if attribute != nil {
					temp, isTrie := attribute.(*trie.Trie)
					if isTrie {
						traverseTrie = temp
					}
				} else {
					childRoute := trie.CreateTrie()
					node.AddAttribute(mapkeys[i], childRoute)
					traverseTrie = childRoute
				}

			} else if i == len(mapkeys)-1 {
				node.AddAttribute(method, handler)
			}

		}
	}
}

func fetchRoute(path string, method string) http.HandlerFunc {
	mapkeys := extractKeys(path)

	var i int
	traverseTrie := route
	var node *trie.TrieNode
	if len(mapkeys) == 1 {
		node = traverseTrie.GetNode(mapkeys[0])
		if node != nil {
			return node.GetAttributeValue(mapkeys[0]).(http.HandlerFunc)
		}
	} else {
		for i = 0; i < len(mapkeys); i++ {

			node = traverseTrie.GetNode(mapkeys[i])

			if node != nil {
				if i < len(mapkeys)-1 {
					traverseTrie = node.GetAttributeValue(mapkeys[i]).(*trie.Trie)
				} else if i == len(mapkeys)-1 {
					return node.GetAttributeValue(method).(http.HandlerFunc)
				}
			} else {
				break
			}
		}
	}
	return nil
}

func extractKeys(path string) []string {
	keys := strings.Split(path, "/")
	mapkeys := make([]string, 0)
	for _, mapkey := range keys {
		if len(mapkey) == 0 {
		} else {
			mapkeys = append(mapkeys, mapkey)
		}
	}

	return mapkeys
}

func (m *Mux) Print() {
	route.Print()
}
