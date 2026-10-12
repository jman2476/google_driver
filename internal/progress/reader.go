package progress

import (
	"fmt"
	"os"
)

type Tracker struct {
	File     *os.File
	Size     int
	Progress int
	TermTracker
}

func (t *Tracker) SetSize(n int) {
	t.Size = n
}

func (t *Tracker) Read(buf []byte) (n int, err error) {
	n, err = t.File.Read(buf)
	if err != nil {
		return 0, err
	}

	t.Progress += n

	return
}

func (t *Tracker) CalcProgress() (float64, error) {
	if t.Size <= 0 {
		return 0.0, fmt.Errorf("file has size of 0 bytes")
	}
	return float64(t.Progress) / float64(t.Size), nil
}

func (t *Tracker) PrintProgress() error {
	progress, err := t.CalcProgress()
	if err != nil {
		return err
	}

	fmt.Printf("Uploading: %d/%d bytes (%.1f%%)\n", t.Progress, t.Size, progress)
	return nil
}
