package matrix

import (
	"fmt"
	"testing"
)

type multiplyTest[T Number] struct {
	A   Matrix[T]
	B   Matrix[T]
	err bool
}

var multiplyIntTests = []multiplyTest[int]{
	multiplyTest[int]{
		Matrix[int]{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
		Matrix[int]{
			{5},
			{4},
			{3},
		},
		false,
	},
	multiplyTest[int]{
		Matrix[int]{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
		Matrix[int]{{5}},
		true,
	},
}

var multiplyFloatTests = []multiplyTest[float64]{
	multiplyTest[float64]{
		Matrix[float64]{
			{0.733, 5.354, 3.989},
			{4.123, 5.963, 12.352},
			{4.554, 8.643, 4.885},
		},
		Matrix[float64]{
			{5.439, 3.420, 9.999},
			{1.111, 2.222, 3.334},
			{9.432, 2.228, 8.765},
		},
		false,
	},
}

func TestMultiplyParallel(t *testing.T) {
	for num, test := range multiplyIntTests {
		defer func() {
			if r := recover(); !test.err && r == nil {
				t.Errorf(fmt.Sprintf("Test %d failed! Should have panicked!", num))
			}
		}()
		result := MultiplyParallel(test.A, test.B)
		expected := Multiply(test.A, test.B)
		ok := true
		for i := 0; i < result.Shape()[0]; i++ {
			for j := 0; j < result.Shape()[1]; j++ {
				ok = (ok && (result[i][j] == expected[i][j]))
			}
		}
		if !ok {
			t.Errorf(fmt.Sprintf("Test %d failed! Result matrices are not the same!", num))
		}
	}
	for num, test := range multiplyFloatTests {
		defer func() {
			if r := recover(); !test.err && r == nil {
				t.Errorf(fmt.Sprintf("Test %d failed! Should have panicked!", num))
			}
		}()
		result := MultiplyParallel(test.A, test.B)
		expected := Multiply(test.A, test.B)
		ok := true
		for i := 0; i < result.Shape()[0]; i++ {
			for j := 0; j < result.Shape()[1]; j++ {
				ok = (ok && (result[i][j]-expected[i][j] < 1e-6))
			}
		}
		if !ok {
			t.Errorf(fmt.Sprintf("Test %d failed! Result matrices are not the same!", num))
		}
	}
}
