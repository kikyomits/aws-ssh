package commands

import (
	"aws-ssh/internal/commands/ecs_port_forward"
	"aws-ssh/internal/commands/factory"
	"aws-ssh/internal/commands/flags"
	"aws-ssh/internal/commands/prerun"

	"github.com/spf13/cobra"
)

// RegisteredCommands returns a realized mapping of available CLI commands in a format that
// the CLI class can consume.
func RegisteredCommands() (cobra.Command, error) {
	root := cobra.Command{
		Use:   "aws-ssh",
		Short: "developer friendly tool to ssh or port forward with AWS ECS",
		Long: "aws-ssh helps developers to execute commands inside Amazon Elastic Container Service (ECS) containers \n" +
			"and set up port forwarding to remote hosts or ports. It simplifies the process of interacting with containers \n" +
			"running in ECS clusters, making it easy to manage and troubleshoot containerized applications.",
		PersistentPreRun: prerun.Setup,
		SilenceUsage:     true,
		SilenceErrors:    true,
	}

	root.AddCommand(ecs_port_forward.New(&factory.AWSFactory{}))
	root.PersistentFlags().String(flags.AWSProfileFlag, "", "AWS Profile")
	root.PersistentFlags().String(flags.RegionFlag, "", "AWS Region")
	root.PersistentFlags().BoolP(flags.VerboseFlag, "v", false, "enable verbose logging")
	return root, nil
}
