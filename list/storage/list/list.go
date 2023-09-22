package list

import "fmt"

type List struct {
	length    int64
	firstNode *node
}

func NewList() *List {
	return &List{length: 0, firstNode: nil}
}

func (l *List) Len() int64 {
	return l.length
}

func (l *List) Add(value int64) int64 {
	if l.firstNode == nil {
		l.firstNode = &node{value: value, next: nil}
	} else {
		current := l.firstNode
		for current.next != nil {
			current = current.next
		}
		current.next = &node{value: value, next: nil}
	}
	l.length++
	return l.length //
}

func (l *List) RemoveByIndex(index int64) {
	if l.firstNode == nil || index < 0 || index >= l.length {
		return
	}

	if index == 0 {
		l.firstNode = l.firstNode.next
	} else {
		current := l.firstNode
		for i := int64(0); i < index-1; i++ {
			current = current.next
		}
		current.next = current.next.next
	}
	l.length--
}

func (l *List) RemoveByValue(value int64) {
	if l.firstNode == nil {
		return
	}

	if l.firstNode.value == value {
		l.firstNode = l.firstNode.next
		l.length--
		return
	}

	current := l.firstNode
	for current.next != nil {
		if current.next.value == value {
			current.next = current.next.next
			l.length--
			return
		}
		current = current.next
	}
}

func (l *List) RemoveAllByValue(value int64) {
	if l.firstNode == nil {
		return
	}

	for l.firstNode != nil && l.firstNode.value == value {
		l.firstNode = l.firstNode.next
		l.length--
	}

	current := l.firstNode
	for current != nil && current.next != nil {
		if current.next.value == value {
			current.next = current.next.next
			l.length--
		} else {
			current = current.next
		}
	}
}

func (l *List) GetByIndex(index int64) (int64, bool) {
	if l.firstNode == nil || index < 0 || index >= l.length {
		return 0, false
	}

	current := l.firstNode
	for i := int64(0); i < index; i++ {
		current = current.next
	}
	return current.value, true
}

func (l *List) GetByValue(value int64) (int64, bool) {
	if l.firstNode == nil {
		return 0, false
	}

	current := l.firstNode
	index := int64(0)
	for current != nil {
		if current.value == value {
			return index, true
		}
		current = current.next
		index++
	}

	return 0, false
}

func (l *List) GetAllByValue(value int64) ([]int64, bool) {
	if l.firstNode == nil {
		return nil, false
	}

	ids := make([]int64, 0)
	current := l.firstNode
	index := int64(0)
	for current != nil {
		if current.value == value {
			ids = append(ids, index)
		}
		current = current.next
		index++
	}

	if len(ids) == 0 {
		return nil, false
	}

	return ids, true
}

func (l *List) GetAll() ([]int64, bool) {
	if l.firstNode == nil {
		return nil, false
	}

	values := make([]int64, 0)
	current := l.firstNode
	for current != nil {
		values = append(values, current.value)
		current = current.next
	}

	return values, true
}

func (l *List) Clear() {
	l.length = 0
	l.firstNode = nil
}

func (l *List) Print() {
	if l.firstNode == nil {
		fmt.Println("no data")
	}
	current := l.firstNode
	for current != nil {
		fmt.Printf("%d ", current.value)
		current = current.next
	}
	fmt.Println()
}
