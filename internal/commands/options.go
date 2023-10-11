package commands

import (
	"aws-ssh/internal/commands/factory"

	"github.com/spf13/cobra"
)

type Options interface {
	Validate(f factory.Factory, cmd *cobra.Command, args []string) error
	Run(f factory.Factory, cmd *cobra.Command, args []string) error
}
