package matrix

func getRow[T Number](ch chan []T, rowA []T, B Matrix[T]) {
	newRow := make([]T, B.Shape()[1])
	for j := 0; j < B.Shape()[1]; j++ {
		var tmp T
		colB := B.GetColumn(j)
		for k := 0; k < len(rowA); k++ {
			tmp += rowA[k] * colB[k]
		}
		newRow[j] = tmp
	}
	ch <- newRow
}

func Multiply[T Number](A Matrix[T], B Matrix[T]) Matrix[T] {
	if A.Shape()[1] != B.Shape()[0] {
		panic("Incorrect matrix dimensions.")
	}
	var C Matrix[T]
	C = make([][]T, A.Shape()[0])
	for i := 0; i < A.Shape()[0]; i++ {
		C[i] = make([]T, B.Shape()[1])
	}
	for i := 0; i < A.Shape()[0]; i++ {
		for j := 0; j < B.Shape()[1]; j++ {
			C[i][j] = 0
			for k := 0; k < A.Shape()[1]; k++ {
				C[i][j] += A[i][k] * B[k][j]
			}
		}
	}
	return C
}

func MultiplyParallel[T Number](A Matrix[T], B Matrix[T]) Matrix[T] {
	if A.Shape()[1] != B.Shape()[0] {
		panic("Incorrect matrix dimensions.")
	}
	rowChannel := make(chan []T, A.Shape()[1])
	var C Matrix[T]
	C = make([][]T, A.Shape()[0])
	for i := 0; i < A.Shape()[0]; i++ {
		C[i] = make([]T, B.Shape()[1])
	}

	for i := 0; i < A.Shape()[0]; i++ {
		go getRow(rowChannel, A[i], B)
		C[i] = <-rowChannel
	}
	return C
}
