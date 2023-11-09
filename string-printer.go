package toolprinter

import (
	"fmt"
	"strings"
)

type StringPrinter struct {
	TestPrinter
}

func NewStringPrinter() *StringPrinter {
	tp := StringPrinter{}
	tp.stdoutPrinter = tp.print
	tp.statusPrinter = func(a ...interface{}) (n int, err error) { text := fmt.Sprint(a...); n = len(text); return }
	return &tp
}

func (sp *StringPrinter) Text() string {
	return strings.Join(sp.lines, "\n")
}
