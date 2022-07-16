package pool

import (
	"fmt"
	"math/rand"
	"sync"
)

// SharablePool a resource pool with limited capacity and resources can be accessed by multiple threads once
type SharablePool struct {
	sync.Mutex
	capacity int
	size     int
	cache    []interface{}
	newFunc  func() interface{}
}

func NewSharablePool(capacity, initSize int, newFunc func() interface{}) *SharablePool {
	if capacity <= 0 || initSize < 0 {
		panic(fmt.Errorf("invalid argument for New Pool"))
	}
	if initSize > capacity {
		initSize = capacity
	}
	p := new(SharablePool)
	p.capacity = capacity
	p.size = 0
	p.cache = make([]interface{}, capacity)
	p.newFunc = newFunc
	p.warmUp(initSize)
	return p
}

func (p *SharablePool) Get() interface{} {
	if p.size < p.capacity {
		p.warmUp(p.capacity)
	}
	return p.cache[p.randSlot()]
}

func (p *SharablePool) GetHashed(i int) interface{} {
	if p.size < p.capacity {
		p.warmUp(p.capacity)
	}
	return p.cache[i%p.capacity]
}

// this function warms up the resource pool to the target size
func (p *SharablePool) warmUp(count int) {
	p.Lock()
	defer p.Unlock()
	for i := p.size; i < count; i++ {
		p.cache[i] = p.newFunc()
		p.size++
	}
}

func (p *SharablePool) randSlot() int {
	i := rand.Int()
	return i % p.capacity
}

func (p *SharablePool) Put(e interface{}) {

}
