package helpers

func Merge(l1, l2 []string) []string {

	if len(l2) > len(l1) {
		return Merge(l2, l1)
	}

	for _, s := range l2 {
		l1 = append(l1, s)
	}

	return l1
}
