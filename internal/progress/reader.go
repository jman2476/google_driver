package progress

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

type Tracker struct {
	File     *os.File
	Size     int64
	Progress int
	TermTracker
	window     int
	lastUpload time.Time
}

func (t *Tracker) SetSize(n int64) {
	t.Size = n
}

func (t *Tracker) SetTime(now time.Time) {
	t.lastUpload = now
}

func (t *Tracker) Read(buf []byte) (n int, err error) {
	n, err = t.File.Read(buf)
	if err != nil {
		return
	}

	t.Progress += n

	err = t.PrintProgressBar(n)
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

func (t *Tracker) PrintProgressBar(n int) error {
	timeStamp := time.Now()
	elapsed := timeStamp.Sub(t.lastUpload)
	t.lastUpload = timeStamp

	t.window = int(os.Stdin.Fd())
	width, _, err := term.GetSize(t.window)
	if err != nil {
		return err
	}

	uploadRate := calcRate(elapsed, n)
	// remaining := float64(t.Size-int64(t.Progress)) / uploadRate

	progress, err := t.CalcProgress()
	if err != nil {
		return err
	}

	// t.calcTimeRemaining(uploadRate)

	var message strings.Builder
	message.Write([]byte("Uploading: "))
	message.Write(fmt.Appendf([]byte{}, "%.2f%%", progress*100))

	var bar strings.Builder
	barWidth := width * 1 / 3
	for i := 0; i <= barWidth; i++ {
		if i <= int(float64(barWidth)*progress) {
			bar.Write([]byte("X"))
		} else {
			bar.Write([]byte("-"))
		}
	}

	message.Write(fmt.Appendf(
		[]byte{}, " [%s] %s  %s remaining\r",
		bar.String(),
		formatRate(uploadRate),
		t.calcTimeRemaining(uploadRate),
	))

	fmt.Print(message.String())
	return nil
}

func calcRate(elapse time.Duration, n int) float64 {
	return float64(n) / elapse.Seconds()
}

func formatRate(rate float64) string {
	exp := math.Log10(rate)

	switch int(exp) / 3 {
	case 4:
		return fmt.Sprintf("%.1f TB/s", rate/math.Pow10(12))
	case 3:
		return fmt.Sprintf("%.1f GB/s", rate/math.Pow10(9))
	case 2:
		return fmt.Sprintf("%.1f MB/s", rate/math.Pow10(6))
	case 1:
		return fmt.Sprintf("%.1f KB/s", rate/1000.0)
	default:
		return fmt.Sprintf("%.1f B/s", rate)
	}
}

func FormatSize(size int64) string {
	exp := math.Log10(float64(size))

	switch int(exp) / 3 {
	case 4:
		return fmt.Sprintf("%.1f TB", float64(size)/math.Pow10(12))
	case 3:
		return fmt.Sprintf("%.1f GB", float64(size)/math.Pow10(9))
	case 2:
		return fmt.Sprintf("%.1f MB", float64(size)/math.Pow10(6))
	case 1:
		return fmt.Sprintf("%.1f MB", float64(size)/1000.0)
	default:
		return fmt.Sprintf("%d B", size)
	}
}

func (t *Tracker) calcTimeRemaining(rate float64) string {
	remaining := t.Size - int64(t.Progress)
	// fmt.Printf("%d B, %.2f B/s", remaining, rate)

	timeRemainder := float64(remaining) / rate
	// fmt.Println(timeRemainder, "s")

	if timeRemainder > 1 {
		seconds := int(timeRemainder) % 60
		minutes := int(timeRemainder) / 60
		hours := minutes / 60
		minutes = minutes % 60

		var remaining strings.Builder
		if hours > 0 {
			remaining.Write(
				fmt.Appendf([]byte{}, "%dh", hours),
			)
		}
		if minutes > 0 {
			remaining.Write(
				fmt.Appendf([]byte{}, "%dm", minutes),
			)
		}
		remaining.Write(
			fmt.Appendf([]byte{}, "%d", seconds),
		)

		return remaining.String()
	}

	return fmt.Sprintf("%.2fs", timeRemainder)
}
