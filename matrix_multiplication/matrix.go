package matrix

type Number interface {
	int | int64 | float32 | float64
}

type Matrix[T Number] [][]T

func (m Matrix[T]) Shape() []int {
	if len(m) == 0 {
		return []int{0, 0}
	}
	return []int{len(m), len(m[0])}
}

func (m Matrix[T]) GetColumn(columnIndex int) []T {
	numRows := len(m)
	column := make([]T, numRows)

	for i := 0; i < numRows; i++ {
		column[i] = m[i][columnIndex]
	}

	return column
}
