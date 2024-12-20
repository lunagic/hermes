package hermesconfig

type Host struct {
	// Name of the host to be used in the SSH comment (note: can be named entries in your SSH config)
	Hostname string
	// Tags to be matched with tasks to target to specific hosts
	Tags []string
}

func (host Host) HasAnyTags(tags ...string) bool {
	for _, tag := range tags {
		for _, hostTag := range host.Tags {
			if tag == hostTag {
				return true
			}
		}
	}

	return false
}
