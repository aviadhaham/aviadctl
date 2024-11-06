package cluster

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

const (
	clusterInstallCommand = "curl -sfL https://get.k3s.io | sh -"
)

var Cmd = &cobra.Command{
	Use:   "cluster",
	Short: "Deploy a k3s k8s cluster",
	Run: func(cmd *cobra.Command, args []string) {
		cmdCheckK3s := exec.Command("/bin/sh", "-c", "which k3s")
		err := cmdCheckK3s.Run()
		if err == nil {
			fmt.Println("k3s is already installed. Exiting.")
		}

		cmdInstallK3s := exec.Command("/bin/sh", "-c", clusterInstallCommand)
		out, err := cmdInstallK3s.CombinedOutput()
		if err != nil {
			log.Fatalf("failed with error: %v, output: %s", err, string(out))
		}
	},
}
