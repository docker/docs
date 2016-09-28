package worker

import (
	"container/heap"
	"sync"
	"time"

	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
)

// This heap returns the oldest time first when popping or peeking
// This fancy heap also has a map so we can look up the location of an element in the heap
// by its job id
type TimeHeap struct {
	//lock   sync.Mutex
	lookup map[string]int
	list   []*schema.Job
}

// This is a wrapper around a timeheap to make it thread safe
type SafeTimeHeap struct {
	lock sync.Mutex
	heap *TimeHeap
}

func NewSafeTimeHeap() *SafeTimeHeap {
	h := &SafeTimeHeap{
		heap: &TimeHeap{
			lookup: map[string]int{},
			list:   []*schema.Job{},
		},
	}
	heap.Init(h.heap)
	return h
}

// the following 5 operations are not locked by a mutex and therefore should not be called directly
// even though they have to be public to support the heap interface
func (h TimeHeap) Len() int {
	return len(h.list)
}
func (h TimeHeap) Less(i, j int) bool {
	return h.list[i].HeartbeatExpiration.After(h.list[j].HeartbeatExpiration)
}
func (h TimeHeap) Swap(i, j int) {
	h.list[i], h.list[j] = h.list[j], h.list[i]
	h.lookup[h.list[i].ID] = i
	h.lookup[h.list[j].ID] = j
}

func (h *TimeHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	entry := x.(*schema.Job)
	h.list = append(h.list, entry)
	h.lookup[entry.ID] = len(h.list) - 1
}

func (h *TimeHeap) Pop() interface{} {
	old := h.list
	n := len(old)
	if n == 0 {
		return nil
	}
	x := old[n-1]
	h.list = old[0 : n-1]
	delete(h.lookup, x.ID)
	return x
}

func (h *SafeTimeHeap) Len() int {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.heap.Len()
}

func (h *SafeTimeHeap) Peek() *schema.Job {
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.peek()
}

func (h *SafeTimeHeap) peek() *schema.Job {
	if len(h.heap.list) > 0 {
		return h.heap.list[len(h.heap.list)-1]
	} else {
		return nil
	}
}

func (h *SafeTimeHeap) PopIfExpired(now time.Time) *schema.Job {
	h.lock.Lock()
	defer h.lock.Unlock()
	peek := h.peek()
	if peek != nil && !peek.HeartbeatExpiration.After(now) {
		tmp := h.heap
		tmp2 := tmp.Pop()
		tmp3 := tmp2.(*schema.Job)
		return tmp3
	}
	return nil
}

func (h *SafeTimeHeap) Update(change *schema.JobChange) {
	h.lock.Lock()
	defer h.lock.Unlock()
	// a job was removed
	if change.NewValue == nil {
		pos, ok := h.heap.lookup[change.OldValue.ID]
		if ok {
			heap.Remove(h.heap, pos)
		}
		delete(h.heap.lookup, change.OldValue.ID)
		return
	}
	// a job was marked as done or at least we don't care any more
	if change.NewValue != nil && schema.StatusIsFinished(change.NewValue.Status) {
		pos, ok := h.heap.lookup[change.NewValue.ID]
		if ok {
			heap.Remove(h.heap, pos)
		}
		delete(h.heap.lookup, change.NewValue.ID)
		return
	}
	// a job was created or updated in every other case
	pos, ok := h.heap.lookup[change.NewValue.ID]
	if ok {
		// update the time if it changed
		h.heap.list[pos] = change.NewValue
		heap.Fix(h.heap, pos)
	} else {
		heap.Push(h.heap, change.NewValue)
	}
}
