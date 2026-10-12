package progress

import (
	"fmt"
	"os"
)

type Tracker struct {
	File     *os.File
	Size     int64
	Progress int
	TermTracker
}

func (t *Tracker) SetSize(n int64) {
	t.Size = n
}

func (t *Tracker) Read(buf []byte) (n int, err error) {
	n, err = t.File.Read(buf)
	if err != nil {
		return
	}

	t.Progress += n

	err = t.PrintProgress()
	if err != nil {
		return
	}

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

	fmt.Printf("Uploading: %d/%d bytes (%.1f%%)\r", t.Progress, t.Size, progress*100)
	return nil
}
