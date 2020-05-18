package pipeline

import (
	"net/http"
)

type Chain struct {
	beforeFlow []Pipe
	afterFlow  []Pipe
	fitFlow    http.HandlerFunc
	gate       http.HandlerFunc
}

func (c Chain) BeforeHandlers(handles ...Pipe) Chain {

	newChain := make([]Pipe, len(handles)+len(c.beforeFlow))
	newChain = append(([]Pipe)(nil), c.beforeFlow...)
	newChain = append(newChain, handles...)
	c.beforeFlow = newChain
	return c
}

func (c Chain) AfterHandlers(handles ...Pipe) Chain {
	newChain := make([]Pipe, len(handles)+len(c.afterFlow))
	newChain = append(([]Pipe)(nil), c.afterFlow...)
	newChain = append(newChain, handles...)
	c.afterFlow = newChain
	return c
}

func NewChain() Chain {
	return Chain{}
}

func (c Chain) Process() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		processPipeLine(c.beforeFlow, w, r)
		gatedFlow(c, w, r)
		processPipeLine(c.afterFlow, w, r)
	}
}

func (c Chain) FitFlow(handler http.HandlerFunc) Chain {
	c.fitFlow = handler
	return c
}

func (c Chain) Gate(handler http.HandlerFunc) Chain {
	c.gate = handler
	return c
}

func flow(handler http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	if handler == nil {
		// do nothing
	} else {
		handler(w, r)
	}
}

func gatedFlow(c Chain, w http.ResponseWriter, r *http.Request) {
	if c.gate == nil {
		flow(c.fitFlow, w, r)
	} else {
		c.gate(w, r)
		flow(c.fitFlow, w, r)
	}
}
