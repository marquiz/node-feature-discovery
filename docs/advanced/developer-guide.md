---
title: "Developer Guide"
layout: default
parent: Advanced
nav_order: 1
---

# Developer Guide
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Building from Source

**Download the source code:**

```
git clone https://github.com/kubernetes-sigs/node-feature-discovery
cd node-feature-discovery
```

### Docker Build


**Build the container image:**<br>
See [customizing the build](#customizing-the-build) below for altering the
container image registry, for example.

```
make
```

**Push the container image:**<br>
Optional, this example with Docker.

```
docker push <IMAGE_TAG>
```

**Change the job spec to use your custom image (optional):**

To use your published image from the step above instead of the
`k8s.gcr.io/nfd/node-feature-discovery` image, edit `image`
attribute in the spec template(s) to the new location
(`<registry-name>/<image-name>[:<version>]`).

### Building Locally

You can also build the binaries locally
```
make build
```

This will compile binaries under `bin/`

### Customizing the Build

There are several Makefile variables that control the build process and the
name of the resulting container image. The following are targeted targeted for
build customization and they can be specified via environment variables or
makefile overrides.

| Variable                   | Description                                                       | Default value
| -------------------------- | ----------------------------------------------------------------- | ----------- |
| HOSTMOUNT_PREFIX           | Prefix of system directories for feature discovery (local builds) | / (*local builds*) /host- (*container builds*)
| IMAGE_BUILD_CMD            | Command to build the image                                        | docker build
| IMAGE_BUILD_EXTRA_OPTS     | Extra options to pass to build command                            | *empty*
| IMAGE_PUSH_CMD             | Command to push the image to remote registry                      | docker push
| IMAGE_REGISTRY             | Container image registry to use                                   | k8s.gcr.io/nfd
| IMAGE_TAG_NAME             | Container image tag name                                          | &lt;nfd version&gt;
| IMAGE_EXTRA_TAG_NAMES      | Additional container image tag(s) to create when building image   | *empty*
| K8S_NAMESPACE              | nfd-master and nfd-worker namespace                               | kube-system
| KUBECONFIG                 | Kubeconfig for running e2e-tests                                  | *empty*
| E2E_TEST_CONFIG            | Parameterization file of e2e-tests (see [example](test/e2e/e2e-test-config.exapmle.yaml)) | *empty*

For example, to use a custom registry:

```
make IMAGE_REGISTRY=<my custom registry uri>
```

Or to specify a build tool different from Docker, It can be done in 2 ways:
1. via environment
```
IMAGE_BUILD_CMD="buildah bud" make
```
2. by overriding the variable value
```
make  IMAGE_BUILD_CMD="buildah bud"
```

### Testing

Unit tests are automatically run as part of the container image build. You can
also run them manually in the source code tree by simply running:
```
make test
```

End-to-end tests are built on top of the e2e test framework of Kubernetes, and,
they required a cluster to run them on. For running the tests on your test
cluster you need to specify the kubeconfig to be used:
```
make e2e-test KUBECONFIG=$HOME/.kube/config
```


## Running Locally

*WORK IN PROGRESS...*


## Documentation

*WORK IN PROGRESS...*
