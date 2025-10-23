package utils

import (
	"fmt"

	"github.com/slendycs/task-2-2/constants"
	"github.com/slendycs/task-2-2/myerrors"
)

func ReadDishesCount() (int, error) {
	var dishCount int

	_, err := fmt.Scan(&dishCount)
	if err != nil || !(constants.MinDishCount <= dishCount && dishCount <= constants.MaxDishCount) {
		return 0, myerrors.ErrIncorrectDishesCount
	}

	return dishCount, nil
}

func ReadDishesRaiting() (int, error) {
	var dishRaiting int

	_, err := fmt.Scan(&dishRaiting)
	if err != nil || !(constants.MinRaiting <= dishRaiting && dishRaiting <= constants.MaxRaiting) {
		return 0, myerrors.ErrIncorrectRaiting
	}

	return dishRaiting, nil
}

func ReadPickedDish(dishCount int) (int, error) {
	var pickedDish int

	_, err := fmt.Scan(&pickedDish)
	if err != nil || !(1 <= pickedDish && pickedDish <= dishCount) {
		return 0, myerrors.ErrIncorrectPickedDish
	}

	return pickedDish, nil
}
