# Kubernetes Generic Controller

The Kubernetes Generic Controller is a controller written in Go. The current version uses the client-go library in version 1.17.0.
A configuration file named config.yaml allows the selective control of Kubernetes resources and if needed, own handle functions can be defined for these resources. The handler defines how the controller should behave when an object is added, updated or deleted.

The controller can run outside of the cluster or on the cluster as a Pod.



## Usage

### Adding custom handler functions

The custom handler has to implement the following interface which is defined in handler.go:

```go
type Handler interface {
	Init(c Config) error
	ObjectCreated(obj interface{})
	ObjectDeleted(obj interface{})
	ObjectUpdated(oldObj, newObj interface{})
}
```

A handler might look something like this (sample_handler.go):
```go
package main

import "log"

type Sample struct {
}

func (t *Sample) Init(c Config) error {

	return nil
}

func (t *Sample) ObjectCreated(obj interface{}) {
	logrus.Info("Sample CREATE function invoked")
}

func (t *Sample) ObjectUpdated(oldObj, newObj interface{}) {
	logrus.Info("Sample UPDATE function invoked")

}

func (t *Sample) ObjectDeleted(obj interface{}) {
	logrus.Info("Sample DELETE function invoked")

}
```

The new custom Type has to be added to config.go. When the config.yaml file is parsed it detects if a handler has been defined. If that's the case the custom handler will be used, if not the default handler will take care of the object:
```go
func ParseEventHandler(conf Config) Handler {

	var eventHandler Handler
	switch {
	case conf.Resources.Rolebinding:
		// Add custom Type here
		eventHandler = new(Sample)
	default:
		eventHandler = new(Default)
	}
	if err := eventHandler.Init(conf); err != nil {
		log.Fatal(err)
	}
	return eventHandler
}
```
### Configuration

The following config would enable a watch for Pods and Deployments. The Sample handler would take care of these objects:
```yaml
resources:
  pod: true
  deployment: true
  replicationcontroller: false
  replicaset: false
  daemonset: false
  services: false
  secret: false
  configmap: false
  rolebinding: false
handler:
  name: Sample
```

### Running the Controller with an out-of-cluster-config:
```bash
export KUBECONFIG=/path/to/config
./kubernetes-controller
```

### Running the Controller with an in-cluster-config:
The Controller can be deployed on Kubernetes as a POD. The application then makes use of the ServiceAccount and its respective ServiceAccount Token.
The cluster-reader ClusterRole is required if you want to control resources across the entire cluster.
```bash
kubectl create clusterrolebinding controller-reader \
  --clusterrole=cluster-reader  \
  --serviceaccount=controller-namespace:controller-reader
```

## Installation

Clone the repository and build it.

```bash
git clone https://github.com/mdnix/kubernetes-generic-controller.git
go build .
```

## Controller internals
![alt text](https://raw.githubusercontent.com/mdnix/kubernetes-generic-controller/master/drawing/architecture.jpeg)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
