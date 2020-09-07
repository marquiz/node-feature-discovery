---
title: "Worker Cmdline Reference"
layout: default
parent: Advanced
nav_order: 4
---

# NFD-Worker Commandline Flags
{: .no_toc }

## Table of Contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

To quickly view available command line flags execute `nfd-worker --help`.
In a docker container:
```
$ docker run k8s.gcr.io/nfd/node-feature-discovery:v0.6.0 nfd-worker --help
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
