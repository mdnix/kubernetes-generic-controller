package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

)

// Start prepares watchers and run their controllers, then waits for process termination signals
func startController(conf Config) {
	clientset, err := getClient()
	if err != nil {
		panic(err.Error())
	}

	if conf.Resources.Pod {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return clientset.CoreV1().Pods("").List(options)
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return clientset.CoreV1().Pods("").Watch(options)
				},
			},
			&apiv1.Pod{},
			0, //Skip resync
			cache.Indexers{},
		)

		eventHandler := ParseEventHandler(conf)
		c := newResourceController(clientset, eventHandler, informer, "pod")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh, "pod")
	}

	if conf.Resources.Deployment {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return clientset.AppsV1().Deployments("").List(options)
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return clientset.AppsV1().Deployments("").Watch(options)
				},
			},
			&appsv1.Deployment{},
			0, //Skip resync
			cache.Indexers{},
		)

		eventHandler := ParseEventHandler(conf)
		c := newResourceController(clientset, eventHandler, informer, "deployment")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh, "deployment")
	}

	if conf.Resources.Replicationcontroller {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return clientset.CoreV1().ReplicationControllers("").List(options)
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return clientset.CoreV1().ReplicationControllers("").Watch(options)
				},
			},
			&apiv1.ReplicationController{},
			0, //Skip resync
			cache.Indexers{},
		)

		eventHandler := ParseEventHandler(conf)
		c := newResourceController(clientset, eventHandler, informer, "replicationcontroller")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh, "replicationcontroller")
	}

	if conf.Resources.Replicaset {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return clientset.AppsV1().ReplicaSets("").List(options)
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return clientset.AppsV1().ReplicaSets("").Watch(options)
				},
			},
			&appsv1.ReplicaSet{},
			0, //Skip resync
			cache.Indexers{},
		)

		eventHandler := ParseEventHandler(conf)
		c := newResourceController(clientset, eventHandler, informer, "replicaset")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh, "replicaset")
	}

	if conf.Resources.Daemonset {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return clientset.AppsV1().DaemonSets("").List(options)
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return clientset.AppsV1().DaemonSets("").Watch(options)
				},
			},
			&appsv1.DaemonSet{},
			0, //Skip resync
			cache.Indexers{},
		)

		eventHandler := ParseEventHandler(conf)
		c := newResourceController(clientset, eventHandler, informer, "daemonset")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh, "daemonset")
	}

	if conf.Resources.Services {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return clientset.CoreV1().Services("").List(options)
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return clientset.CoreV1().Services("").Watch(options)
				},
			},
			&apiv1.Service{},
			0, //Skip resync
			cache.Indexers{},
		)

		eventHandler := ParseEventHandler(conf)
		c := newResourceController(clientset, eventHandler, informer, "service")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh, "service")
	}

	if conf.Resources.Secret {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return clientset.CoreV1().Secrets("").List(options)
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return clientset.CoreV1().Secrets("").Watch(options)
				},
			},
			&apiv1.Secret{},
			0, //Skip resync
			cache.Indexers{},
		)

		eventHandler := ParseEventHandler(conf)
		c := newResourceController(clientset, eventHandler, informer, "secret")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh, "secret")
	}

	if conf.Resources.Configmap {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return clientset.CoreV1().ConfigMaps("").List(options)
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return clientset.CoreV1().ConfigMaps("").Watch(options)
				},
			},
			&apiv1.ConfigMap{},
			0, //Skip resync
			cache.Indexers{},
		)
		eventHandler := ParseEventHandler(conf)
		c := newResourceController(clientset, eventHandler, informer, "configmap")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh, "configmap")
	}

	if conf.Resources.Rolebinding {
		informer := cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return clientset.RbacV1().RoleBindings("").List(options)
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return clientset.RbacV1().RoleBindings("").Watch(options)
				},
			},
			&rbacv1.RoleBinding{},
			0, //Skip resync
			cache.Indexers{},
		)

		eventHandler := ParseEventHandler(conf)
		c := newResourceController(clientset, eventHandler, informer, "rolebinding")
		stopCh := make(chan struct{})
		defer close(stopCh)

		go c.Run(stopCh, "rolebinding")
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm
}

