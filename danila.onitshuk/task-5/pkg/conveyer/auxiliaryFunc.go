package conveyer

func (c *Conveyer) makeChannel(name string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.channels[name]; !ok {
		c.channels[name] = make(chan string, c.channelSize)
	}
}

func (c *Conveyer) makeChannels(names ...string) {
	for _, name := range names {
		c.makeChannel(name)
	}
}
