package mySort

import (
	"testing"
)

func TestBubbleSort(t *testing.T) {
	ar := []int{5, 35, 12, 24, 21, 98, 6, 10}
	expected := []int{5, 6, 10, 12, 21, 24, 35, 98}

	BubbleSort(ar)

	for inx, elem := range ar {
		if elem != expected[inx] {
			t.Errorf("Bubble sort is not correct. Expect %d, got %d", expected, ar)
			break
		}
	}
}

func TestSelectSort(t *testing.T) {
	ar := []int{5, 35, 12, 24, 21, 98, 6, 10}
	expected := []int{5, 6, 10, 12, 21, 24, 35, 98}

	SelectSort(ar)

	for inx, elem := range ar {
		if elem != expected[inx] {
			t.Errorf("Bubble sort is not correct. Expect %d, got %d", expected, ar)
			break
		}
	}
}
