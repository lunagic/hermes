package hermesconfig

type Job struct {
	Name     string
	HostTags []string `yaml:"host_tags"`
	Tasks    []string
}
