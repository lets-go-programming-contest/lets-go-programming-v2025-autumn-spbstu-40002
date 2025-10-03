package main

import (
	"fmt"

	"github.com/slendycs/task-2-2/utils"
)

func main() {
	var (
		dishCount    int
		pickedDishID int
	)

	dishesRaitingList := make([]int, 0)

	// Get the number of dishes.
	dishCount, ok := utils.ReadDishesCount()
	if !ok {
		return
	}

	// Create a list of dish ratings.
	for range dishCount {
		// get the raiting of dish
		dishRaiting, ok := utils.ReadDishesRaiting()
		if !ok {
			return
		}

		dishesRaitingList = append(dishesRaitingList, dishRaiting)
	}

	// Get the selected dish.
	pickedDishID, ok = utils.ReadPickedDish(dishCount)
	if !ok {
		return
	}

	fmt.Println(utils.GetPreferredDish(pickedDishID, dishesRaitingList))
}
