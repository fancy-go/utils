package seq_util

import (
	"iter"
)

func Map[F any, T any](seq iter.Seq[F], fn func(F) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

func Map12[F, T1, T2 any](seq iter.Seq[F], fn func(F) (T1, T2)) iter.Seq2[T1, T2] {
	return func(yield func(T1, T2) bool) {
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

func FlatOuterSlice[F any](s []iter.Seq[F]) iter.Seq[F] {
	return func(yield func(F) bool) {
		for _, innerSeq := range s {
			for item := range innerSeq {
				if !yield(item) {
					return
				}
			}
		}
	}
}

func FlatInnerResultSlice[F any](s iter.Seq2[[]F, error]) iter.Seq2[F, error] {
	var zero F
	return func(yield func(F, error) bool) {
		for innerSeq, err := range s {
			if err != nil {
				if !yield(zero, err) {
					return
				}
			}
			for _, item := range innerSeq {
				if !yield(item, nil) {
					return
				}
			}
		}
	}
}

func FlatInnerSlice[F any](s iter.Seq[[]F]) iter.Seq[F] {
	return func(yield func(F) bool) {
		for innerSeq := range s {
			for _, item := range innerSeq {
				if !yield(item) {
					return
				}
			}
		}
	}
}

func Flat12[F1 any, F2 any](seq iter.Seq[iter.Seq2[F1, F2]]) iter.Seq2[F1, F2] {
	return func(yield func(F1, F2) bool) {
		for innerSeq := range seq {
			for k, v := range innerSeq {
				if !yield(k, v) {
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

func Log[T any](seq iter.Seq[T], fn func(int64, T), sampleRate int64) iter.Seq[T] {
	var count int64
	return func(yield func(T) bool) {
		for v := range seq {
			count++
			if count%sampleRate == 0 {
				fn(count, v)
			}
			if !yield(v) {
				return
			}
		}
	}
}

func Log2[T1 any, T2 any](seq iter.Seq2[T1, T2], fn func(int64, T1, T2), sampleRate int64) iter.Seq2[T1, T2] {
	var count int64
	return func(yield func(T1, T2) bool) {
		for k, v := range seq {
			count++
			if count%sampleRate == 0 {
				fn(count, k, v)
			}
			if !yield(k, v) {
				return
			}
		}
	}
}

func Drain[T any](seq iter.Seq[T], fn func(T) bool) iter.Seq[T] {
	if fn == nil {
		fn = func(_ T) bool { return true }
	}
	return func(yield func(T) bool) {
		for v := range seq {
			if !fn(v) {
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}

func Reduce[T any, R any](seq iter.Seq[T], fn func(R, T) R, init R) R {
	for v := range seq {
		init = fn(init, v)
	}
	return init
}

func ReduceUntil[T any, R any](seq iter.Seq[T], fn func(R, T) (R, bool), init R) R {
	var cont bool
	for v := range seq {
		init, cont = fn(init, v)
		if !cont {
			break
		}
	}
	return init
}

func Reduce2[T1 any, T2 any, R any](seq iter.Seq2[T1, T2], fn func(R, T1, T2) R, init R) R {
	for k, v := range seq {
		init = fn(init, k, v)
	}
	return init
}

func ReduceUntil2[T1 any, T2 any, R any](seq iter.Seq2[T1, T2], fn func(R, T1, T2) (R, bool), init R) R {
	var cont bool
	for k, v := range seq {
		init, cont = fn(init, k, v)
		if !cont {
			break
		}
	}
	return init
}

func Count[T any](seq iter.Seq[T]) int {
	count := 0
	for _ = range seq {
		count++
	}
	return count
}

func Count2[T1 any, T2 any](seq iter.Seq2[T1, T2]) int {
	count := 0
	for _, _ = range seq {
		count++
	}
	return count
}

func MapResult[T any, R any](seq iter.Seq[T], fn func(T) (R, error)) iter.Seq2[R, error] {
	return func(yield func(R, error) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

func MapResult2[T1 any, T2 any, R any](seq iter.Seq2[T1, T2], fn func(T1, T2) (R, error)) iter.Seq2[R, error] {
	return func(yield func(R, error) bool) {
		for k, v := range seq {
			if !yield(fn(k, v)) {
				return
			}
		}
	}
}

func UnwrapOrStop[T any](seq iter.Seq2[T, error], stopError *error, stopData *T, cb func(error, T) error) iter.Seq[T] {
	return func(yield func(T) bool) {
		for k, v := range seq {
			if v != nil {
				if cb != nil {
					v = cb(v, k)
					if v == nil {
						continue
					}
				}
				if stopError != nil {
					*stopError = v
				}
				if stopData != nil {
					*stopData = k
				}
				return
			}
			if !yield(k) {
				return
			}
		}
	}

}

type Result[T any] struct {
	Data T
	Err  error
}

func Buffer[T any](seq iter.Seq[T], size int) iter.Seq[T] {
	return func(yield func(T) bool) {
		ch := make(chan T, size)
		breakSignal := make(chan struct{})
		defer close(breakSignal)
		go func() {
			for v := range seq {
				select {
				case ch <- v:
				case <-breakSignal:
					return
				}
			}
			close(ch)
		}()
		for v := range ch {
			if !yield(v) {
				return
			}
		}
	}
}

func Buffer2[F1, F2 any](seq iter.Seq2[F1, F2], size int) iter.Seq2[F1, F2] {
	return func(yield func(F1, F2) bool) {
		ch := make(chan Pair[F1, F2], size)
		breakSignal := make(chan struct{})
		defer close(breakSignal)
		go func() {
			for k, v := range seq {
				select {
				case ch <- Pair[F1, F2]{K: k, V: v}:
				case <-breakSignal:
					return
				}
			}
			close(ch)
		}()
		for item := range ch {
			if !yield(item.K, item.V) {
				return
			}
		}
	}
}

func Batch[T any](seq iter.Seq[T], batchSize int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		var batch []T
		for v := range seq {
			batch = append(batch, v)
			if len(batch) == batchSize {
				if !yield(batch) {
					return
				}
				batch = nil
			}
		}
		if len(batch) > 0 {
			if !yield(batch) {
				return
			}
		}
	}
}

func FromSlice2[F1, F2 any](slice []Pair[F1, F2]) iter.Seq2[F1, F2] {
	return func(yield func(F1, F2) bool) {
		for _, item := range slice {
			if !yield(item.K, item.V) {
				return
			}
		}
	}
}

func Filter[T any](seq iter.Seq[T], fn func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if fn(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Filter2[T1 any, T2 any](seq iter.Seq2[T1, T2], fn func(T1, T2) bool) iter.Seq2[T1, T2] {
	return func(yield func(T1, T2) bool) {
		for k, v := range seq {
			if fn(k, v) {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}
