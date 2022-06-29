package lib

import (
	"time"

	"github.com/spf13/pflag"
)

type timeValue time.Time

func newTimeValue(p *time.Time) *timeValue {
	return (*timeValue)(p)
}

func (s *timeValue) Set(val string) error {
	pt, err := time.Parse(time.RFC3339, val)

	if err != nil {
		return err
	}

	*s = timeValue(pt)
	return nil
}
func (s *timeValue) Type() string {
	return "time"
}

func (s *timeValue) String() string {
	if s == nil {
		return ""
	}
	return time.Time(*s).Format(time.RFC3339)
}

func TimeVar(f *pflag.FlagSet, p *time.Time, name string) {
	f.Var(newTimeValue(p), name, "2012-11-01T22:08:41+00:00")
}
