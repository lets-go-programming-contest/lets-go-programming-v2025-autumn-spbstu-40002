package main

import "fmt"

func main() {
    var departmentCount int
    fmt.Scan(&departmentCount)

    for i := 0; i < departmentCount; i++ {
        var employeesCount int
        fmt.Scan(&employeesCount)

        minTemp := 15
        maxTemp := 30

        for j := 0; j < employeesCount; j++ {
            var operator string
            var temp int
            fmt.Scan(&operator, &temp)

            if operator == ">=" {
                if temp > minTemp {
                    minTemp = temp
                }
            } else if operator == "<=" {
                if temp < maxTemp {
                    maxTemp = temp
                }
            }

            if minTemp <= maxTemp {
                fmt.Println(minTemp)
            } else {
                fmt.Println(-1)
            }
        }
        if i < departmentCount-1 {
            fmt.Println()
        }
    }
}
