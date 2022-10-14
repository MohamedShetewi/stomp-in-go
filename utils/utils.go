package utils

import (
	"errors"
	"golang.org/x/exp/constraints"
	"math"
)

type Comparable interface {
	constraints.Integer | constraints.Float
}

func RemoveIndex[T any](slice []*T, idx int) error {
	if idx < 0 || idx >= len(slice) {
		return errors.New("invalid index")
	}
	slice[idx] = slice[len(slice)-1]
	slice = slice[:len(slice)-1]

	return nil
}

func Max[T Comparable](x, y T) T {
	return T(math.Max(float64(x), float64(y)))
}
