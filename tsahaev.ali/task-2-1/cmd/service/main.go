package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Читаем N и K
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	scanner.Scan()
	k, _ := strconv.Atoi(scanner.Text())

	// Обрабатываем каждый отдел
	for i := 0; i < n; i++ {
		// Диапазон температур для текущего отдела
		minTemp := 15
		maxTemp := 30
		possible := true

		// Обрабатываем каждого сотрудника в отделе
		for j := 0; j < k; j++ {
			scanner.Scan()
			line := scanner.Text()
			parts := strings.Split(line, " ")

			operator := parts[0]
			temp, _ := strconv.Atoi(parts[1])

			// Обновляем диапазон в зависимости от требования
			if operator == ">=" {
				if temp > minTemp {
					minTemp = temp
				}
			} else if operator == "<=" {
				if temp < maxTemp {
					maxTemp = temp
				}
			}

			// Проверяем, возможен ли диапазон
			if minTemp <= maxTemp && possible {
				// Выбираем минимально возможную температуру (как в примере)
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
				possible = false
			}
		}
	}
}
