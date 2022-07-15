package route

import (
	"github.com/StellarisJAY/mnet/interface/network"
	"log"
)

type MapRouter struct {
	apis        map[byte]network.Handler
	workerQueue []chan network.HandlerContext
}

func MakeMapRouter() *MapRouter {
	return &MapRouter{apis: make(map[byte]network.Handler), workerQueue: make([]chan network.HandlerContext, 8)}
}

func (r *MapRouter) StartWorkers() {
	for i := 0; i < len(r.workerQueue); i++ {
		// create worker task queue
		r.workerQueue[i] = make(chan network.HandlerContext, 1<<16)
		go func(i int) {
			for {
				select {
				case ctx, ok := <-r.workerQueue[i]:
					// stop goroutine if work queue is closed
					if !ok {
						break
					}
					r.work(ctx)
				}
			}
		}(i)
		log.Println("Worker ", i, " started")
	}
}

func (r *MapRouter) Close() {
	// close worker queues, this will trigger worker goroutine to close
	for _, queue := range r.workerQueue {
		close(queue)
	}
}

func (r *MapRouter) Execute(ctx network.HandlerContext) {
	r.work(ctx)
}

func (r *MapRouter) Submit(ctx network.HandlerContext) {
	worker := ctx.GetPacket().ID() % uint32(len(r.workerQueue))
	r.workerQueue[worker] <- ctx
}

func (r *MapRouter) Register(typeCode byte, handler network.Handler) {
	r.apis[typeCode] = handler
}

func (r *MapRouter) work(ctx network.HandlerContext) {
	handler, ok := r.apis[ctx.GetPacket().Type()]
	if ok {
		// handle error after executing handler
		defer func() {
			if err := ctx.GetError(); err != nil {
				handler.HandleError(err, ctx)
			}
		}()
		// pre->handle->post sequence
		handler.PreHandle(ctx)
		// stop handler execution when user closes context
		if ctx.IsClosed() {
			return
		}
		handler.Handle(ctx)
		if ctx.IsClosed() {
			return
		}
		handler.PostHandle(ctx)
	}
}
