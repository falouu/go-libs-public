package iter

import "iter"

func FindAny[T any](i iter.Seq[T], condition func(elem T) bool) *T {
	for elem := range i {
		if condition(elem) {
			return &elem
		}
	}
	return nil
}

func AnyMatch[T any](i iter.Seq[T], condition func(elem T) bool) bool {
	return FindAny(i, condition) != nil
}
