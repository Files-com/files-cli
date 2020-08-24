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
	return time.Time(*s).Format(time.RFC3339)
}

func TimeVarP(f *pflag.FlagSet, p *time.Time, name, shorthand string) {
	f.VarP(newTimeValue(p), name, shorthand, "2012-11-01T22:08:41+00:00")
}
