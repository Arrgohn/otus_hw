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
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	back  *ListItem
	front  *ListItem
	len  int
}

func (l *list) Len() int {
	{ return l.len }
}

func (l *list) Front() *ListItem {
	if l.len == 0 {
		return nil
	}
	return l.front
}

func (l *list) Back() *ListItem {
	if l.len == 0 {
		return nil
	}
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	e := &ListItem{Value: v}

	if l.front == nil {
		l.front = e
		l.back = e

		l.len++
		return e
	}

	l.front.Prev = e
	e.Next = l.front
	l.front = e

	l.len++
	return e
}

func (l *list) PushBack(v interface{}) *ListItem {
	e := &ListItem{Value: v}

	if l.front == nil {
		l.front = e
		l.back = e

		l.len++
		return e
	}

	l.back.Next = e
	e.Prev = l.back
	l.back = e

	l.len++
	return e
}

func (l *list) Remove(i *ListItem) {
	if i == l.front {
		i.Next.Prev = i.Prev
		l.front = i.Next
	}

	if i == l.back {
		i.Prev.Next = i.Next
		l.back = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
