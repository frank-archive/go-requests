package client

// Use appends the passed middleware to the end of the list
func (c *Client) Use(mw Middleware) {
	c.mwLock.Lock()
	defer c.mwLock.Unlock()
	c.middlewares = append(c.middlewares, mw)
	c.buildHandler()
}

func (c *Client) UseIndex(mw Middleware, idx int) {
	c.mwLock.Lock()
	defer c.mwLock.Unlock()
	if len(c.middlewares) <= idx {
		c.middlewares = append(c.middlewares, mw)
	} else {
		c.middlewares = append(c.middlewares[:idx+1], c.middlewares[idx:]...)
		c.middlewares[idx] = mw
	}
	c.buildHandler()
}

func (c *Client) buildHandler() {
	// rebuilds the entire middleware chain
	chain := c.request
	for _, mw := range c.middlewares {
		chain = mw(chain)
	}
	c.chainedHandler = chain
}
