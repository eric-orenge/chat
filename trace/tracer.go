package trace

import (
	"fmt"
	"io"
)

//Tracing is a practice by which we log or print key steps in the flow of a program to
// make what is going on under the covers visible
type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}
func Off() Tracer {
	return &nilTracer{}
}

type Tracer interface {
	Trace(...interface{})
}
type tracer struct {
	out io.Writer
}

func (t tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
