## kreate

`kreate` is a command line utility that can be used to create the Kubernetes
resources that are not support by `kubectl create` command.

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

> [!NOTE]  
> The code structure and the entire idea is just copied from `kubectl` project.

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

## Installing `kreate`

`kreate` is not available currently using OS package managers. You will have to install it by
downloading the released binary from GitHub.

Go to releases page and download the appropriate binary for your operating system and architecture,
using either curl or wget commands. And move it to your `PATH`.

You can figure out the operating system details using the below command

```
uname -a
Darwin vsap-mac-HW32YFVFWC 23.6.0 Darwin Kernel Version 23.6.0: Mon Jul 29 21:14:30 PDT 2024; root:xnu-10063.141.2~1/RELEASE_ARM64_T6030 arm64
```

### MacOS

- Download, respective binary (artifacts) from releases, you can specify the expected value for the version var

```
export VERSION=0.0.1
wget https://github.com/viveksinghggits/kreate/releases/download/v${VERSION}/kreate_Darwin_arm64.tar.gz
```

- Extract the downloaded .tar.gz file

```
tar xf kreate_Darwin_arm64.tar.gz
```

- Move the binary to PATH

```
mv kreate /usr/local/bin
```

## Using with kubectl

Since `kreate` is a standalone binary it can easily be used with kubectl as kubectl plugin.
More about kubectl plugins can be read [here](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/).

Download the `kreate` binary for your OS and ARCH and move the binary to the `$PATH` after renaming
it to `kubectl-kreate`. And that's all, now you can easily use `kreate` with standard `kubectl` tool
like shown below.

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
