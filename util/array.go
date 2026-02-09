package util

func ContainsFunc[T comparable](list []T) func(T) bool {
	set := make(map[T]struct{}, len(list))
	for _, v := range list {
		set[v] = struct{}{}
	}
	return func(e T) bool {
		_, ok := set[e]
		return ok
	}
}
