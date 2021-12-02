package s3io

import (
	"container/list"
	"time"
)

type qb struct {
	when time.Time
	s    []byte
}

// A BufferPool is a pool of byte buffers.
type BufferPool interface {
	Get() chan []byte
	Give() chan []byte
	Quit()
}

type bp struct {
	makes   int
	get     chan []byte
	give    chan []byte
	quit    chan struct{}
	timeout time.Duration
	bufsz   int64
	sizech  chan int64

	closed bool
}

func (pool *bp) Get() chan []byte {
	return pool.get
}

func (pool *bp) Give() chan []byte {
	return pool.give
}

func (pool *bp) Quit() {
	if !pool.closed {
		pool.closed = true
		close(pool.quit)
	}
}

func bufferPool(bufsz int64) BufferPool {
	sp := &bp{
		get:     make(chan []byte),
		give:    make(chan []byte),
		quit:    make(chan struct{}),
		timeout: time.Second * 30,
		sizech:  make(chan int64),
	}
	go func() {
		q := new(list.List)
		for {
			if q.Len() == 0 {
				q.PushFront(qb{when: time.Now(), s: make([]byte, bufsz)})
				sp.makes++
			}

			e := q.Front()

			timeout := time.NewTimer(sp.timeout)
			select {
			// incoming work for the pool
			case b := <-sp.give:
				timeout.Stop()
				q.PushFront(qb{when: time.Now(), s: b})
			// do work on this byte slice
			case sp.get <- e.Value.(qb).s:
				timeout.Stop()
				q.Remove(e)
			case <-timeout.C:
				// free unused slices older than timeout
				e := q.Front()
				for e != nil {
					n := e.Next()
					if time.Since(e.Value.(qb).when) > sp.timeout {
						q.Remove(e)
						e.Value = nil
					}
					e = n
				}
			case sz := <-sp.sizech: // update buffer size, free buffers
				bufsz = sz
			case <-sp.quit:
				return
			}
		}

	}()
	return sp
}
