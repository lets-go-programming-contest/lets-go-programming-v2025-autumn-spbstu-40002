package main

import "fmt"

func main() {
var N int
fmt.Scan(&N)

for i := 0; i < N; i++ {
var K int
fmt.Scan(&K)

minTemp := 15
maxTemp := 30

for j := 0; j < K; j++ {
var operator string
var temp int
fmt.Scan(&operator, &temp)

switch operator {
case ">=":
if temp > minTemp {
minTemp = temp
}
case "<=":
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
}
}
