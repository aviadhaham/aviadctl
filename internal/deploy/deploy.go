package deploy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

var Cmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy 2 pods in the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
		if err != nil {
			log.Fatalf("Error loading kubeconfig: %v", err)
		}

		dynamicClient, err := dynamic.NewForConfig(config)
		if err != nil {
			log.Fatalf("Error creating clientset: %v", err)
		}

		manifestDir := "manifests"

		err = applyManifests(dynamicClient, manifestDir)
		if err != nil {
			log.Fatalf("Error applying manifests: %v", err)
		}

		fmt.Println("Manifests applied successfully")
	},
}

func applyManifests(dynamicClient dynamic.Interface, manifestDir string) error {
	err := filepath.Walk(manifestDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".yml" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %v", path, err)
		}

		decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), 1<<23)
		for {
			obj := &unstructured.Unstructured{}
			err = decoder.Decode(obj)
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("error decoding YAML from %s: %v", path, err)
			}

			// Skip empty objects
			if len(obj.Object) == 0 {
				continue
			}

			gvk := obj.GroupVersionKind()
			gvr := schema.GroupVersionResource{
				Group:    gvk.Group,
				Version:  gvk.Version,
				Resource: strings.ToLower(gvk.Kind) + "s",
			}

			ns := obj.GetNamespace()
			if ns == "" {
				ns = "default"
			}

			// Apply the object
			_, err = dynamicClient.Resource(gvr).Namespace(ns).Create(context.TODO(), obj, metav1.CreateOptions{})
			if err != nil {
				if errors.IsAlreadyExists(err) {
					_, err = dynamicClient.Resource(gvr).Namespace(obj.GetNamespace()).Update(context.TODO(), obj, metav1.UpdateOptions{})
					if err != nil {
						return fmt.Errorf("error updating resource from %s: %v", path, err)
					}
					fmt.Printf("Updated resource: %s/%s\n", obj.GetKind(), obj.GetName())
				} else {
					return fmt.Errorf("error creating resource from %s: %v", path, err)
				}
			} else {
				fmt.Printf("Created resource: %s/%s\n", obj.GetKind(), obj.GetName())
			}
		}
		return nil
	})

	return err
}
