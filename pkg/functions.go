package pkg

import "sync"

type Pubsub struct {
	channels []chan struct{}
	lock     *sync.RWMutex
}

func Filter[T any](Slices []T, fn func(T) bool) []T {

	filtered := make([]T, 0, len(Slices))

	for _, item := range Slices {
		if fn(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func Newpubsub() *Pubsub {
	return &Pubsub{
		channels: make([]chan struct{}, 0),
		lock:     new(sync.RWMutex),
	}
}

func (p *Pubsub) Subscribe() (<-chan struct{}, func()) {
	p.lock.Lock()
	defer p.lock.Unlock()
	c := make(chan struct{}, 1)
	p.channels = append(p.channels, c)
	return c, func() {
		p.lock.Lock()
		defer p.lock.Unlock()

		for i, channel := range p.channels {
			if channel == c {
				p.channels = append(p.channels[:i], p.channels[i+1:]...)
				close(c)
				return
			}
		}

	}
}

func (p *Pubsub) Publish() {
	p.lock.RLock()
	defer p.lock.RUnlock()
	for _, channel := range p.channels {
		channel <- struct{}{}
	}
}
