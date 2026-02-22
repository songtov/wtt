package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set via ldflags at build time: -ldflags "-X github.com/songtov/wtt/cmd.Version=v1.0.0"
var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Args:  cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println(Version)
	},
}
