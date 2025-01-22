package cli

import (
	"log"
	"os"

	"github.com/lunagic/hermes/hermes"
	"github.com/lunagic/hermes/hermes/hermesconfig"
	"github.com/lunagic/environment-go/environment"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	WorkingDirectory string `env:"HERMES_WORKING_DIRECTORY"`
}

func Cmd() *cobra.Command {
	cliConfig := &Config{
		WorkingDirectory: ".",
	}
	if err := environment.New().Decode(cliConfig); err != nil {
		log.Fatal(err)
	}

	// Change working directory to the one listed in the configuration only
	// if there isn't a file already in our working directory.
	if _, err := os.Stat("hermes.yaml"); err != nil {
		if err := os.Chdir(cliConfig.WorkingDirectory); err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.OpenFile("hermes.yaml", os.O_RDONLY, 0770)
	if err != nil {
		log.Fatal(err)
	}

	config := hermesconfig.Config{}
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal(err)
	}

	root := &cobra.Command{
		Use:   "hermes",
		Short: "Shell Command Delivery Tool",
	}

	run := &cobra.Command{
		Use:   "run",
		Short: "Run an 'operation' found in the hermes.yaml file",
	}
	root.AddCommand(run)

	for operationName, operation := range config.Operations {
		run.AddCommand(&cobra.Command{
			Use:          operationName,
			Short:        operation.Description,
			SilenceUsage: true,
			RunE: func(cmd *cobra.Command, args []string) error {
				h, err := hermes.New(config)
				if err != nil {
					log.Fatal(err)
				}

				return h.Execute(operationName)
			},
		})
	}

	return root
}
