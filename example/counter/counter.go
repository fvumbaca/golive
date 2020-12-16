package main

import (
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/fvumabca/golive"
)

func main() {
	ep := golive.NewEndpoint()

	http.Handle("/", ep.ViewHandler(NewCounter))
	http.Handle("/live", ep.LiveHandler(NewCounter))

	log.Println("listening on :3000")
	log.Fatalln(http.ListenAndServe(":3000", nil))
}

func NewCounter() golive.LiveView {
	return &Counter{}
}

type Counter struct {
	Count int
}

func (c *Counter) Update(event golive.Event) {
	switch event.Action() {
	case "inc":
		c.Count++
	case "clear":
		c.Count = 0
	}
}

var counterTpl, _ = template.New("").Parse(`
<h1>My Count {{.Count}}</h1>
<button golive-onclick="inc">Increment</button>
<button golive-onclick="clear">Clear</button>
`)

func (c *Counter) Render(out io.Writer) error {
	return counterTpl.Execute(out, c)
}
