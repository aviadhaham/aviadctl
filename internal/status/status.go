package status

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	NAMESPACE = "default"
)

var Cmd = &cobra.Command{
	Use:   "status",
	Short: "Print out status table containing the status of pod names and status in default namespace",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
		if err != nil {
			log.Fatalf("Error loading kubeconfig: %v", err)
		}

		clienset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("Error creating clientset: %v", err)
		}

		podList, err := clienset.CoreV1().Pods(NAMESPACE).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error creating watch: %v", err)
		}

		tbl := table.New("Namespace", "Name", "Status")
		counter := 0
		for _, pod := range podList.Items {
			counter++
			tbl.AddRow(counter, pod.Name, pod.Status.Phase)
			fmt.Println()
		}
		tbl.Print()

	},
}
