package golive

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Endpoint struct {
	sessions map[string]*session
	mux      sync.Mutex
}

func NewEndpoint() *Endpoint {
	return &Endpoint{
		sessions: make(map[string]*session),
	}
}

func (e *Endpoint) ViewHandler(lv func() LiveView) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var b bytes.Buffer
		lv().Render(&b)          // TODO: handle error
		renderHTML(&b, rw, true) // TODO: handle error
	})
}

func (e *Endpoint) LiveHandler(lv func() LiveView) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// TODO: replace with upgrader
		conn, err := websocket.Upgrade(rw, r, nil, 1024, 1024) // TODO: refine buffer sizes
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		s := session{
			ID:      fmt.Sprint(rand.Int()), // TODO: better session id generation
			conn:    conn,
			Handler: lv(),
		}
		e.mux.Lock()
		e.sessions[s.ID] = &s // TODO: check for existing session
		e.mux.Unlock()

		go s.liveLoop()
	})
}
