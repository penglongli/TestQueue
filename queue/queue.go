package queue

type Element interface{}
type compare func(s1, s2 interface{}) int

type Queue interface {
	// Inserts the specified element into this queue
	Offer(e Element)
	// Retrieves and removes the head of this queue
	Poll() Element
	// Retrieves, but does not remove, the head of this queue
	Peek() Element
	// Get element with index
	Get(index int) Element
	// Set element with index
	Set(index int, E Element)
	// Size of this queue
	Size() int
	// Clear data of queue
	Clear()
}
