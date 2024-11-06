package cmd

import (
	"os"

	"github.com/aviadhaham/aviadctl/internal/cluster"
	"github.com/aviadhaham/aviadctl/internal/deploy"
	"github.com/aviadhaham/aviadctl/internal/status"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aviadctl",
	Short: "The aviadctl command",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(cluster.Cmd)
	rootCmd.AddCommand(deploy.Cmd)
	rootCmd.AddCommand(status.Cmd)
}
