package matrix

import (
	"fmt"
	"math"
	"testing"
)

type multiplyTest[T Number] struct {
	A        Matrix[T]
	B        Matrix[T]
	expected Matrix[T]
}

var multiplyIntTests = []multiplyTest[int]{
	{
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
		Matrix[int]{
			{22},
			{58},
			{94},
		},
	},
}

var multiplyFloatTests = []multiplyTest[float64]{
	{
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
		Matrix[float64]{
			{47.559329, 23.29094, 60.143088},
			{145.553954, 54.870702, 169.371799},
			{80.446899, 45.663206, 117.168233},
		},
	},
}

var multiplyErrorTests = []multiplyTest[int]{
	{
		Matrix[int]{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
		Matrix[int]{{5}},
		Matrix[int]{{}},
	},
}

func TestMultiplyParallel(t *testing.T) {
	for num, test := range multiplyIntTests {
		t.Run(fmt.Sprintf("Test %d", num+1),
			func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Unexpected error: %s", r)
					}
				}()
				result := MultiplyParallel(test.A, test.B)
				ok := true
				for i := 0; i < result.Shape()[0]; i++ {
					for j := 0; j < result.Shape()[1]; j++ {
						ok = (ok && (result[i][j] == test.expected[i][j]))
					}
				}
				if !ok {
					t.Errorf("Result does not equal expected!")
				}
			})
	}
	for num, test := range multiplyFloatTests {
		t.Run(fmt.Sprintf("Test %d", num+1),
			func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Unexpected error: %s", r)
					}
				}()
				result := MultiplyParallel(test.A, test.B)
				ok := true
				for i := 0; i < result.Shape()[0]; i++ {
					for j := 0; j < result.Shape()[1]; j++ {
						ok = (ok && (math.Abs(result[i][j]-test.expected[i][j]) < 1e-6))
					}
				}
				if !ok {
					t.Errorf("Result does not equal expected!")
				}
			})
	}
	for num, test := range multiplyErrorTests {
		t.Run(fmt.Sprintf("Test %d", num+1),
			func(t *testing.T) {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Should have panicked!")
					}
				}()
				MultiplyParallel(test.A, test.B)
			})
	}
}
