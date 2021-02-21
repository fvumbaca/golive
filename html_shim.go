package golive

import (
	"errors"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// func injectHTMLHashIDs(root *html.Node) error {
// 	f := []*html.Node{root}
// 	var i int64 = 0
// 	for len(f) > 0 {
// 		f[0].Attr = append(f[0].Attr, html.Attribute{
// 			Key: "id",
// 			Val: strconv.FormatInt(i, 10),
// 		})
// 		i++
// 		c := f[0].FirstChild
// 		for c != nil {
// 			f = append(f, c)
// 			c = c.NextSibling
// 		}
// 		f = f[1:]
// 	}
// 	return nil
// }

func injectLiveClientShim(input io.Reader, output io.Writer, stateID string) error {
	root, err := html.Parse(input)
	if err != nil {
		return err
	}

	js := strings.ReplaceAll(javascript, "${STATE_ID}", stateID)

	jsNodes, err := html.ParseFragment(strings.NewReader(js), nil)
	if err != nil {
		return err
	}

	if len(jsNodes) > 1 {
		panic("Unsupported html shim")
	}

	err = _addJSShim(root, jsNodes[0])
	if err != nil {
		return err
	}

	return html.Render(output, root)
}

func _addJSShim(root *html.Node, jsNode *html.Node) error {

	body := root.FirstChild.FirstChild

	for body != nil && body.Data != "body" {
		body = body.NextSibling
	}

	if body == nil {
		return errors.New("no body found in html tree")
	}

	body.LastChild.NextSibling = jsNode
	body.LastChild = jsNode

	return nil
}

const javascript = `<script type="text/javascript">
// Injected by golive

var conn;

document.addEventListener("click", function(evnt){
	var kind = "click"
    var elem = event.target
    if (elem != undefined) {
        var action = elem.getAttribute("golive-onclick")
        if (action != undefined) {
            var value = elem.value
            console.log("Clicking with event:", action, "and value:", elem.value)
            conn.send(JSON.stringify({kind, action, value}))
        }
    }
});


document.addEventListener("change", function(evnt) {
	var kind = "change"
    var elem = document.getElementById(event.target.id)
    if (elem != undefined) {
        var action = elem.getAttribute("golive-onchange")
        if (action != undefined) {
            var value = elem.value
            console.log("Changing with event:", action, "and value:", elem.value)
            conn.send(JSON.stringify({kind, action, value}))
        }
    }
});

window.addEventListener("beforeunload", function(event) {
	kind = "unload"
	conn.send(JSON.stringify({kind}))
});


window.onload = function () {
    console.log("Golive is live")
    if (window["WebSocket"]) {

        conn = new WebSocket("ws://localhost:3000/live?state_id=${STATE_ID}&page="+window.location.pathname)
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
            var changes = JSON.parse(evt.data)
			changes.forEach(function(i) {
				var {id, html} = i
				console.log(i)
				if (id == undefined) {
					document.body.innerHTML = html
				} else {
					document.getElementById(id).innerHTML = html
				}
			})
            // var {id, html} = JSON.parse(evt.data)

            // if (id == undefined) {
            //     document.body.innerHTML = html
            // } else {
            //     document.getElementById(id).innerHTML = html
            // }
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
