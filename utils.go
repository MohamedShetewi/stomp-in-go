package main

func removeIndex[T comparable](slice *[]*T, idx int) {

	(*slice)[idx] = (*slice)[len(*slice)-1]
	*slice = (*slice)[:len(*slice)-1]
}