func newResourceController(client kubernetes.Interface, eventHandler Handler, informer cache.SharedIndexInformer, resourceType string) *Controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	var newEvent Event
	var err error
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			newEvent.key, err = cache.MetaNamespaceKeyFunc(obj)
			newEvent.eventType = "create"
			newEvent.resourceType = resourceType
			logrus.WithField("resource", resourceType).Infof("Processing add to %v: %s", resourceType, newEvent.key)
			if err == nil {
				queue.Add(newEvent)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			newEvent.key, err = cache.MetaNamespaceKeyFunc(old)
			newEvent.eventType = "update"
			newEvent.resourceType = resourceType
			logrus.WithField("resource", resourceType).Infof("Processing update to %v: %s", resourceType, newEvent.key)
			if err == nil {
				queue.Add(newEvent)
			}
		},
		DeleteFunc: func(obj interface{}) {
			newEvent.key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			newEvent.eventType = "delete"
			newEvent.resourceType = resourceType
			newEvent.namespace = GetObjectMetaData(obj).Namespace
			logrus.WithField("resource", resourceType).Infof("Processing delete to %v: %s", resourceType, newEvent.key)
			if err == nil {
				queue.Add(newEvent)
			}
		},
	})

	return &Controller{
		logger:       logrus.WithField("resource", resourceType),
		clientset:    client,
		informer:     informer,
		queue:        queue,
		eventHandler: eventHandler,
	}
}

// Run starts the controller controller
func (c *Controller) Run(stopCh <-chan struct{}, resourceType string) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	c.logger.Infof("Starting %s controller", resourceType)
	serverStartTime = time.Now().Local()

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	c.logger.Infof("%s controller synced and ready", resourceType)

	wait.Until(c.runWorker, time.Second, stopCh)
}

func (c *Controller) runWorker() {
	// processNextWorkItem will automatically wait until there's work available
	for c.processNextItem() {
		// continue looping
	}
}

// HasSynced is required for the cache.Controller interface.
func (c *Controller) HasSynced() bool {
	return c.informer.HasSynced()
}

// LastSyncResourceVersion is required for the cache.Controller interface.
func (c *Controller) LastSyncResourceVersion() string {
	return c.informer.LastSyncResourceVersion()
}

// processNextWorkItem deals with one key off the queue.  It returns false
// when it's time to quit.
func (c *Controller) processNextItem() bool {
	newEvent, quit := c.queue.Get()

	if quit {
		return false
	}
	defer c.queue.Done(newEvent)
	err := c.processItem(newEvent.(Event))
	if err == nil {
		// No error, reset the ratelimit counters
		c.queue.Forget(newEvent)
	} else if c.queue.NumRequeues(newEvent) < maxRetries {
		c.logger.Errorf("Error processing %s (will retry): %v", newEvent.(Event).key, err)
		c.queue.AddRateLimited(newEvent)
	} else {
		// err != nil and too many retries
		c.logger.Errorf("Error processing %s (giving up): %v", newEvent.(Event).key, err)
		c.queue.Forget(newEvent)
		utilruntime.HandleError(err)
	}
	return true
}

func (c *Controller) processItem(newEvent Event) error {
	obj, _, err := c.informer.GetIndexer().GetByKey(newEvent.key)
	if err != nil {
		return fmt.Errorf("Error fetching object with key %s from store: %v", newEvent.key, err)
	}
	// get object's metedata
	objectMeta := GetObjectMetaData(obj)

	// process events based on its type
	switch newEvent.eventType {
	case "create":
		// compare CreationTimestamp and serverStartTime and alert only on latest events
		// Could be Replaced by using Delta or DeltaFIFO
		if objectMeta.CreationTimestamp.Sub(serverStartTime).Seconds() > 0 {
			c.eventHandler.ObjectCreated(obj)
		}
	case "update":
		//c.eventHandler.ObjectUpdated(obj)

	case "delete":
		c.eventHandler.ObjectDeleted(obj)
	}
	return nil
}
