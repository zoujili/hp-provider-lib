package provider

// Provider.
// TODO: Explain what a provider is.
type Provider interface {
	Init() error  // Initializes the Provider as it is added to the Stack.
	Close() error // Stops the Provider (if runnable) and closes any open routines.
}

// RunProvider.
// A RunProvider differs from a normal Provider in that it has a blocking Run() method.
// This is mostly used for creating HTTP services.
type RunProvider interface {
	Provider

	Run() error      // Blocking function, starting up any services provided by the Provider.
	IsRunning() bool // Returns true only after the Provider has fully started up (making it usable by other functions).
}

// Abstract Provider.
type AbstractProvider struct {
	Provider
}

// Extend if the Provider needs to be initialized.
func (p *AbstractProvider) Init() error {
	return nil
}

// Extend if the Provider needs to be closed.
func (p *AbstractProvider) Close() error {
	return nil
}

// Abstract RunProvider.
// Does not extend the Run() method, since Providers that don't actually run shouldn't be a RunProvider.
type AbstractRunProvider struct {
	RunProvider

	running bool
}

// Extend if the Provider needs to be initialized.
func (p *AbstractRunProvider) Init() error {
	return nil
}

// Extend if the Provider needs to be closed.
func (p *AbstractRunProvider) Close() error {
	return nil
}

// Returns true after the Provider has started.
func (p *AbstractRunProvider) IsRunning() bool {
	return p.running
}

// Allows providers to set themselves as running.
func (p *AbstractRunProvider) SetRunning(running bool) {
	p.running = running
}
