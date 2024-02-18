package list

func FirstWhere[T any](iter []T, predicate func(T) bool) (*T, bool) {
	for _, t := range iter {
		if predicate(t) {
			return &t, true
		}
	}

	return nil, false
}

func Map[T, K any](iter []T, mapper func(T) K) []K {
	toReturn := make([]K, len(iter))

	for i, m := range iter {
		toReturn[i] = mapper(m)
	}

	return toReturn
}
