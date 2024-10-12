## kreate

`kreate` a command line utility that can be used to create the Kubernetes
resources that are not support by `kubectl create` command.

Initially we were planning to just support generating the `yaml` representation
of the resource using `-oyaml` argument so that it can be changed manually and
then the resource can be created.
But eventually settled on just supporting creation as well.

```bash
./kreate --help
Create Kubernetes resources that are not supported by `kubectl create`

Usage:
  kreate [flags]
  kreate [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  pvc         Create PVC resources

Flags:
  -h, --help   help for kreate

Use "kreate [command] --help" for more information about a command.
```

## Examples

- Create a PVC resource yaml with PVC name `my-claim`

```
./kreate pvc my-claim --dry-run=client -oyaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  creationTimestamp: null
  name: my-claim
  namespace: default
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: ""
status: {}
```

- Create a PVC resource yaml with PVC name `my-claim`, storage class `gp2-csi` and size 20Gi

```
./kreate pvc my-claim --size 20 --storageclass gp2-csi  --dry-run=client -oyaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  creationTimestamp: null
  name: my-claim
  namespace: default
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 20Gi
  storageClassName: gp2-csi
status: {}
```

## Using with kubectl

Since `kreate` is a standalone binary it can easily be used with kubectl as kubectl plugin.
More about kubectl plugins can be read [here](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/).

Download the `kreate` binary for your OS and ARCH and move the binary to the `$PATH` after renaming
it to `kubectl-kreate`. And that's all, now you can easily use `kreate` with standard `kubectl` tool.

```
kubectl kreate pvc my-claim --size 30 --dry-run=client -oyaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  creationTimestamp: null
  name: my-claim
  namespace: default
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 30Gi
  storageClassName: ""
status: {}
```


## Supported Resources

- [PersistentVolumeClaim](https://pkg.go.dev/k8s.io/api/core/v1#PersistentVolumeClaim)


> [!NOTE]  
> The code structure and the entire idea is just copied from `kubectl` project.
