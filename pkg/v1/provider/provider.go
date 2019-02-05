package provider

// Provider.
// Enables an application to add a piece of functionality very quickly.
// This normally results a connection to an external service being setup.
// If you want to use a Provider to start an HTTP service or any other blocking functionality, use the RunProvider.
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

// Override if the Provider needs to be initialized.
func (p *AbstractProvider) Init() error {
	return nil
}

// Override if the Provider needs to be closed.
func (p *AbstractProvider) Close() error {
	return nil
}

// Abstract RunProvider.
// Does not extend the Run() method, since Providers that don't actually run shouldn't be a RunProvider.
type AbstractRunProvider struct {
	RunProvider

	running bool
}

// Override if the RunProvider needs to be initialized.
func (p *AbstractRunProvider) Init() error {
	return nil
}

// Override if the RunProvider needs to be closed. Make sure to update call SetRunning(false).
func (p *AbstractRunProvider) Close() error {
	p.SetRunning(false)
	return nil
}

// Returns true after the RunProvider has started.
func (p *AbstractRunProvider) IsRunning() bool {
	return p.running
}

// Used by extending providers to update their running status. Should be called with true once the Run() method has almost finished (just before the blocking part).
func (p *AbstractRunProvider) SetRunning(running bool) {
	p.running = running
}
