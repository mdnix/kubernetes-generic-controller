package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const maxRetries = 10

var serverStartTime time.Time

// Config struct generated from config.yaml
type Config struct {
	Resources struct {
		Pod                   bool
		Deployment            bool
		Replicationcontroller bool
		Replicaset            bool
		Daemonset             bool
		Services              bool
		Secret                bool
		Configmap             bool
		Rolebinding           bool
	} `yaml:"resources"`
	Handler struct {
		Name string
	} `yaml:"handler"`
}

// Event indicate the informerEvent
type Event struct {
	key          string
	eventType    string
	namespace    string
	resourceType string
}

// Controller object
type Controller struct {
	logger       *logrus.Entry
	clientset    kubernetes.Interface
	queue        workqueue.RateLimitingInterface
	informer     cache.SharedIndexInformer
	eventHandler Handler
}
