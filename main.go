package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			os.Exit(e.Sys().(syscall.WaitStatus).ExitStatus())
		}
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	argv := os.Args[1:]

	// if -json flag presents simply delegate everything to `go test`
	for i := range argv {
		if argv[i] == "-json" {
			cmd := exec.Command("go", append([]string{"test"}, argv...)...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}
	}

	r, w := io.Pipe()
	errc := make(chan error, 1)
	go func() {
		errc <- colorize(r)
	}()
	defer func() {
		_ = r.Close() // signal colorize to stop
		if e := <-errc; e != nil {
			fmt.Fprintf(os.Stderr, "colorize error: %s\n", e)
		}
	}()

	cmd := exec.Command("go", append([]string{"test", "-json"}, argv...)...)
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func colorize(r io.Reader) error {
	stream := json.NewDecoder(r)
	states := map[string]map[string]string{}

	for {
		var ev event
		if err := stream.Decode(&ev); err != nil {
			if err == io.EOF || err == io.ErrClosedPipe {
				break
			}
			return err
		}

		// events without output describe package / test states,
		// so we don't need to print anything we just maintain states map
		if ev.Output == "" {
			switch ev.Action {
			case "run":
				if states[ev.Package] == nil {
					states[ev.Package] = map[string]string{}
				}
			case "pass", "fail", "skip":
				if ev.Test == "" {
					// stop tracking entire package, can only pass or fail
					delete(states, ev.Package)
				} else {
					// stop tracking a single test
					delete(states[ev.Package], ev.Test)
				}
			}
			continue
		}

		color := getOutputColor(ev.Output)

		// output events for the test should be colored the same way, for example:
		// --- FAIL: TestFail (0.00s)
		//     example_test.go:11: failure reason
		if state := getOutputState(ev.Output); state != "" {
			states[ev.Package][ev.Test] = state
		} else if state := states[ev.Package][ev.Test]; state != "" {
			color = getOutputColor(state)
		}

		for _, c := range color {
			fmt.Printf("\033[%dm", c)
		}
		fmt.Print(ev.Output)
		fmt.Print("\033[0m")
	}
	return nil
}

func getOutputState(s string) string {
	switch {
	case strings.HasPrefix(s, stateFail):
		return stateFail
	case strings.HasPrefix(s, statePass):
		return statePass
	case strings.HasPrefix(s, stateSkip):
		return stateSkip
	default:
		return ""
	}
}

func getOutputColor(s string) []int {
	for prefix, color := range colors {
		if strings.HasPrefix(s, prefix) {
			return color
		}
	}
	return nil
}

const (
	stateFail = "--- FAIL"
	statePass = "--- PASS"
	stateSkip = "--- SKIP"
)

var colors = map[string][]int{
	// "=== RUN":   0,
	stateFail:   {termRed},
	statePass:   {termGreen},
	stateSkip:   {termYellow},
	"=== PAUSE": {termDarkGray},
	"=== CONT":  {termDarkGray},
	"PASS":      {termBold, termGreen},
	"ok":        {termBold, termGreen},
	"FAIL":      {termBold, termRed},
	"?":         {termBold, termYellow},
}

const (
	termBold     = 1
	termRed      = 31
	termGreen    = 32
	termYellow   = 33
	termDarkGray = 90
)

// go doc test2json
type event struct {
	Time    time.Time
	Action  string
	Package string
	Test    string
	Elapsed float64
	Output  string
}