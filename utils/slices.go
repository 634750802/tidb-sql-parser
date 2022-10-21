package utils

func Map[T any, R any](slice []T, mapper func(T) R) []R {
	res := make([]R, len(slice))
	for i := range slice {
		res[i] = mapper(slice[i])
	}
	return res
}

func EraseType[T any](slice []T) []any {
	return Map(slice, func(t T) any {
		return t
	})
}
