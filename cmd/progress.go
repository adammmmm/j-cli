package cmd

import (
	"fmt"
	"os"
	"sync/atomic"
)

type Progress struct {
	total   int64
	current int64
}

func NewProgress(total int) *Progress {
	return &Progress{total: int64(total)}
}

func (p *Progress) Increment() {
	cur := atomic.AddInt64(&p.current, 1)
	percent := float64(cur) / float64(p.total) * 100
	fmt.Fprintf(os.Stderr, "Processed %d/%d devices (%.2f%%)\r", cur, p.total, percent)

	if cur == p.total {
		fmt.Fprintln(os.Stderr)
	}
}
