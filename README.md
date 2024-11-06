## Steps to Run the Exercise

1. Build the binary:

   ```bash
   go build -o aviadctl
   ```

2. Run the cluster setup:

   ```bash
   ./aviadctl cluster
   ```

   **_After running `cluster`, make sure you set up the kubeconfig file._**

   ```bash
   mkdir ~/.kube
   sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
   export KUBECONFIG=~/.kube/config
   ```

3. Wait for the cluster to be completely initialized, and then run:

   ```bash
   ./aviadctl deploy
   ```

4. Check the status of the pods:
   ```bash
   ./aviadctl status
   ```

## Notes

- `KUBECONFIG` environment variable must be set and must be accessible by the user running the `aviadctl` binary.
- You can also `go install` the `aviadctl` binary and add it to your PATH:
  ```bash
  go install github.com/aviadhaham/aviadctl@latest
  ```
  Then you can run it from anywhere:
  ```bash
  aviadctl <command>
  ```
- To access the Wordpress UI running in the deployed pod, you can utilize `kubectl port-forward` (the container port is `8080`).
- To access the Wordpress Admin UI, you can reach `/wp-admin`.
- To get the password for the admin user (named `user`), use:
  ```bash
  kubectl get secret --namespace default wp-wordpress -o jsonpath="{.data.wordpress-password}" | base64 --decode
  ```
