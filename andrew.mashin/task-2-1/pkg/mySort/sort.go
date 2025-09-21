package mySort

func BubbleSort(ar []int) {
	for i := 0; i < len(ar); i++ {
		for j := len(ar) - 1; j > i; j-- {
			if ar[j-1] > ar[j] {
				ar[j-1], ar[j] = ar[j], ar[j-1]
			}
		}
	}
}

func SelectSort(ar []int) {
	for i := 0; i < len(ar)-1; i++ {
		min := i
		for j := i + 1; j < len(ar); j++ {
			if ar[min] > ar[j] {
				min = j
			}
		}
		if min != i {
			ar[i], ar[min] = ar[min], ar[i]
		}
	}
}
