# Demo Broker

Demo broker is an implementation of the [Open Service
Broker API](https://github.com/openservicebrokerapi/servicebroker) and is based on
[`osb-starter-pack`](https://github.com/pmorie/osb-starter-pack).

The catalog is retrieved from a remote address. By default:

``
$ export CATALOG_PATH=https://raw.githubusercontent.com/cheld/demo-broker/master/samples/catalog.json
``

Customize as needed. This functionality can be used to
* simulate a catalog look & feel.
* implement the a REST service that generates this JSON as a first implementation step.


## Prerequisites

You'll need:

- [`go`](https://golang.org/dl/)
- A running [Kubernetes](https://github.com/kubernetes/kubernetes) (or [openshift](https://github.com/openshift/origin/)) cluster
- The [service-catalog](https://github.com/kubernetes-incubator/service-catalog)
  [installed](https://github.com/kubernetes-incubator/service-catalog/blob/master/docs/install.md)
  in that cluster


## Getting started

You can `go get` this repo or `git clone` it to start poking around right away.

The project comes ready with a minimal example service that you can easily
deploy and begin iterating on.

### Get the project

```console
$ go get github.com/cheld/demo-broker/cmd/servicebroker
```

Or clone the repo:

```console
$ cd $GOPATH/src && mkdir -p github.com/cheld && cd github.com/cheld && git clone git://github.com/cheld/demo-broker
```

Change into the project directory:

```console
$ cd $GOPATH/src/github.com/cheld/demo-broker
```

Deploy to OpenShift cluster by passing a custom image and tag name.
Note: You must already be logged into an OpenShift cluster.
This also pushes the generated image with docker.

```console
$ IMAGE=myimage TAG=latest make push deploy-openshift
```

Running either of these flavors of deploy targets will build the broker binary,
build the image, deploy the broker into your Kubernetes, and add a
`ClusterServiceBroker` to the service-catalog.
