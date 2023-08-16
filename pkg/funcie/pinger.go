package funcie

// Pinger is an interface for pinging applications.
type Pinger interface {
	Ping(app Application) error
}
