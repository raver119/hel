package hel

// Configuration is a struct that contains all the configuration options for the server.
type Configuration struct {
	BlockUntilLaunched bool
	BindAddress        string
}

var DefaultConfiguration = Configuration{
	BlockUntilLaunched: false,
	BindAddress:        "0.0.0.0",
}
