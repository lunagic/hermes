package hermes

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/aaronellington/hermes/hermes/hermesconfig"
)

type Service struct {
	hermes hermesconfig.Config
	logger Logger
}

func New(config hermesconfig.Config) (*Service, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Service{
		logger: loggerLogger{
			Logger: log.New(os.Stdout, "", log.LstdFlags),
		},
		hermes: config,
	}, nil
}

func (service *Service) Execute(operationName string) error {
	operation := service.hermes.Operations[operationName]
	anyFailure := false

	for _, taskName := range operation.Tasks {
		service.logger.Log(fmt.Sprintf("TASK: %s", taskName))
		task := service.hermes.Tasks[taskName]
		hosts := service.hermes.GetHostsByTag(task.HostTags...)

		hostWaitGroup := &sync.WaitGroup{}
		for _, host := range hosts {
			hostWaitGroup.Add(1)
			go func(host hermesconfig.Host) {
				defer hostWaitGroup.Done()
				if !service.runTask(host, task) {
					anyFailure = true
				}
			}(host)
		}
		hostWaitGroup.Wait()
	}

	if anyFailure {
		return errors.New("was not completely successful")
	}

	return nil
}

func (service *Service) runTask(host hermesconfig.Host, task hermesconfig.Task) bool {
	if task.If != "" {
		if err := service.shellExec(host, task, task.If); err != nil {
			service.logger.Info(fmt.Sprintf("SKIPPED: %s", host.Hostname))
			return true
		}
	}

	if err := service.shellExec(host, task, strings.Join(task.Commands, " && ")); err != nil {
		service.logger.Error(fmt.Sprintf("ERROR: %s %s", host.Hostname, err.Error()))
		return false
	}

	service.logger.Success(fmt.Sprintf("SUCCESS: %s", host.Hostname))

	return true
}

func (service *Service) shellExec(host hermesconfig.Host, task hermesconfig.Task, command string) error {
	destination := host.Hostname
	if task.User != "" {
		destination = task.User + "@" + host.Hostname
	}

	errorBuffer := bytes.NewBuffer([]byte{})

	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-o", "ConnectTimeout=5", destination, command)
	cmd.Stderr = errorBuffer

	if err := cmd.Run(); err != nil {
		return errors.New(errorBuffer.String())
	}

	return nil
}
