package list

func FirstWhere[T any](iter []T, predicate func(T) bool) (*T, int, bool) {
	for i, t := range iter {
		if predicate(t) {
			return &t, i, true
		}
	}

	return nil, 0, false
}
