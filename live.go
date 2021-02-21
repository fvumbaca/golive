package golive

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

type ViewUpdate struct {
	ID   string `json:"id,omitempty"`
	HTML string `json:"html,omitempty"`
}

type LiveHub struct {
	pageMap   map[string]Page
	statesMap map[string]interface{}
}

func NewHub() LiveHub {
	return LiveHub{
		pageMap:   make(map[string]Page),
		statesMap: make(map[string]interface{}),
	}
}

func (l LiveHub) RegisterPage(path string, page Page) {
	l.pageMap[path] = page
}

func (l LiveHub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stateID := r.URL.Query().Get("state_id")
	pageName := r.URL.Query().Get("page")

	page, pageExists := l.pageMap[pageName]
	if !pageExists {
		log.Printf("page \"%s\" does not exist\n", pageName)
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		var event Event
		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Println("read unmarshal:", err)
			break
		}

		if event.Kind == "unload" {
			log.Printf("Unloading state %s", stateID)
			delete(l.statesMap, stateID)
			break
		}

		if page.Controller != nil {
			l.statesMap[stateID] = page.Controller(r.Context(), event, l.statesMap[stateID])
		}

		updatedPage := page.View(r.Context(), l.statesMap[stateID])
		// TODO: Use a better diff engine here
		// this is only temporary
		updates := []ViewUpdate{
			{HTML: string(updatedPage)},
		}

		err = c.WriteJSON(&updates)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
