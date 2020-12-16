package golive

import (
	"io"

	"golang.org/x/net/html"
)

func renderHTML(in io.Reader, out io.Writer, jsShim bool) error {
	htmlRoot, err := html.Parse(in)
	if err != nil {
		return err
	}

	err = injectHTMLHashIDs(htmlRoot)
	if err != nil {
		return err
	}

	if jsShim {
		err = liveSocketShim(htmlRoot)
		if err != nil {
			return err
		}
	}

	return html.Render(out, htmlRoot)
}
