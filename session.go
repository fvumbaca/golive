package golive

import (
	"bytes"
	"fmt"

	"github.com/gorilla/websocket"
)

type session struct {
	ID      string
	conn    *websocket.Conn
	Handler LiveView
}

func (s *session) liveLoop() {
	defer func() {
		s.conn.Close()
	}()

	for {
		var event BrowserEvent
		err := s.conn.ReadJSON(&event)
		if err != nil {
			// TODO: handle this properly
			fmt.Println("error:", err)
			continue
		}

		s.Handler.Update(event)

		var b bytes.Buffer
		err = s.Handler.Render(&b)
		if err != nil {
			// TODO: handle this properly
			fmt.Println("error:", err)
			continue
		}

		var bb bytes.Buffer
		renderHTML(&b, &bb, false)

		err = s.conn.WriteJSON(&elementUpdate{
			HTML: bb.String(),
		})

		// pl, err := json.Marshal(&elementUpdate{
		// 	HTML: bb.String(),
		// })
		if err != nil {
			// TODO: handle this properly
			fmt.Println("error:", err)
			continue
		}

		// s.send <- pl
	}
}

type Event interface {
	Action() string
	Value() string
}

type BrowserEvent struct {
	Actn string `json:"action"`
	Val  string `json:"value"`
}

func (b BrowserEvent) Action() string {
	return b.Actn
}

func (b BrowserEvent) Value() string {
	return b.Val
}

type elementUpdate struct {
	ID   string `json:"id,omitempty"`
	HTML string `json:"html"`
}
