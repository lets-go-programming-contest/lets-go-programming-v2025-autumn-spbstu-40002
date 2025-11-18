package conveyer

func (c *Conveyer) makeChannel(name string) {
	if _, ok := c.channels[name]; !ok {
		c.channels[name] = make(chan string, c.channelSize)
	}
}

func (c *Conveyer) makeChannels(names ...string) {
	for _, name := range names {
		c.makeChannel(name)
	}
}

func (c *Conveyer) closeChannels() {
	for _, ch := range c.channels {
		close(ch)
	}
}
