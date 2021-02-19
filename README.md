# jobless

Simple Kubernetes job watcher to cleanup completed jobs

It watches for successfully completed jobs (`status.successful=1`) in the namespace in which it's deployed and simply deletes the jobs

## Usage

```sh
kubectl apply -f deployment.yaml
```
