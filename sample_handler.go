package main

import (
	"github.com/sirupsen/logrus"
)

// Sample handler implements Handler interface,
type Sample struct {
}

// Init initializes handler configuration
// Do nothing for default handler
func (t *Sample) Init(c Config) error {

	return nil
}

// ObjectCreated sends events on object creation
func (t *Sample) ObjectCreated(obj interface{}) {
	logrus.Info("Sample CREATE function invoked")
}

// ObjectUpdated sends events on object updation
func (t *Sample) ObjectUpdated(oldObj, newObj interface{}) {
	logrus.Info("Sample UPDATE function invoked")

}

// ObjectDeleted sends events on object deletion
func (t *Sample) ObjectDeleted(obj interface{}) {
	logrus.Info("Sample DELETE function invoked")

}
