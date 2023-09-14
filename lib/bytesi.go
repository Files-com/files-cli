package lib

import (
	"fmt"
	"math"
)

func ByteCountSI(b int64) string {
	return ByteCountSIFloat64(float64(b))
}

func ByteCountSIFloat64(b float64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", int64(b))
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	c := b / float64(div)
	if math.IsNaN(c) || c < 0 {
		return "0 B"
	}

	return fmt.Sprintf("%.1f %cB",
		c, "kMGTPE"[exp])
}
