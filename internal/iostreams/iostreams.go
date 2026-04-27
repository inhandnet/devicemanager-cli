package iostreams

import (
	"io"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
)

type IOStreams struct {
	In       io.Reader
	Out      io.Writer
	ErrOut   io.Writer
	outIsTTY bool
	termOut  *termenv.Output
	JQExpr   string
}

func (s *IOStreams) IsStdoutTTY() bool {
	return s.outIsTTY
}

func (s *IOStreams) TermOutput() *termenv.Output {
	return s.termOut
}

func System() *IOStreams {
	out := os.Stdout
	isTTY := isTerminal(out)
	tOut := termenv.NewOutput(out)
	if !isTTY {
		tOut = termenv.NewOutput(out, termenv.WithProfile(termenv.Ascii))
	}
	return &IOStreams{
		In:       os.Stdin,
		Out:      out,
		ErrOut:   os.Stderr,
		outIsTTY: isTTY,
		termOut:  tOut,
	}
}

func isTerminal(f *os.File) bool {
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}

// Color helpers

func Bold(s string) string   { return "\033[1m" + s + "\033[0m" }
func Green(s string) string  { return "\033[32m" + s + "\033[0m" }
func Red(s string) string    { return "\033[31m" + s + "\033[0m" }
func Yellow(s string) string { return "\033[33m" + s + "\033[0m" }
func Gray(s string) string   { return "\033[90m" + s + "\033[0m" }
