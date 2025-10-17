package buffet

import (
	"container/heap"
	"fmt"
)

func ChoosingDish(dishPreferenceID int, dishes *Heap) {
	if dishPreferenceID < 1 || dishPreferenceID > dishes.Len() {
		fmt.Println(ErrInvalidNumberDishes)

		return
	}

	for range dishPreferenceID - 1 {
		heap.Pop(dishes)
	}

	fmt.Println((*dishes)[0])
}
