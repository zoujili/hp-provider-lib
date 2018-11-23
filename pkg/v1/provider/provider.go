package provider

// Provider ...
type Provider interface {
	Init() error
	Close() error
}

// RunProvider ...
type RunProvider interface {
	Provider

	Run() error
	IsRunning() bool
}
