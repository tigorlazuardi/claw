package types

// NSFWMode represents the NSFW handling mode for devices
type NSFWMode int

const (
	// NSFWModeUnspecified indicates no specific NSFW mode is set
	NSFWModeUnspecified NSFWMode = 0
	// NSFWModeBlock blocks NSFW content
	NSFWModeBlock NSFWMode = 1
	// NSFWModeAccept accepts NSFW content
	NSFWModeAccept NSFWMode = 2
	// NSFWModeOnly only accepts NSFW content
	NSFWModeOnly NSFWMode = 3
)

// String returns a human-readable representation of the NSFWMode
func (n NSFWMode) String() string {
	switch n {
	case NSFWModeUnspecified:
		return "unspecified"
	case NSFWModeBlock:
		return "block"
	case NSFWModeAccept:
		return "accept"
	case NSFWModeOnly:
		return "only"
	default:
		return "unknown"
	}
}