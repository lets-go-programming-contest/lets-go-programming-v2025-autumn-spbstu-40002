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
	dishCount, success := utils.ReadDishesCount()
	if !success {
		return
	}

	// Create a list of dish ratings.
	for range dishCount {
		// get the raiting of dish
		dishRaiting, success := utils.ReadDishesRaiting()
		if !success {
			return
		}

		dishesRaitingList = append(dishesRaitingList, dishRaiting)
	}

	// Get the selected dish.
	pickedDishID, success = utils.ReadPickedDish(dishCount)
	if !success {
		return
	}

	fmt.Println(utils.GetPreferredDish(pickedDishID, dishesRaitingList))
}
