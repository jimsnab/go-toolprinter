package toolprinter

import (
	"fmt"
	"strings"
)

type TestPrinter struct {
	defaultPrinter
	partial bool
	lines   []string
}

func NewTestPrinter() *TestPrinter {
	tp := TestPrinter{}
	tp.stdoutPrinter = tp.print
	tp.statusPrinter = func(a ...interface{}) (n int, err error) { text := fmt.Sprint(a...); n = len(text); return }
	tp.lines = make([]string, 0)
	return &tp
}

func (tp *TestPrinter) print(args ...interface{}) (n int, err error) {
	text := fmt.Sprint(args...)
	n = len(text)

	parts := strings.Split(text, "\n")

	end := len(parts) - 1
	pos := 0
	if tp.partial {
		tp.lines[len(tp.lines)-1] += parts[0]
		if pos == end {
			return
		}
		pos++
		tp.partial = false
	}
	for pos < end {
		tp.lines = append(tp.lines, parts[pos])
		pos++
	}
	if len(parts[pos]) > 0 {
		tp.lines = append(tp.lines, parts[pos])
		tp.partial = true
	}
	return
}

func (tp *TestPrinter) GetStatusText() string {
	return tp.lastStatusText
}

func (tp *TestPrinter) GetLines() []string {
	return tp.lines
}
