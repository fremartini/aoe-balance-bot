package list

func FirstWhere[T any](iter []T, predicate func(T) bool) (*T, bool) {
	for _, t := range iter {
		if predicate(t) {
			return &t, true
		}
	}

	return nil, false
}

func Where[T any](iter []T, predicate func(T) bool) []T {
	toReturn := []T{}

	for _, t := range iter {
		if predicate(t) {
			toReturn = append(toReturn, t)
		}
	}

	return toReturn
}

func Map[T, K any](iter []T, mapper func(T) K) []K {
	toReturn := make([]K, len(iter))

	for i, m := range iter {
		toReturn[i] = mapper(m)
	}

	return toReturn
}

func Contains[T comparable](iter []T, elem T) bool {
	for _, a := range iter {
		if a == elem {
			return true
		}
	}
	return false
}
