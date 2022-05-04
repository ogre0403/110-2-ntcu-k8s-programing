package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	config, err := clientcmd.BuildConfigFromFlags("", path.Join(home, ".kube/config"))
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace := "default"

	desired := corev1.ConfigMap{Data: map[string]string{"foo": "bar"}}
	desired.Namespace = namespace
	desired.GenerateName = "crud-typed-simple-"

	// List
	cmList, err := client.
		CoreV1().
		ConfigMaps(namespace).
		List(
			context.Background(),
			metav1.ListOptions{},
		)
	if err != nil {
		panic(err.Error())
	}

	for _, c := range cmList.Items {
		fmt.Printf("Existing ConfigMap name: %s\n", c.Name)
	}
	prompt()

	// Create
	created, err := client.
		CoreV1().
		ConfigMaps(namespace).
		Create(
			context.Background(),
			&desired,
			metav1.CreateOptions{},
		)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Created ConfigMap %s/%s\n", namespace, created.GetName())

	if !reflect.DeepEqual(created.Data, desired.Data) {
		panic("Created ConfigMap has unexpected data")
	}

	prompt()

	// Read
	read, err := client.
		CoreV1().
		ConfigMaps(namespace).
		Get(
			context.Background(),
			created.GetName(),
			metav1.GetOptions{},
		)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Read ConfigMap %s/%s\n", namespace, read.GetName())
	fmt.Println(read.Data["foo"])

	if !reflect.DeepEqual(read.Data, desired.Data) {
		panic("Read ConfigMap has unexpected data")
	}
	prompt()

	// Update
	read.Data["foo"] = "qux"
	updated, err := client.
		CoreV1().
		ConfigMaps(namespace).
		Update(
			context.Background(),
			read,
			metav1.UpdateOptions{},
		)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Updated ConfigMap %s/%s\n", namespace, updated.GetName())
	fmt.Println(read.Data["foo"])

	if !reflect.DeepEqual(updated.Data, read.Data) {
		panic("Updated ConfigMap has unexpected data")
	}
	prompt()

	// Delete
	err = client.
		CoreV1().
		ConfigMaps(namespace).
		Delete(
			context.Background(),
			created.GetName(),
			metav1.DeleteOptions{},
		)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Deleted ConfigMap %s/%s\n", namespace, created.GetName())
}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}
