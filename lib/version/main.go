package version

import (
	"fmt"
	"strconv"
	"strings"
)

func New(v string) (Version, error) {
	parts := strings.Split(strings.TrimSpace(strings.Replace(v, "v", "", 1)), ".")
	var err error
	var major int64
	if len(parts) >= 1 {
		major, err = strconv.ParseInt(parts[0], 10, 32)
		if err != nil {
			return Version{}, err
		}
	}
	var minor int64
	if len(parts) >= 2 {
		minor, err = strconv.ParseInt(parts[1], 10, 32)
		if err != nil {
			return Version{}, err
		}
	}
	var patch int64
	if len(parts) == 3 {
		patch, err = strconv.ParseInt(parts[2], 10, 32)
		if err != nil {
			return Version{}, err
		}
	}

	return Version{major, minor, patch}, nil
}

type Version struct {
	Major int64
	Minor int64
	Patch int64
}

func (v Version) Equal(o Version) bool {
	if o.Major != v.Major {
		return false
	}
	if o.Minor != v.Minor {
		return false
	}
	if o.Patch != v.Patch {
		return false
	}

	return true
}

func (v Version) Greater(o Version) bool {
	if o.Major > v.Major {
		return true
	}
	if o.Major < v.Major {
		return false
	}
	if o.Minor > v.Minor {
		return true
	}
	if o.Minor < v.Minor {
		return false
	}
	if o.Patch > v.Patch {
		return true
	}

	return false
}

func (v Version) String() string {
	return fmt.Sprintf("%v.%v.%v", v.Major, v.Minor, v.Patch)
}
