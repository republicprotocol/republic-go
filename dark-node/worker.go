package node

import (
	"github.com/republicprotocol/republic-go/order"
)

type OrderFragmentWorkerPool struct {
	work    chan *order.Fragment
	workers []OrderFragmentWorker
}

func NewOrderFragmentWorkerPool(buffer, workers int) *OrderFragmentWorkerPool {
	w := make(chan *order.Fragment, buffer)
	ws := make([]OrderFragmentWorker, workers)
	for i := range ws {
		ws[i] = NewOrderFragmentWorker(w)
		go ws[i].Run()
	}
	return &OrderFragmentWorkerPool{
		work:    w,
		workers: ws,
	}
}

type OrderFragmentWorker struct {
	work chan *order.Fragment
}

func NewOrderFragmentWorker(work chan *order.Fragment) OrderFragmentWorker {
	return OrderFragmentWorker{
		work: work,
	}
}

func (worker *OrderFragmentWorker) Run() {
	for work := range worker.work {
		worker.Process(work)
	}
}

func (worker *OrderFragmentWorker) Process(orderFragment *order.Fragment) {
	worker.
}
