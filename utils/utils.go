package utils

import "errors"

func RemoveIndex[T any](slice []*T, idx int) error {
	if idx < 0 || idx >= len(slice) {
		return errors.New("invalid index")
	}
	slice[idx] = slice[len(slice)-1]
	slice = slice[:len(slice)-1]

	return nil
}
