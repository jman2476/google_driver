package progress

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
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

	err = t.PrintProgressBar()
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

func (t *Tracker) PrintProgressBar() error {
	t.window = int(os.Stdin.Fd())
	width, _, err := term.GetSize(t.window)
	if err != nil {
		return err
	}

	progress, err := t.CalcProgress()
	if err != nil {
		return err
	}

	var message strings.Builder
	message.Write([]byte("Uploading: "))
	message.Write(fmt.Appendf([]byte{}, "%.2f%%", progress*100))

	var bar strings.Builder
	barWidth := width * 1 / 2
	for i := 0; i <= barWidth; i++ {
		if i <= int(float64(barWidth)*progress) {
			bar.Write([]byte("X"))
		} else {
			bar.Write([]byte("-"))
		}
	}

	message.Write(fmt.Appendf([]byte{}, " [%s]\r", bar.String()))

	fmt.Print(message.String())
	return nil
}
