# Node feature discovery for [Kubernetes](https://kubernetes.io)

[![Go Report Card](https://goreportcard.com/badge/github.com/kubernetes-sigs/node-feature-discovery)](https://goreportcard.com/report/github.com/kubernetes-sigs/node-feature-discovery)
[![Prow Build](https://prow.k8s.io/badge.svg?jobs=pull-node-feature-discovery-build-image)](https://prow.k8s.io/job-history/gs/kubernetes-jenkins/pr-logs/directory/pull-node-feature-discovery-build-image)

- [Command line interface](#command-line-interface)
- [Getting started](#getting-started)
  - [System requirements](#system-requirements)
  - [Usage](#usage)
- [Targeting nodes with specific features](#targeting-nodes-with-specific-features)
- [Community](#community)
- [License](#license)
- [Demo](#demo)

## Command line interface

You can run NFD in stand-alone Docker containers e.g. for testing
purposes. This is useful for checking features-detection.

### NFD-Master

When running as a standalone container labeling is expected to fail because
Kubernetes API is not available. Thus, it is recommended to use `--no-publish`
command line flag. E.g.
```
$ docker run --rm --name=nfd-test <NFD_CONTAINER_IMAGE> nfd-master --no-publish
2019/02/01 14:48:21 Node Feature Discovery Master <NFD_VERSION>
2019/02/01 14:48:21 gRPC server serving on port: 8080
```

Command line flags of nfd-master:
```
$ docker run --rm <NFD_CONTAINER_IMAGE> nfd-master --help
...
Usage:
  nfd-master [--prune] [--no-publish] [--label-whitelist=<pattern>] [--port=<port>]
     [--ca-file=<path>] [--cert-file=<path>] [--key-file=<path>]
     [--verify-node-name] [--extra-label-ns=<list>] [--resource-labels=<list>]
     [--kubeconfig=<path>]
  nfd-master -h | --help
  nfd-master --version

  Options:
  -h --help                       Show this screen.
  --version                       Output version and exit.
  --prune                         Prune all NFD related attributes from all nodes
                                  of the cluster and exit.
  --kubeconfig=<path>             Kubeconfig to use [Default: ]
  --port=<port>                   Port on which to listen for connections.
                                  [Default: 8080]
  --ca-file=<path>                Root certificate for verifying connections
                                  [Default: ]
  --cert-file=<path>              Certificate used for authenticating connections
                                  [Default: ]
  --key-file=<path>               Private key matching --cert-file
                                  [Default: ]
  --verify-node-name              Verify worker node name against CN from the TLS
                                  certificate. Only has effect when TLS authentication
                                  has been enabled.
  --no-publish                    Do not publish feature labels
  --label-whitelist=<pattern>     Regular expression to filter label names to
                                  publish to the Kubernetes API server.
                                  NB: the label namespace is omitted i.e. the filter
                                  is only applied to the name part after '/'.
                                  [Default: ]
  --extra-label-ns=<list>         Comma separated list of allowed extra label namespaces
                                  [Default: ]
  --resource-labels=<list>        Comma separated list of labels to be exposed as extended resources.
                                  [Default: ]
```

### NFD-Worker

In order to run nfd-worker as a "stand-alone" container against your
standalone nfd-master you need to run them in the same network namespace:
```
$ docker run --rm --network=container:nfd-test <NFD_CONTAINER_IMAGE> nfd-worker
2019/02/01 14:48:56 Node Feature Discovery Worker <NFD_VERSION>
...
```
If you just want to try out feature discovery without connecting to nfd-master,
pass the `--no-publish` flag to nfd-worker.

Command line flags of nfd-worker:
```
$ docker run --rm <CONTAINER_IMAGE_ID> nfd-worker --help
...
nfd-worker.

  Usage:
  nfd-worker [--no-publish] [--sources=<sources>] [--label-whitelist=<pattern>]
     [--oneshot | --sleep-interval=<seconds>] [--config=<path>]
     [--options=<config>] [--server=<server>] [--server-name-override=<name>]
     [--ca-file=<path>] [--cert-file=<path>] [--key-file=<path>]
  nfd-worker -h | --help
  nfd-worker --version

  Options:
  -h --help                   Show this screen.
  --version                   Output version and exit.
  --config=<path>             Config file to use.
                              [Default: /etc/kubernetes/node-feature-discovery/nfd-worker.conf]
  --options=<config>          Specify config options from command line. Config
                              options are specified in the same format as in the
                              config file (i.e. json or yaml). These options
                              will override settings read from the config file.
                              [Default: ]
  --ca-file=<path>            Root certificate for verifying connections
                              [Default: ]
  --cert-file=<path>          Certificate used for authenticating connections
                              [Default: ]
  --key-file=<path>           Private key matching --cert-file
                              [Default: ]
  --server=<server>           NFD server address to connecto to.
                              [Default: localhost:8080]
  --server-name-override=<name> Name (CN) expect from server certificate, useful
                              in testing
                              [Default: ]
  --sources=<sources>         Comma separated list of feature sources.
                              [Default: cpu,custom,iommu,kernel,local,memory,network,pci,storage,system,usb]
  --no-publish                Do not publish discovered features to the
                              cluster-local Kubernetes API server.
  --label-whitelist=<pattern> Regular expression to filter label names to
                              publish to the Kubernetes API server.
                              NB: the label namespace is omitted i.e. the filter
                              is only applied to the name part after '/'.
                              [Default: ]
  --oneshot                   Label once and exit.
  --sleep-interval=<seconds>  Time to sleep between re-labeling. Non-positive
                              value implies no re-labeling (i.e. infinite
                              sleep). [Default: 60s]
```
**NOTE** Some feature sources need certain directories and/or files from the
host mounted inside the NFD container. Thus, you need to provide Docker with the
correct `--volume` options in order for them to work correctly when run
stand-alone directly with `docker run`. See the
[template spec](https://github.com/kubernetes-sigs/node-feature-discovery/blob/master/nfd-worker-daemonset.yaml.template)
for up-to-date information about the required volume mounts.

## Getting started

For a stable version with ready-built images see the
[latest released version](https://github.com/kubernetes-sigs/node-feature-discovery/tree/v0.6.0) ([release notes](https://github.com/kubernetes-sigs/node-feature-discovery/releases/latest)).

If you want to use the latest development version (master branch) you need to
[build your own custom image](#building-from-source).

### System requirements

1. Linux (x86_64/Arm64/Arm)
1. [kubectl][kubectl-setup] (properly set up and configured to work with your
   Kubernetes cluster)
1. [Docker][docker-down] (only required to build and push docker images)

### Usage

#### nfd-master

Nfd-master runs as a deployment (with a replica count of 1), by default
it prefers running on the cluster's master nodes but will run on worker
nodes if no master nodes are found.

For High Availability, you should simply increase the replica count of
the deployment object. You should also look into adding [inter-pod](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity)
affinity to prevent masters from running on the same node.
However note that inter-pod affinity is costly and is not recommended
in bigger clusters.

You can use the template spec provided to deploy nfd-master, or
use `nfd-master.yaml` generated by `Makefile`. The latter includes
`image:` and `namespace:` definitions that match the latest built
image. Example:
```
make IMAGE_TAG=<IMAGE_TAG>
docker push <IMAGE_TAG>
kubectl create -f nfd-master.yaml
```
Nfd-master listens for connections from nfd-worker(s) and connects to the
Kubernetes API server to add node labels advertised by them.

If you have RBAC authorization enabled (as is the default e.g. with clusters
initialized with kubeadm) you need to configure the appropriate ClusterRoles,
ClusterRoleBindings and a ServiceAccount in order for NFD to create node
labels. The provided template will configure these for you.


#### nfd-worker

Nfd-worker is preferably run as a Kubernetes DaemonSet. There is an
example spec (`nfd-worker-daemonset.yaml.template`) that can be used
as a template, or, as is when just trying out the service. Similarly
to nfd-master above, the `Makefile` also generates
`nfd-worker-daemonset.yaml` from the template that you can use to
deploy the latest image. Example:
```
make IMAGE_TAG=<IMAGE_TAG>
docker push <IMAGE_TAG>
kubectl create -f nfd-worker-daemonset.yaml
```

Nfd-worker connects to the nfd-master service to advertise hardware features.

When run as a daemonset, nodes are re-labeled at an interval specified using
the `--sleep-interval` option. In the
[template](https://github.com/kubernetes-sigs/node-feature-discovery/blob/master/nfd-worker-daemonset.yaml.template#L26)
the default interval is set to 60s which is also the default when no
`--sleep-interval` is specified. Also, the configuration file is re-read on
each iteration providing a simple mechanism of run-time reconfiguration.

Feature discovery can alternatively be configured as a one-shot job. There is
an example script in this repo that demonstrates how to deploy the job in the cluster.

```
./label-nodes.sh [<IMAGE_TAG>]
```

The label-nodes.sh script tries to launch as many jobs as there are Ready nodes.
Note that this approach does not guarantee running once on every node.
For example, if some node is tainted NoSchedule or fails to start a job for some other reason, then some other node will run extra job instance(s) to satisfy the request and the tainted/failed node does not get labeled.

#### nfd-master and nfd-worker in the same Pod

You can also run nfd-master and nfd-worker inside a single pod (skip the `sed`
part if running the latest released version):
```
sed -E s',^(\s*)image:.+$,\1image: <YOUR_IMAGE_REPO>:<YOUR_IMAGE_TAG>,' nfd-daemonset-combined.yaml.template > nfd-daemonset-combined.yaml
kubectl apply -f nfd-daemonset-combined.yaml
```
Similar to the nfd-worker setup above, this creates a DaemonSet that schedules
an NFD Pod an all worker nodes, with the difference that the Pod also also
contains an nfd-master instance. In this case no nfd-master service is run on
the master node(s), but, the worker nodes are able to label themselves.

This may be desirable e.g. in single-node setups.

#### TLS authentication

NFD supports mutual TLS authentication between the nfd-master and nfd-worker
instances.  That is, nfd-worker and nfd-master both verify that the other end
presents a valid certificate.

TLS authentication is enabled by specifying `--ca-file`, `--key-file` and
`--cert-file` args, on both the nfd-master and nfd-worker instances.
The template specs provided with NFD contain (commented out) example
configuration for enabling TLS authentication.

The Common Name (CN) of the nfd-master certificate must match the DNS name of
the nfd-master Service of the cluster. By default, nfd-master only check that
the nfd-worker has been signed by the specified root certificate (--ca-file).
Additional hardening can be enabled by specifying --verify-node-name in
nfd-master args, in which case nfd-master verifies that the NodeName presented
by nfd-worker matches the Common Name (CN) of its certificate. This means that
each nfd-worker requires a individual node-specific TLS certificate.


#### Usage demo

[![asciicast](https://asciinema.org/a/247316.svg)](https://asciinema.org/a/247316)

### Configuration options

Nfd-worker supports a configuration file. The default location is
`/etc/kubernetes/node-feature-discovery/nfd-worker.conf`, but,
this can be changed by specifying the`--config` command line flag.
Configuration file is re-read on each labeling pass (determined by
`--sleep-interval`) which makes run-time re-configuration of nfd-worker
possible.

Worker configuration file is read inside the container, and thus, Volumes and
VolumeMounts are needed to make your configuration available for NFD. The
preferred method is to use a ConfigMap which provides easy deployment and
re-configurability.  For example, create a config map using the example config
as a template:
```
cp nfd-worker.conf.example nfd-worker.conf
vim nfd-worker.conf  # edit the configuration
kubectl create configmap nfd-worker-config --from-file=nfd-worker.conf
```
Then, configure Volumes and VolumeMounts in the Pod spec (just the relevant
snippets shown below):
```
...
  containers:
      volumeMounts:
        - name: nfd-worker-config
          mountPath: "/etc/kubernetes/node-feature-discovery/"
...
  volumes:
    - name: nfd-worker-config
      configMap:
        name: nfd-worker-config
...
```
You could also use other types of volumes, of course. That is, hostPath if
different config for different nodes would be required, for example.

The (empty-by-default)
[example config](https://github.com/kubernetes-sigs/node-feature-discovery/blob/master/nfd-worker.conf.example)
is used as a config in the NFD Docker image. Thus, this can be used as a default
configuration in custom-built images.

Configuration options can also be specified via the `--options` command line
flag, in which case no mounts need to be used. The same format as in the config
file must be used, i.e. JSON (or YAML). For example:
```
--options='{"sources": { "pci": { "deviceClassWhitelist": ["12"] } } }'
```
Configuration options specified from the command line will override those read
from the config file.

Currently, the only available configuration options are related to the
[CPU](#cpu-features), [PCI](#pci-features) and [Kernel](#kernel-features)
feature sources.

## Targeting Nodes with Specific Features

Nodes with specific features can be targeted using the `nodeSelector` field. The
following example shows how to target nodes with Intel TurboBoost enabled.

```yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    env: test
  name: golang-test
spec:
  containers:
    - image: golang
      name: go1
  nodeSelector:
    feature.node.kubernetes.io/cpu-pstate.turbo: 'true'
```

For more details on targeting nodes, see [node selection][node-sel].

## Community

You can reach us via the following channels:

- [#node-feature-discovery](https://kubernetes.slack.com/messages/node-feature-discovery) channel in
  [Kubernetes Slack](slack.k8s.io)
- [SIG-Node](https://groups.google.com/g/kubernetes-sig-node) mailing list
- File an [issue](https://github.com/kubernetes-sigs/node-feature-discovery/issues/new) in this repository


## Governance

This is a [SIG-node](https://github.com/kubernetes/community/blob/master/sig-node/README.md)
subproject, hosted under the
[Kubernetes SIGs](https://github.com/kubernetes-sigs) organization in
Github. The project was established in 2016 as a
[Kubernetes Incubator](https://github.com/kubernetes/community/blob/master/incubator.md)
project and migrated to Kubernetes SIGs in 2018.

## License

This is open source software released under the [Apache 2.0 License](LICENSE).

## Demo

A demo on the benefits of using node feature discovery can be found in [demo](demo/).

<!-- Links -->
[docker-down]: https://docs.docker.com/install
[kubectl-setup]: https://kubernetes.io/docs/tasks/tools/install-kubectl
[node-sel]: http://kubernetes.io/docs/user-guide/node-selection
