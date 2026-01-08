package main

import (
"fmt"
"github.com/Xkoex/task-2-2/internal/heap"
)

func main() {
var n, pref int
fmt.Scan(&n)
if n < 1 || n > 10000 {
fmt.Println("invalid number of dishes")
return
}

dishes := make([]int, n)
for i := 0; i < n; i++ {
fmt.Scan(&dishes[i])
if dishes[i] < -10000 || dishes[i] > 10000 {
fmt.Println("invalid rating")
return
}
}

fmt.Scan(&pref)
if pref < 1 || pref > n {
fmt.Println("invalid preference number")
return
}

h := &heap.MaxHeap{}
for _, d := range dishes {
h.Push(d)
}

for i := 1; i < pref; i++ {
h.Pop()
}

fmt.Println(h.Pop())
}
