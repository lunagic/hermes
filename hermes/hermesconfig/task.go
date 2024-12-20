package hermesconfig

type Task struct {
	// User override for the SSH connection
	User string
	// If set, this command must succeed in order to run the commands on the given host
	If string
	// This task will execute on Any hosts that match any of these tags. If this list is empty, all hosts are included
	HostTags []string `yaml:"host_tags"`
	// All commands are executed in order and stops if any fail
	Commands []string
}
