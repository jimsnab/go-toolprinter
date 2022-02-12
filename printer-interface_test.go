package toolprinter

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/jimsnab/go-testutils"
)

var captureStdout = testutils.CaptureStdout
var expectError = testutils.ExpectError
var expectBool = testutils.ExpectBool
var expectPanicError = testutils.ExpectPanicError
var expectString = testutils.ExpectString
var expectValue = testutils.ExpectValue

var Prn ToolPrinter

func SetPrinter(toolPrinter ToolPrinter) ToolPrinter {
	prior := Prn
	Prn = toolPrinter

	return prior
}


type testTerminal struct {
	redirected bool
	badSize    bool
	smallWidth bool
}

func (t *testTerminal) IsTerminal(fd int) bool {
	return !t.redirected
}

func (t *testTerminal) GetSize(fd int) (int, int, error) {
	if t.redirected || t.badSize {
		return 0, 0, fmt.Errorf("file descriptor is not a terminal")
	}
	if t.smallWidth {
		return 10, 25, nil
	} else {
		return 140, 60, nil
	}
}

func TestDefaultPrinterStatus(t *testing.T) {

	xterm = &testTerminal{}
	SetPrinter(NewToolPrinter())

	output := captureStdout(
		t,
		func() {
			Prn.Status("test 12", "3")
			Prn.Status("test 345")
			Prn.Clear()
		},
	)

	expectString(t, "test 123\b\b\b345\b\b\b\b\b\b\b\b        \b\b\b\b\b\b\b\b", output)
}

func TestDefaultPrinterStatusf(t *testing.T) {

	xterm = &testTerminal{}
	SetPrinter(NewToolPrinter())

	output := captureStdout(
		t,
		func() {
			Prn.Statusf("%s", "test 123")
			Prn.Statusf("%s %d", "test", 345)
			Prn.Clear()
		},
	)

	expectString(t, "test 123\b\b\b345\b\b\b\b\b\b\b\b        \b\b\b\b\b\b\b\b", output)
}

func TestDefaultPrinterStatusRedirected(t *testing.T) {

	xterm = &testTerminal{redirected: true}
	SetPrinter(NewToolPrinter())

	output := captureStdout(
		t,
		func() {
			Prn.Status("test 123")
			Prn.Status("test 345")
			Prn.Clear()
		},
	)

	expectString(t, "", output)
}

func TestDefaultPrinterStatusBadTerminal(t *testing.T) {

	xterm = &testTerminal{badSize: true}
	SetPrinter(NewToolPrinter())

	output := captureStdout(
		t,
		func() {
			Prn.Status("test 123")
			Prn.Status("test 345")
			Prn.Clear()
		},
	)

	expectString(t, "", output)
}

func TestDefaultPrinterStatusTruncated(t *testing.T) {

	xterm = &testTerminal{smallWidth: true}
	SetPrinter(NewToolPrinter())

	output := captureStdout(
		t,
		func() {
			Prn.Status("test 123 test 456 test 888")
			Prn.Clear()
		},
	)

	expectString(t, "test 123 \b\b\b\b\b\b\b\b\b         \b\b\b\b\b\b\b\b\b", output)
}

func TestDefaultPrinterStatusErased(t *testing.T) {

	xterm = &testTerminal{}
	prn := NewToolPrinter()

	output := captureStdout(
		t,
		func() {
			prn.Status("test 1234567")
			prn.Status("test!")
			prn.Clear()
		},
	)

	expectString(t, "test 1234567\b\b\b\b\b\b\b\b!       \b\b\b\b\b\b\b\b\b\b\b\b     \b\b\b\b\b", output)

	prn = NewToolPrinter()

	output = captureStdout(
		t,
		func() {
			prn.Status("test!")
			prn.Status("bar")
		},
	)

	expectString(t, "test!\b\b\b\b\bbar  \b\b", output)

	prn = NewToolPrinter()

	output = captureStdout(
		t,
		func() {
			prn.Status("test!")
			prn.Status("testing!")
		},
	)

	expectString(t, "test!\bing!", output)
}

