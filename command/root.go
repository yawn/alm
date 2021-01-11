package command

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           app,
	SilenceErrors: true,
	SilenceUsage:  true,
}
