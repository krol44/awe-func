package aweFunc

// SliceChunk return slice with slices by division on chunkSize
func SliceChunk[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

// SliceReverse return flipped over slice
func SliceReverse[T any](sli []T) []T {
	result := make([]T, len(sli))
	length := len(sli)
	for index := range sli {
		result[(length-1)-index] = sli[index]
	}
	return result
}
