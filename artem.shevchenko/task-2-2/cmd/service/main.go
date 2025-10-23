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
	dishCount, err := utils.ReadDishesCount()
	if err != nil {
		fmt.Println(err)
	}

	// Create a list of dish ratings.
	for range dishCount {
		// get the raiting of dish
		dishRaiting, err := utils.ReadDishesRaiting()
		if err != nil {
			fmt.Println(err)
		}

		dishesRaitingList = append(dishesRaitingList, dishRaiting)
	}

	// Get the selected dish.
	pickedDishID, err = utils.ReadPickedDish(dishCount)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(utils.GetPreferredDish(pickedDishID, dishesRaitingList))
}
