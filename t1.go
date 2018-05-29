package main

import (
	"fmt"
	"TestQueue/queue"
)

type Orange struct {
	weight float32
}

func compareOrange(s1, s2 interface{}) int {
	if s1.(Orange).weight > s2.(Orange).weight {
		return 1
	} else {
		return 0
	}
}

func main() {
	a1 := Orange{weight: 6.32}
	a2 := Orange{weight: 2.32}
	a3 := Orange{weight: 1.11}
	a4 := Orange{weight: 1.01}
	a5 := Orange{weight: 1.22}
	a6 := Orange{weight: 3.45}
	a7 := Orange{weight: 8.69}
	a8 := Orange{weight: 3.4}

	list := queue.NewLinkedList(compareOrange)
	list.Offer(a1)
	list.Offer(a2)
	list.Offer(a3)
	list.Offer(a4)
	list.Offer(a5)
	list.Offer(a6)
	list.Offer(a7)
	list.Offer(a8)

	for {
		item := list.Poll()
		if item == nil {
			break
		}
		fmt.Println(item.(Orange).weight)
	}
}
