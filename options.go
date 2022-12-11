package hel

type Option func(*Configuration)

// WithBlockingLaunch function sets the BlockUntilLaunched option.
func WithBlockingLaunch(block bool) Option {
	return func(c *Configuration) {
		c.BlockUntilLaunched = block
	}
}

// WithBindAddress function sets the BindAddress option.
func WithBindAddress(address string) Option {
	return func(c *Configuration) {
		c.BindAddress = address
	}
}
