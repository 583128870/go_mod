package cache

import (
	"container/heap"
	"container/list"
	"container/ring"
	"fmt"
	"log"
	"testing"
	"time"
)

func Test(t *testing.T) {
	statistics := NewMemCache(&testStringModel{})
	//value, rel, err := statistics.GetCacheInfo("5")
	//if err != nil {
	//	fmt.Println(err)
	//}
	for {
		time.Sleep(1 * time.Second)
		allData, err := statistics.GetCacheData(false)
		if err != nil {
			log.Printf("获取缓存数据错误：%s", err.Error())
			continue
		}
		fmt.Println(allData)
	}

}

// An IntHeap is a min-heap of ints.
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// This example inserts several ints into an IntHeap, checks the minimum,
// and removes them in order of priority.
func Test_intHeap(t *testing.T) {
	h := &IntHeap{2, 1, 5}
	heap.Init(h)
	heap.Push(h, 3)
	fmt.Printf("minimum: %d\n", (*h)[0])
	fmt.Println(*h)
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}
	for len(*h) > 2 {
		fmt.Println(heap.Pop(h))
	}
	for h.Len() > 0 {
		fmt.Printf("%d ", heap.Pop(h))
	}
	fmt.Println(*h)
	// Output:
	// minimum: 1
	// 1 2 3 5
}

func Test_List(t *testing.T) {
	// Create a new list and put some numbers in it.
	l := list.New()
	e4 := l.PushBack(4)
	e1 := l.PushFront(1)
	l.InsertBefore(3, e4)
	l.InsertAfter(2, e1)

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
	// Output:
	// 1
	// 2
	// 3
	// 4
}
func Test_ExampleRing_Do(t *testing.T) {
	// Create a new ring of size 5
	r := ring.New(5)

	// Get the length of the ring
	n := r.Len()

	// Initialize the ring with some integer values
	for i := 0; i < n; i++ {
		r.Value = i
		r = r.Next()
	}

	// Iterate through the ring and print its contents
	r.Do(func(p any) {
		fmt.Println(p.(int))
	})

	for {
		r = r.Next()
		fmt.Println(r.Value)
		time.Sleep(1 * time.Second)
	}
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}
