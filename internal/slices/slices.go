package slices

// Map maps a slice to a new slice using the provided function `fn`.
func Map[T, Out any](fn func(T) Out, in []T) []Out {
	out := make([]Out, len(in))
	for i, v := range in {
		out[i] = fn(v)
	}
	return out
}

// Contains returns whether the slice `s` contains the value `v`.
func Contains[T comparable](v T, s []T) bool {
	for _, e := range s {
		if e == v {
			return true
		}
	}
	return false
}

// Unique returns a new slice with all duplicate values removed.
func Unique[S ~[]T, T comparable](s S) S {
	out := make(S, 0, len(s))
	uniq := make(map[T]struct{}, len(s))
	for _, v := range s {
		if _, ok := uniq[v]; ok {
			continue
		}
		uniq[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
