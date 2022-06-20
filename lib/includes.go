package lib

func Includes(item string, includes []string) bool {
	for _, v := range includes {
		if v == item {
			return true
		}
	}

	return false
}
