package slices

func Map[T, Out any](fn func(T) Out, in []T) []Out {
	out := make([]Out, len(in))
	for i, v := range in {
		out[i] = fn(v)
	}
	return out
}
