---
title: "Master Cmdline Reference"
layout: default
parent: Advanced
nav_order: 3
---

# NFD-Master Commandline Flags
{: .no_toc }

## Table of Contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

To quickly view available command line flags execute `nfd-master --help`.
In a docker container:
```
$ docker run k8s.gcr.io/nfd/node-feature-discovery:v0.6.0 nfd-master --help
nfd-master.

  Usage:
  nfd-master [--no-publish] [--label-whitelist=<pattern>] [--port=<port>]
     [--ca-file=<path>] [--cert-file=<path>] [--key-file=<path>]
     [--verify-node-name] [--extra-label-ns=<list>] [--resource-labels=<list>]
  nfd-master -h | --help
  nfd-master --version

  Options:
  -h --help                       Show this screen.
  --version                       Output version and exit.
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

### `--version

Output version and exit.
