package golive

import (
	"errors"
	"io"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func injectHTMLHashIDs(root *html.Node) error {
	f := []*html.Node{root}
	var i int64 = 0
	for len(f) > 0 {
		f[0].Attr = append(f[0].Attr, html.Attribute{
			Key: "id",
			Val: strconv.FormatInt(i, 10),
		})
		i++
		c := f[0].FirstChild
		for c != nil {
			f = append(f, c)
			c = c.NextSibling
		}
		f = f[1:]
	}
	return nil
}

func ShimInlineHTML(input io.Reader, output io.Writer) error {
	root, err := html.Parse(input)
	if err != nil {
		return err
	}

	err = liveSocketShim(root)
	if err != nil {
		return err
	}

	return html.Render(output, root)
}

func liveSocketShim(root *html.Node) error {
	jsNode, err := html.ParseFragment(strings.NewReader(javascript), nil)
	if err != nil {
		return err
	}

	body := root.FirstChild.FirstChild

	for body != nil && body.Data != "body" {
		body = body.NextSibling
	}

	if body == nil {
		return errors.New("no body found in html tree")
	}

	body.LastChild.NextSibling = jsNode[0]
	body.LastChild = jsNode[0]

	return nil
}

const javascript = `<script type="text/javascript">
var conn;

document.addEventListener("click", function(evnt){
    var elem = document.getElementById(event.target.id)
    if (elem != undefined) {
        var action = elem.getAttribute("golive-onclick")
        if (action != undefined) {
            var value = elem.value
            console.log("Clicking with event:", action, "and value:", elem.value)
            conn.send(JSON.stringify({action, value}))
        }
    }
});

document.addEventListener("change", function(evnt) {
    var elem = document.getElementById(event.target.id)
    if (elem != undefined) {
        var action = elem.getAttribute("golive-onchange")
        if (action != undefined) {
            var value = elem.value
            console.log("Changing with event:", action, "and value:", elem.value)
            conn.send(JSON.stringify({action, value}))
        }
    }
});

window.onload = function () {
    console.log("Golive is live")
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://localhost:3000/live")
        conn.onopen = function() {
            console.log("Connection opening...")
        }
        conn.onclose = function () {
            console.log("Connection Closed")
        }
        conn.onload = function() {
            console.log("on load")
        }
        conn.onmessage = function (evt) {
            // console.log(evt.data)
            var {id, html} = JSON.parse(evt.data)
            if (id == undefined) {
                document.body.innerHTML = html
            } else {
                document.getElementById(id).innerHTML = html
            }
        }
    } else {
        console.log("This browser does not support websockets")
    }
}

function handlerBuilder(eventName) {
    return function(e) {
        console.log("would send event of type " + eventName)
    }
}

</script>`
