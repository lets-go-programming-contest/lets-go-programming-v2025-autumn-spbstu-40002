package main

import (
	"fmt"

	"github.com/slendycs/task-2-2/utils"
)

func main() {
	var (
		dishCount    int
		pickedDishId int
	)

	dishesRaitingList := make([]int, 0)

	// get the number of dishes
	dishCount, ok := utils.ReadDishesCount()
	if !ok {
		return
	}

	// сreate a list of dish ratings
	for range dishCount {
		// get the raiting of dish
		dishRaiting, ok := utils.ReadDishesRaiting()
		if !ok {
			return
		}

		dishesRaitingList = append(dishesRaitingList, dishRaiting)
	}

	// get the selected dish
	pickedDishId, ok = utils.ReadPickedDish(dishCount)
	if !ok {
		return
	}

	fmt.Println(utils.GetPreferredDish(pickedDishId, dishesRaitingList))
}
