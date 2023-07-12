package matrix

import (
	"sync"
)

func getRow[T Number](rowA []T, B Matrix[T], wg *sync.WaitGroup) []T {
	defer wg.Done()
	newRow := make([]T, B.Shape()[1])
	for j := 0; j < B.Shape()[1]; j++ {
		var tmp T
		colB := B.GetColumn(j)
		for k := 0; k < len(rowA); k++ {
			tmp += rowA[k] * colB[k]
		}
		newRow[j] = tmp
	}
	return newRow
}

func MultiplyParallel[T Number](A Matrix[T], B Matrix[T]) Matrix[T] {
	if A.Shape()[1] != B.Shape()[0] {
		panic("Incorrect matrix dimensions.")
	}
	var C Matrix[T]
	C = make([][]T, A.Shape()[0])
	for i := 0; i < A.Shape()[0]; i++ {
		C[i] = make([]T, B.Shape()[1])
	}
	var wg sync.WaitGroup

	for i := 0; i < A.Shape()[0]; i++ {
		wg.Add(1)
		go func(i int) {
			C[i] = getRow(A[i], B, &wg)
		}(i)
	}

	wg.Wait()
	return C
}
