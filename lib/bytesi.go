package lib

import "fmt"

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
	return fmt.Sprintf("%.1f %cB",
		b/float64(div), "kMGTPE"[exp])
}
