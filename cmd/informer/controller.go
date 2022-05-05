package main

import (
	"fmt"
	"gitlab.com/ogre0403/110-2-ntcu-k8s-programing/pkg/util"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/informers"
	batchinformer "k8s.io/client-go/informers/batch/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"time"
)

type ConfigMapController struct {
	informerFactory informers.SharedInformerFactory
	informer        batchinformer.JobInformer
	clientSet       *kubernetes.Clientset
}

// Run starts shared informers and waits for the shared informer cache to
// synchronize.
func (c *ConfigMapController) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so
	// far.
	c.informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.informer.Informer().HasSynced) {
		return fmt.Errorf("Failed to sync")
	}
	return nil
}

func (c *ConfigMapController) onAdd(obj interface{}) {
	job := obj.(*batchv1.Job)
	fmt.Printf("Informer event: Job ADDED %s/%s\n", job.GetNamespace(), job.GetName())
}

func (c *ConfigMapController) onUpdate(old, new interface{}) {
	job := old.(*batchv1.Job)
	fmt.Printf("Informer event: Job UPDATED %s/%s\n", job.GetNamespace(), job.GetName())

	_, err := util.GetConfigMap(c.clientSet, namespace, "test-configmap")
	if err != nil && errors.IsNotFound(err) {
		cm, _ := util.CreateConfigMap(c.clientSet, namespace, "test-configmap")
		fmt.Printf("----Create ConfigMap when Job UPDATED Event %s/%s\n", cm.GetNamespace(), cm.GetName())
	}

}

func (c *ConfigMapController) onDelete(obj interface{}) {
	job := obj.(*batchv1.Job)
	fmt.Printf("Informer event: Job DELETED %s/%s\n", job.GetNamespace(), job.GetName())

	if err := util.DeleteConfigMap(c.clientSet, namespace, "test-configmap"); err == nil {
		fmt.Printf("----Delete ConfigMap when Job DELETE Event %s/%s\n", namespace, "test-configmap")
	}
}

// NewConfigMapController creates a ConfigMapController
func NewConfigMapController(client *kubernetes.Clientset) *ConfigMapController {
	factory := informers.NewSharedInformerFactoryWithOptions(client, 5*time.Second, informers.WithNamespace(namespace))
	informer := factory.Batch().V1().Jobs()

	c := &ConfigMapController{
		informerFactory: factory,
		informer:        informer,
		clientSet:       client,
	}
	informer.Informer().AddEventHandler(
		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: c.onAdd,
			// Called on resource update and every resyncPeriod on existing resources.
			UpdateFunc: c.onUpdate,
			// Called on resource deletion.
			DeleteFunc: c.onDelete,
		},
	)
	return c
}
