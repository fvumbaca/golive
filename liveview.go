package golive

import "io"

type LiveView interface {
	Update(Event)
	Render(io.Writer) error
}
