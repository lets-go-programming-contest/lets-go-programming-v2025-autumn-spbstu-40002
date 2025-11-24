package conveyer

import "sync"

type conveyor struct {
	chSize int
	mu     sync.RWMutex
	chans  map[string]chan string
}

func newConveyor(size int) *conveyor {
	return &conveyor{
		chSize: size,
		chans:  make(map[string]chan string),
	}
}

func (c *conveyor) ensureChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ch, ok := c.chans[name]; ok {
		return ch
	}
	ch := make(chan string, c.chSize)
	c.chans[name] = ch
	return ch
}

func (c *conveyor) getChan(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ch, ok := c.chans[name]
	return ch, ok
}

func (c *conveyor) closeAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, ch := range c.chans {
		close(ch)
		delete(c.chans, k)
	}
}
