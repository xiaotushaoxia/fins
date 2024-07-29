package fins

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

func checkIsWordMemoryArea(memoryArea byte) error {
	if memoryArea == MemoryAreaDMWord ||
		memoryArea == MemoryAreaARWord ||
		memoryArea == MemoryAreaHRWord ||
		memoryArea == MemoryAreaWRWord {
		return nil
	}
	return IncompatibleMemoryAreaError{memoryArea}
}

func checkIsBitMemoryArea(memoryArea byte) error {
	if memoryArea == MemoryAreaDMBit ||
		memoryArea == MemoryAreaARBit ||
		memoryArea == MemoryAreaHRBit ||
		memoryArea == MemoryAreaWRBit {
		return nil
	}
	return IncompatibleMemoryAreaError{memoryArea}
}

type atomicByte struct {
	m sync.Mutex
	v byte
}

func (ab *atomicByte) load() byte {
	ab.m.Lock()
	defer ab.m.Unlock()
	return ab.v
}

func (ab *atomicByte) increment() byte {
	ab.m.Lock()
	defer ab.m.Unlock()
	ab.v++
	return ab.v
}

func waitMoment(ctx context.Context, d time.Duration) {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return
	case <-timer.C:
		return
	}
}

type syncRespSlice struct {
	m  sync.Mutex
	rs [256]chan *response
}

func (s *syncRespSlice) getW(i byte) chan<- *response {
	s.m.Lock()
	defer s.m.Unlock()
	return s.rs[i]
}

func (s *syncRespSlice) getR(i byte) <-chan *response {
	s.m.Lock()
	defer s.m.Unlock()
	return s.rs[i]
}

func (s *syncRespSlice) set(i byte, c chan *response) {
	s.m.Lock()
	defer s.m.Unlock()
	s.rs[i] = c
}

// singleflightOne
// idea from github.com/golang/x/sync/singleflight and remove return value and Key/Group
type singleflightOne struct {
	mu    sync.Mutex
	wg    *sync.WaitGroup
	count atomic.Int32
}

func (sr *singleflightOne) do(f func()) {
	wg, first := sr.prepare()
	if !first {
		wg.Wait()
		return
	}
	defer func() { // cleanup even if panic
		sr.mu.Lock()
		defer sr.mu.Unlock()
		sr.wg = nil
		wg.Done()
	}()
	f()
}

func (sr *singleflightOne) prepare() (wg0 *sync.WaitGroup, first bool) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	if wg := sr.wg; wg != nil {
		return wg, false
	}
	wg := new(sync.WaitGroup)
	wg.Add(1)
	sr.wg = wg
	return wg, true
}
