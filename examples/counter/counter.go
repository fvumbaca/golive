package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/fvumabca/golive"
)

func main() {
	h := golive.NewHub()
	http.Handle("/live", h)

	p := golive.NewPage(controller, view)
	h.RegisterPage("/", p)
	http.Handle("/", p)

	log.Println("listening on :3000")
	log.Fatalln(http.ListenAndServe(":3000", nil))
}

func controller(ctx context.Context, e golive.Event, state interface{}) interface{} {
	if v, ok := state.(int); ok {
		switch e.Action {
		case "inc":
			return v + 1
		case "clear":
			return 0
		}
	}
	return 1
}

func view(ctx context.Context, info interface{}) []byte {
	if v, ok := info.(int); ok {
		return []byte(fmt.Sprintf(`<h1>Count: %d</h1>
		<button golive-onclick="inc">Increment</button>
		<button golive-onclick="clear">Clear</button>`, v))
	}

	return []byte(`<button golive-onclick="inc">Start Counting</button>`)
}
