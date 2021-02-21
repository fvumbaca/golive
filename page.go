package golive

import (
	"bytes"
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Event struct {
	Kind   string `json:"kind,omitempty"`
	Action string `json:"action"`
	Value  string `json:"value"`
}

type View func(context.Context, interface{}) []byte
type Controller func(context.Context, Event, interface{}) interface{}

type Page struct {
	View
	Controller
}

func NewPage(c Controller, v View) Page {
	return Page{
		Controller: c,
		View:       v,
	}
}

func (p Page) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO: Ensure all bytes are written
	content := p.View(req.Context(), nil)
	err := injectLiveClientShim(bytes.NewReader(content), w, uuid.NewString())
	if err != nil {
		log.Println("html shim: ", err)
	}
}
