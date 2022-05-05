package util

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path"
)

func GetClientSet(isOut bool) *kubernetes.Clientset {
	var clientset *kubernetes.Clientset
	if isOut {
		// creates the out-cluster config
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		config, err := clientcmd.BuildConfigFromFlags("", path.Join(home, ".kube/config"))
		if err != nil {
			panic(err.Error())
		}
		// creates the clientset
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	} else {
		// creates the in-cluster config
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		// creates the clientset
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	}

	return clientset
}

func GetConfigMap(client kubernetes.Interface, namespace, name string) (*corev1.ConfigMap, error) {
	cm, err := client.
		CoreV1().
		ConfigMaps(namespace).
		Get(
			context.Background(),
			name,
			metav1.GetOptions{},
		)
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func CreateConfigMap(client kubernetes.Interface, namespace, name string) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{Data: map[string]string{"foo": "bar"}}
	cm.Namespace = namespace
	cm.Name = name

	cm, err := client.
		CoreV1().
		ConfigMaps(namespace).
		Create(
			context.Background(),
			cm,
			metav1.CreateOptions{},
		)
	if err != nil {
		return nil, err
	}

	return cm, nil
}

func DeleteConfigMap(client kubernetes.Interface, namespace, name string) error {
	err := client.
		CoreV1().
		ConfigMaps(namespace).
		Delete(
			context.Background(),
			name,
			metav1.DeleteOptions{},
		)

	return err
}
