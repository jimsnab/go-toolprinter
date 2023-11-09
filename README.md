# Tool Printer

CLI tools typically need to show progress status temporarily in the
midst of regular console output. This *Tool Printer* interface was created
as a way to separate "progress bar" printing from "final output".

The interface is:

```go
type ToolPrinter interface {
	Status(text ...interface{})
	Statusf(format string, args ...interface{})
	Clear()
	ChattyStatus(text ...interface{})
	ChattyStatusf(format string, args ...interface{})
	SetCounterMax(max int, text ...interface{})
	UpdateCountStatus(extraStatusText ...interface{})
	Count()
	PauseStatus()
	ResumeStatus()
	Println(text ...interface{})
	Printlnf(format string, args ...interface{})
	BeginPrint(text ...interface{})
	ContinuePrint(text ...interface{})
	ContinuePrintf(format string, args ...interface{})
	EndPrint(text ...interface{})
	EndPrintIfStarted()
	DateRangeStatus(from time.Time, to time.Time, purpose ...interface{})
	VerbosePrintln(text ...interface{})
	VerbosePrintlnf(format string, args ...interface{})
	EnableVerbose(enabled bool)
}
```

A ttl implementation is provided via `toolprinter.NewToolPrinter()`. This
implmentation will print status neatly to a ttl, while suppressing it if
`stdout` is redirected to a non-ttl device such as a file.

Example Use:

```go
package main

import (
	"github.com/jimsnab/go-toolprinter"
)

func main() {
	prn := toolprinter.NewToolPrinter()

	for i := 0; i < 100; i++ {
		prn.Statusf("Processing i at %d", i)
		if i%9 == 0 {
			prn.Printlnf("i is divisble by 9 at %d", i)
		}
	}

	prn.Clear()
}
```

[try it](https://go.dev/play/p/_2CoJnMNG4-)

## Status
`Status()` and `Statusf()` print a single line status message, without a new line. `Clear()` erases the status message.
`ChattyStatus()` and `ChattyStatusf()` print status once every 250 ms; for status within high speed processing.

`PauseStatus()` clears the status message (if any), typically because regular printing is desired. `ResumeStatus()`
re-prints the latest status message.

Date range status can be printed with `DateRangeStatus()`.

## Progress Counter

`SetCounterMax()` starts a progress status line, such as:

```go
  prn.SetCounterMax("Reading", 1000)
```

produces the following status message:

```text
Reading 1 of 1000 0%
```

Calling `Count()` increments the counter, and the status message with progress percentage is updated. It won't go beyond
100%.

Suppose your program is processing three items, and the processing of each item takes several minutes, and
the processing of one item can be broken down into individual named steps. For example, let's say they are
"querying", "sorting", "analyzing", "generating output". The program can call `UpdateCountStatus()` to
print querying, sorting, etc., to achieve status such as:

```text
Processing 1 of 3 33% querying
```

then

```text
Processing 1 of 3 33% sorting
```

then

```text
Processing 1 of 3 33% analyzing
```

and so on.

## Printing Regular Output

To use `fmt`, `log` or something similar that produces stdout output, first call `PauseStatus()`, then print to stdout, then
call `ResumeStatus()`.

When you can avoid `fmt` and `log`, use the following `ToolPrinter` functions to print regular output:

* `Println()` and `Printlnf()` - prints a single line of regular output. Any status message is cleared before printing the line,
  and restored after printing the line.
* `BeginPrint()`, `ContinuePrint()`, `ContinuePrintf()` and `EndPrint()` - prints regular output in multiple print calls.
  `BeginPrint()` erases the status message, if any, and `ContinuePrint()`/`ContinuePrintf()` print text exactly as it is
  specified and does not add a line break. `EndPrint()` prints a line break, then reprints the status message, if any.
* `EndPrintIfStarted()` will perform `EndPrint()` only if printing was started with `BeginPrint()` but not ended.

## Verbose Printing

Tools often have different levels of printing. This library supports two - regular and verbose.

* `VerbosePrintln()` and `VerbosePrintlnf()` only print when verbose is enabled.
* `EnableVerbose()` controls whether verbose printing is on or off

## String Printer

Print to a string buffer using `NewStringPrinter()`. Get the content via `sp.Text()` where `sp` is your string printer instance.

## Test Printer

For unit tests, a test version of the `ToolPrinter` interface is provided. Instead of printing, it captures output lines
into an array, giving unit tests an easy way to check for expected output. Status text is not captured, except for the
current status line that a unit test can examine.

Obtain a test printer instance with `NewTestPrinter()`.

Additional methods of a test printer are:

* `GetLines()` returns a string array of printed lines
* `String()` returns the line-parsed input as a single string
* `GetStatusText()` returns the current status text
