package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Prev  *ListItem
	Next  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newListItem := ListItem{
		Value: v,
		Prev:  nil,
		Next:  nil,
	}

	if l.len == 0 {
		l.front = &newListItem
		l.back = &newListItem
	} else {
		newListItem.Next = l.front
		l.front.Prev = &newListItem
		l.front = &newListItem
	}

	l.len++

	return &newListItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newListItem := ListItem{
		Value: v,
		Prev:  nil,
		Next:  nil,
	}

	if l.len == 0 {
		l.front = &newListItem
		l.back = &newListItem
	} else {
		newListItem.Prev = l.back
		l.back.Next = &newListItem
		l.back = &newListItem
	}

	l.len++

	return &newListItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil && i.Next == nil {
		return
	}
	if i.Prev == nil {
		i.Next.Prev = nil
		l.front = i.Next
	}
	if i.Next == nil {
		i.Prev.Next = nil
		l.back = i.Prev
	}
	if i.Prev != nil && i.Next != nil {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil { // значит i уже является первым
		return
	}
	l.Remove(i)
	i.Prev = nil
	i.Next = l.front
	l.front.Prev = i
	l.front = i
	l.len++ // потому что после Remove() длина уменьшилась
}
