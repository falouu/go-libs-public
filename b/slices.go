package b

func Map[S any, D any](src []S, mapper func(input S) D) []D {
	dst := []D{}
	for _, elem := range src {
		dst = append(dst, mapper(elem))
	}
	return dst
}
