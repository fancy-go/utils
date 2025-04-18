package seq_util

import "iter"

func Map[F any, T any](seq iter.Seq[F], fn func(F) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

func Map21[F1 any, F2 any, T any](seq iter.Seq2[F1, F2], fn func(F1, F2) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for k, v := range seq {
			if !yield(fn(k, v)) {
				return
			}
		}
	}
}

func Flat[F any](seq iter.Seq[iter.Seq[F]]) iter.Seq[F] {
	return func(yield func(F) bool) {
		for innerSeq := range seq {
			for item := range innerSeq {
				if !yield(item) {
					return
				}
			}
		}
	}
}

func FromSlice[T any](slice []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range slice {
			if !yield(item) {
				return
			}
		}
	}
}

type Pair[K any, V any] struct {
	K K
	V V
}

func FromPairSlice[K, V any](slice []Pair[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, item := range slice {
			if !yield(item.K, item.V) {
				return
			}
		}
	}
}

func Concat[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range seqs {
			for item := range seq {
				if !yield(item) {
					return
				}
			}
		}
	}
}

func Concat2[K, V any](seqs ...iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, seq := range seqs {
			for k, v := range seq {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}
