package pipeline

import (
	"net/http"
)

type Pipe struct {
	Flow  http.HandlerFunc
	Async bool
}

func processPipeLine(flow []Pipe, w http.ResponseWriter, r *http.Request) {
	for _, f := range flow {
		if f.Async == true {
			go f.Flow(w, r)
		} else {
			f.Flow(w, r)
		}
	}
}
