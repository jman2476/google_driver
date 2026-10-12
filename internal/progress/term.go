package progress

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type TermTracker struct {
	term  *term.Terminal
	state *term.State
	// window int
}

func (t *Tracker) SetRawMode() error {
	t.window = int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(t.window)
	if err != nil {
		return fmt.Errorf("unable to start raw mode: %v", err)
	}
	t.state = oldState

	return nil
}

func (t *Tracker) RestoreTerm() error {
	return term.Restore(t.window, t.state)
}
