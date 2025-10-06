package main

import (
	"fmt"
	"github.com/rekottt/task-2-2/kth"
	"os"
)

func main() {
	var itemCount int
	_, scanErr := fmt.Fscan(os.Stdin, &itemCount)
	if scanErr != nil {
		return
	}
	if itemCount < 1 || itemCount > 10000 {
		return
	}

	values := make([]int, itemCount)
	for i := 0; i < itemCount; i++ {
		_, scanErr = fmt.Fscan(os.Stdin, &values[i])
		if scanErr != nil {
			return
		}
		if values[i] < -10000 || values[i] > 10000 {
			return
		}
	}

	var position int
	_, scanErr = fmt.Fscan(os.Stdin, &position)
	if scanErr != nil {
		return
	}
	if position < 1 || position > itemCount {
		return
	}

	result, err := kth.KthMostPreferred(values, position)
	if err != nil {
		return
	}

	fmt.Println(result)
}