func TestDefaultPrinterChattyStatus(t *testing.T) {

	xterm = &testTerminal{}
	prn := NewToolPrinter()

	output := captureStdout(
		t,
		func() {
			prn.ChattyStatus("test 1", "23")
			prn.ChattyStatus("test 345")
			time.Sleep(1200 * time.Millisecond)
			prn.ChattyStatus("test 678")
			prn.Clear()
		},
	)

	expectString(t, "test 123\b\b\b678\b\b\b\b\b\b\b\b        \b\b\b\b\b\b\b\b", output)
}

func TestDefaultPrinterChattyStatusf(t *testing.T) {

	xterm = &testTerminal{}
	prn := NewToolPrinter()

	output := captureStdout(
		t,
		func() {
			prn.ChattyStatusf("%s %d", "test", 123)
			prn.ChattyStatusf("test 345")
			time.Sleep(1200 * time.Millisecond)
			prn.ChattyStatusf("%s", "test 678")
			prn.Clear()
		},
	)

	expectString(t, "test 123\b\b\b678\b\b\b\b\b\b\b\b        \b\b\b\b\b\b\b\b", output)
}

func TestDefaultPrinterPercentStatus(t *testing.T) {

	xterm = &testTerminal{}
	prn := NewToolPrinter()

	output := captureStdout(
		t,
		func() {
			prn.SetCounterMax(3, "te", "st")
			prn.Count()
			time.Sleep(300 * time.Millisecond)
			prn.Count()
			time.Sleep(300 * time.Millisecond)
			prn.Count()
			prn.Clear()
		},
	)

	expectString(t, "test 1 of 3 33%\b\b\b\b\b\b\b\b\b\b2 of 3 67%\b\b\b\b\b\b\b\b\b\b3 of 3 100%\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b                \b\b\b\b\b\b\b\b\b\b\b\b\b\b\b\b", output)
}

func TestDefaultPrinterPercentStatusQuick(t *testing.T) {

	xterm = &testTerminal{}
	prn := NewToolPrinter()

	output := captureStdout(
		t,
		func() {
			prn.SetCounterMax(100, "test")
			for i := 0; i < 110; i++ {
				prn.Count()
			}
		},
	)

	expectString(t, "test 1 of 100 1%\b\b\b\b\b\b\b\b\b\b\b99 of 100 99%\b\b\b\b\b\b\b\b\b\b\b\b\b100 of 100 100%", output)
}

func TestDefaultPrinterPercentStatusUpdate(t *testing.T) {

	xterm = &testTerminal{}
	prn := NewToolPrinter()

	output := captureStdout(
		t,
		func() {
			prn.SetCounterMax(100, "test")
			prn.Count()
			time.Sleep(300 * time.Millisecond)
			prn.UpdateCountStatus("")
			time.Sleep(300 * time.Millisecond)
			prn.UpdateCountStatus("pa", "ss")
		},
	)

	expectString(t, "test 1 of 100 1% pass", output)
}

func TestDefaultPausedStatus(t *testing.T) {

	xterm = &testTerminal{}
	prn := NewToolPrinter()

	output := captureStdout(
		t,
		func() {
			prn.Status("test 123")
			prn.PauseStatus()
			prn.Status("test 345")
			prn.ResumeStatus()
			prn.Status("test 340")

			// ok to call resume too many times
			prn.ResumeStatus()

			prn.Clear()
		},
	)

	expectString(t, "test 123\b\b\b\b\b\b\b\b        \b\b\b\b\b\b\b\btest 345\b0\b\b\b\b\b\b\b\b        \b\b\b\b\b\b\b\b", output)
}

func TestPrintInStatus(t *testing.T) {

	xterm = &testTerminal{}
	prn := NewToolPrinter()

	output := captureStdout(
		t,
		func() {
			prn.Status("test 123")
			prn.Println("TE", "ST")
			prn.Status("test 345")
		},
	)

	expectString(t, "test 123\b\b\b\b\b\b\b\b        \b\b\b\b\b\b\b\bTEST\ntest 123\b\b\b345", output)
}

func TestPrintlnfInStatus(t *testing.T) {

	xterm = &testTerminal{}
	prn := NewToolPrinter()

	output := captureStdout(
		t,
		func() {
			prn.Status("test 123")
			prn.Printlnf("%s", "TEST")
			prn.Status("test 345")
		},
	)

	expectString(t, "test 123\b\b\b\b\b\b\b\b        \b\b\b\b\b\b\b\bTEST\ntest 123\b\b\b345", output)
}

