package hermesconfig

import "fmt"

type Config struct {
	Hosts      []Host
	Tasks      map[string]Task
	Operations map[string]Operation
}

func (config Config) Validate() error {
	for operationName, operation := range config.Operations {
		for _, taskName := range operation.Tasks {
			if _, found := config.Tasks[taskName]; !found {
				return fmt.Errorf("unknown task '%s' in operation '%s'", taskName, operationName)
			}
		}
	}

	return nil
}

func (config Config) GetHostsByTag(tags ...string) []Host {
	if len(tags) == 0 {
		return config.Hosts
	}

	hosts := []Host{}
	for _, host := range config.Hosts {
		if !host.HasAnyTags(tags...) {
			continue
		}

		hosts = append(hosts, host)
	}

	return hosts
}
