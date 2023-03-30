package housekeeper

import (
	"sync"
	"time"
)

type Garbage interface {
	ID() int
	TrashOut(i time.Time)
}

type Housekeeper struct {
	lock 		*sync.Mutex
	stock		[]Garbage
	Interval	time.Duration
}

func New(i time.Duration) *Housekeeper {
	return &Housekeeper{
		lock: &sync.Mutex{},
		stock: make([]Garbage, 0),
		Interval: i,
	}
}

func (h *Housekeeper) Place(n Garbage) {
	h.lock.Lock()
	h.stock = append(h.stock, n)
	h.lock.Unlock()
}

func (h *Housekeeper) Remove(id int) {
	h.lock.Lock()
	el := -1
    for i := range h.stock{
		if h.stock[i].ID() == id {
			el = i
			break
		}
	}

	if el != -1 {
		h.stock[el] = h.stock[len(h.stock) - 1]
		h.stock = h.stock[:len(h.stock) - 1]
	}

	h.lock.Unlock()
}

func (h *Housekeeper) Clear() {
	h.lock.Lock()
	h.stock = nil
	h.lock.Unlock()
}

func (h *Housekeeper) Start() {
	for {
		<-time.After(h.Interval)
		limit := time.Now()
		h.lock.Lock()
		for i := range h.stock {
			h.stock[i].TrashOut(limit)
		}
		h.lock.Unlock()
	}
}