package hermesconfig

type Operation struct {
	// List of task keys to run as part of this operation
	Tasks []string
	// The description to be used in the CLI help output
	Description string
}
