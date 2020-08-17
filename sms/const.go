package sms

// Driver identifies source code management driver.
type Driver int

// Drivers we support
const (
	DriverUnknown Driver = iota
	DriverFdi
)

func (d Driver) String() (s string) {
	switch d {
	case DriverFdi:
		return "fdi"
	default:
		return "unknown"
	}
}
