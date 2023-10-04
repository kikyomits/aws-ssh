package prerun

import (
	"aws-ssh/internal/commands/flags"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Setup(cmd *cobra.Command, args []string) {
	v, err := cmd.Flags().GetBool(flags.VerboseFlag)
	if err != nil {
		log.WithError(err).Errorf("Unknown error to parse verbose flag")
	}
	isVerbose := err == nil && v

	setupLogger(isVerbose)
}

func setupLogger(verbose bool) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
		FieldMap: log.FieldMap{
			log.FieldKeyTime: "timestamp",
			log.FieldKeyMsg:  "message",
		},
	})
}
