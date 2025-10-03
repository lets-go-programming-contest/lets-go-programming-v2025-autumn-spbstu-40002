package utils

import (
	"fmt"

	"github.com/slendycs/task-2-2/constants"
	"github.com/slendycs/task-2-2/myErrors"
)

func ReadDishesCount() (int, bool) {
	var dishCount int

	_, err := fmt.Scan(&dishCount)
	if err != nil || !(constants.MinDishCount <= dishCount && dishCount <= constants.MaxDishCount) {
		fmt.Println(myErrors.ErrIncorrectDishesCount)

		return 0, false
	}

	return dishCount, true
}

func ReadDishesRaiting() (int, bool) {
	var dishRaiting int

	_, err := fmt.Scan(&dishRaiting)
	if err != nil || !(constants.MinRaiting <= dishRaiting && dishRaiting <= constants.MaxRaiting) {
		fmt.Println(myErrors.ErrIncorrectRaiting)

		return 0, false
	}

	return dishRaiting, true
}

func ReadPickedDish(dishCount int) (int, bool) {
	var pickedDish int

	_, err := fmt.Scan(&pickedDish)
	if err != nil || !(1 <= pickedDish && pickedDish <= dishCount) {
		fmt.Println(myErrors.ErrIncorrectPickedDish)

		return 0, false
	}

	return pickedDish, true
}
