package conveyer

func (c *Conveyer) makeChannel(name string) {
	_, ok := c.channels[name]
	if !ok {
		c.channels[name] = make(chan string, c.size)
	}
}

func (c *Conveyer) makeChannels(names ...string) {
	for _, name := range names {
		c.makeChannel(name)
	}
}
