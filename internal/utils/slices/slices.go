package slices

func Map[T any, M any](a []T, f func(T) M) []M {
	n := make([]M, len(a))
	for i, e := range a {
		n[i] = f(e)
	}
	return n
}

func Contains[T comparable](a []T, elem T) bool {
	for i := range a {
		if a[i] == elem {
			return true
		}
	}
	return false
}

func Filter[T comparable](a []T, elem T) []T {
	result := make([]T, 0, len(a))
	for i := range a {
		if a[i] == elem {
			continue
		}
		result = append(result, a[i])
	}
	return result
}