func TestPrintParts(t *testing.T) {

	xterm = &testTerminal{}
	prn := NewToolPrinter()

	output := captureStdout(
		t,
		func() {
			prn.BeginPrint("")
			prn.ContinuePrint("TES", "T")
			prn.ContinuePrintf("%s", "-ABC")
			prn.EndPrint("")
		},
	)

	expectString(t, "TEST-ABC\n", output)

	prn = NewToolPrinter()

	output = captureStdout(
		t,
		func() {
			prn.BeginPrint("T", "EST")
			prn.EndPrint("-", "ABC")
		},
	)

	expectString(t, "TEST-ABC\n", output)

	prn = NewToolPrinter()

	output = captureStdout(
		t,
		func() {
			prn.EndPrintIfStarted()
		},
	)

	expectString(t, "", output)

	prn = NewToolPrinter()

	output = captureStdout(
		t,
		func() {
			prn.BeginPrint("TEST")
			prn.EndPrintIfStarted()
		},
	)

	expectString(t, "TEST\n", output)
}

func TestPrintInParts(t *testing.T) {
	xterm = &testTerminal{}
	prn := NewToolPrinter()

	prn.BeginPrint("")
	expectPanicError(
		t,
		fmt.Errorf("in a nested print"),
		func() { prn.BeginPrint("Illegal") },
	)
	expectPanicError(
		t,
		fmt.Errorf("in a nested print"),
		func() { prn.Println("Illegal") },
	)

	prn = NewToolPrinter()
	expectPanicError(
		t,
		fmt.Errorf("segmented printing didn't begin yet"),
		func() { prn.ContinuePrint("Illegal") },
	)
	expectPanicError(
		t,
		fmt.Errorf("segmented printing didn't begin yet"),
		func() { prn.EndPrint("Illegal") },
	)
}

func TestPrintDateRange(t *testing.T) {
	xterm = &testTerminal{}
	prn := NewToolPrinter()

	output := captureStdout(
		t,
		func() {
			dt, err := time.Parse("2006-01-02 03:04:05 MST", "2022-01-01 12:00:00 EST")
			if err != nil {
				t.Error(err)
			} else {
				prn.DateRangeStatus(dt, dt, "tes", "ting")
			}
		},
	)

	expectString(t, "testing for 2022-01-01 12:00:00 EST", output)

	prn = NewToolPrinter()

	output = captureStdout(
		t,
		func() {
			start, err := time.Parse("2006-01-02 15:04:05 MST", "2022-01-01 12:00:00 EST")
			if err != nil {
				t.Error(err)
			} else {
				end, err := time.Parse("2006-01-02 15:04:05 MST", "2022-01-01 13:00:00 EST")
				if err != nil {
					t.Error(err)
				} else {
					prn.DateRangeStatus(start, end, "testing")
				}
			}
		},
	)

	expectString(t, "testing between 2022-01-01 12:00:00 EST and 2022-01-01 13:00:00 EST", output)
}

func TestExerciseDefaultTerminal(t *testing.T) {
	_, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
		return
	}

	dt := defaultTerminal{}
	expectBool(t, false, dt.IsTerminal(int(w.Fd())))
	_, _, err = dt.GetSize(int(w.Fd()))
	expectError(t, syscall.ENOTTY, err)
}

func TestVerbose(t *testing.T) {
	prn := NewTestPrinter()

	// off initially
	prn.VerbosePrintln("initial input")
	expectValue(t, 0, len(prn.GetLines()))

	// turn it on
	prn.EnableVerbose(true)
	prn.VerbosePrintln("verbose 1")
	expectValue(t, 1, len(prn.GetLines()))
	expectString(t, "verbose 1", prn.GetLines()[0])

	prn.EnableVerbose(true)	// no op
	prn.VerbosePrintln("verbose 2")
	expectValue(t, 2, len(prn.GetLines()))
	expectString(t, "verbose 2", prn.GetLines()[1])

	prn.EnableVerbose(false)
	prn.VerbosePrintln("verbose 3")
	expectValue(t, 2, len(prn.GetLines()))

	prn.EnableVerbose(false)
	prn.VerbosePrintln("verbose 4")
	expectValue(t, 2, len(prn.GetLines()))
}