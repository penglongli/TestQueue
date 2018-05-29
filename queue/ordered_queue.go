package queue

import (
	"sync"
)

var (
	head     = &Node{}
	tail     = &Node{}
)

const DEPTH = 5

type Node struct {
	value Element
	prev  [DEPTH]*Node
	next  [DEPTH]*Node
}

type linkedList struct {
	sync.RWMutex
	head        *Node
	tail        *Node
	size        int
	compareWith compare
}

func NewLinkedList(f compare) *linkedList {
	list := &linkedList{
		head:        head,
		tail:        tail,
		compareWith: f,
	}
	list.head.next[0] = tail
	list.tail.prev[0] = head
	return list
}

func (list *linkedList) Offer(E interface {
	Element
}) {
	list.Lock()
	defer func() {
		list.size++
		list.Unlock()
	}()

	newNode := &Node{value: E}
	if list.size == 0 {
		list.insertAfter(newNode, head)
		return
	}

	currNode := head.next[0]
	for {
		// 当前节点 < 插入节点
		if list.compareWith(currNode.value, E) == 0 {
			for i := 0; i < DEPTH; i++ {
				t := currNode.prev[i]
				if t == head {
					list.insertAfter(newNode, head)
					return
				}

				if list.compareWith(t.value, E) != 0 {
					list.insertAfter(newNode, t)
					return
				}
			}
		}

		// 当前节点 >= 插入节点
		nextNode := currNode.next[DEPTH-1]
		if nextNode == nil || nextNode == tail {
			for i := 0; i < DEPTH; i++ {
				t := currNode.next[i]
				if t == tail {
					list.insertBefore(newNode, tail)
					return
				}

				if list.compareWith(t.value, E) == 0 {
					list.insertBefore(newNode, t)
					return
				}
			}
		}
		currNode = nextNode
	}
}

func (list *linkedList) Poll() Element {
	list.Lock()
	defer func() {
		list.size--
		list.Unlock()
	}()

	result := head.next[0]
	if result == nil || result == tail {
		return nil
	}

	for i := 0; i < DEPTH; i++ {
		currNode := result.next[i]
		if currNode == nil {
			break
		}

		currNode.prev[i] = head
		head.next[i] = currNode
	}
	return result.value
}

func (list *linkedList) Peek() Element {
	list.RLock()
	defer list.RUnlock()

	result := head.next[0]
	if result == nil || result == tail {
		return nil
	}
	return result.value
}

func (list *linkedList) Get(index int) Element {
	list.RLock()
	defer list.RUnlock()

	n := list.indexOf(index).value
	return n
}

func (list *linkedList) Set(index int, E Element) {
	list.Lock()
	defer list.Unlock()

	n := list.indexOf(index)
	if n == nil {
		return
	}
	n.value = E
}

func (list *linkedList) Size() int {
	return list.size
}

func (list *linkedList) Clear() {
	list.Lock()
	defer list.Unlock()

	head     = &Node{}
	tail     = &Node{}
	list.head.next[0] = tail
	list.tail.prev[0] = head
	list.size = 0
}

// Attention: 调用前需要加锁
func (list *linkedList) indexOf(index int) *Node {
	size := list.size
	if index < 0 || index > size {
		return nil
	}
	if index < DEPTH {
		return head.next[index]
	}
	if size-DEPTH <= index {
		return tail.prev[size-index-1]
	}

	if size/2 >= index {
		currNode, currIndex := head, 0

		for {
			if currIndex+DEPTH > index {
				break
			}
			currNode = currNode.next[DEPTH-1]
			currIndex = currIndex + DEPTH - 1
		}
		return currNode.next[index-currIndex]
	} else {
		currNode, currIndex := tail, size-1
		for {
			if currIndex-DEPTH < index {
				break
			}
			currNode = currNode.prev[DEPTH-1]
			currIndex = currIndex - DEPTH + 1
		}
		return currNode.prev[currIndex-index-1]
	}
}

func (list *linkedList) insertBefore(newNode *Node, node *Node) {
	prevNode := node.prev[0]
	list.insertAfter(newNode, prevNode)
}

func (list *linkedList) insertAfter(newNode *Node, node *Node) {
	leftNode, rightNode := node, node.next[0]

	// 左侧链表与 newNode 建立关联
	for i := 0; i < DEPTH-1; i++ {
		currNode := leftNode.prev[i]
		if currNode == nil {
			break
		}

		for j := DEPTH - 1; j > i+1; j-- {
			currNode.next[j] = currNode.next[j-1]
		}
		currNode.next[i+1] = newNode
		newNode.prev[i+1] = currNode
	}
	for j := DEPTH - 1; j > 0; j-- {
		leftNode.next[j] = leftNode.next[j-1]
	}
	leftNode.next[0] = newNode
	newNode.prev[0] = leftNode

	// 右侧链表与 newNode 建立关联
	for i := 0; i < DEPTH-1; i++ {
		currNode := rightNode.next[i]
		if currNode == nil {
			break
		}

		for j := DEPTH - 1; j > i+1; j-- {
			currNode.prev[j] = currNode.prev[j-1]
		}
		currNode.prev[i+1] = newNode
		newNode.next[i+1] = currNode
	}
	for j := DEPTH - 1; j > 0; j-- {
		rightNode.prev[j] = rightNode.prev[j-1]
	}
	rightNode.prev[0] = newNode
	newNode.next[0] = rightNode
}
