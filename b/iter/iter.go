package iter

import "iter"

type Seq[T any] = iter.Seq[T]

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

func Map[T any, K any](i iter.Seq[T], mapper func(elem T) K) Seq[K] {
	return func(yield func(K) bool) {
		for elem := range i {
			if !yield(mapper(elem)) {
				return
			}
		}
	}
}
