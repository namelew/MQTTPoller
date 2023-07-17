package waitgroup

import "sync"

type WaitGroup struct {
	members int
	done    chan interface{}
	m       *sync.Mutex
}

func New() *WaitGroup {
	return &WaitGroup{
		members: 0,
		done:    make(chan interface{}, 1),
		m:       &sync.Mutex{},
	}
}

func (w *WaitGroup) Add(n int) {
	w.m.Lock()
	w.members = n
	w.m.Unlock()
}

func (w *WaitGroup) Done() {
	w.m.Lock()

	if w.members > 0 {
		w.members--
	}

	if w.members <= 0 {
		w.done <- true
	}

	w.m.Unlock()
}

func (w *WaitGroup) Wait() {
	<-w.done
}
