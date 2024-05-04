package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Back().Prev
		require.Equal(t, 20, middle.Value)

		middle = l.Front().Next // 20
		l.Remove(middle)        // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("elements can be of different types", func(t *testing.T) {
		l := NewList()

		l.PushFront(10)         // [10]
		l.PushBack(2.0)         // [10, 2.0]
		l.PushFront('r')        // ['r', 10, 2.0]
		l.PushBack("string")    // ['r', 10, 2.0, "string]
		l.PushFront(struct{}{}) // [struct{}{}, 'r', 10, 2.0, "string]

		elems := make([]interface{}, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value)
		}
		require.Equal(t, []interface{}{struct{}{}, int32(114), 10, 2.0, "string"}, elems)
	})

	t.Run("removing single element does not panic", func(t *testing.T) {
		l := NewList()
		l.PushFront(1)
		require.NotPanics(t, func() { l.Remove(l.Front()) })
	})
}
