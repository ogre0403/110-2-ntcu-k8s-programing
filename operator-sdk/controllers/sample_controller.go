/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	samplev1alpha1 "github.com/ogre0403/110-2-ntcu-k8s-programing/api/v1alpha1"
)

// SampleReconciler reconciles a Sample object
type SampleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=sample.ntcu.edu.tw,resources=samples,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=sample.ntcu.edu.tw,resources=samples/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=sample.ntcu.edu.tw,resources=samples/finalizers,verbs=update

//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Sample object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *SampleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// TODO(user): your logic here
	// Fetch the Memcached instance
	sample := &samplev1alpha1.Sample{}
	err := r.Get(ctx, req.NamespacedName, sample)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Sample resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get sample")
		return ctrl.Result{}, err
	}

	foundJob := &batchv1.Job{}
	err = r.Get(ctx, types.NamespacedName{Name: sample.Name, Namespace: sample.Namespace}, foundJob)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.newJob(sample)
		log.Info("Creating a new Job", "Job.Namespace", dep.Namespace, "Job.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Job", "Job.Namespace", dep.Namespace, "Job.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Job")
		return ctrl.Result{}, err
	}

	foundCM := &corev1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Name: sample.Name, Namespace: sample.Namespace}, foundCM)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.newConfigMap(sample)
		log.Info("Creating a new ConfigMap", "ConfigMap.Namespace", dep.Namespace, "ConfigMap.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", dep.Namespace, "ConfigMap.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get ConfigMap")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *SampleReconciler) newConfigMap(s *samplev1alpha1.Sample) *corev1.ConfigMap {

	cm := &corev1.ConfigMap{
		Data: map[string]string{s.Spec.Key: s.Spec.Value},
	}
	cm.Name = s.Name
	cm.Namespace = s.Namespace

	ctrl.SetControllerReference(s, cm, r.Scheme)
	return cm
}

func (r *SampleReconciler) newJob(s *samplev1alpha1.Sample) *batchv1.Job {
	job := &batchv1.Job{
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: s.Name,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   s.Spec.Image,
						Name:    s.Name,
						Command: []string{"cat", fmt.Sprintf("/tmp/cm/%s", s.Spec.Key)},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      s.Name,
								MountPath: "/tmp/cm",
							},
						},
					}},
					Volumes: []corev1.Volume{
						{
							Name: s.Name,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: s.Name,
									},
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}
	job.Name = s.Name
	job.Namespace = s.Namespace

	ctrl.SetControllerReference(s, job, r.Scheme)
	return job
}

// SetupWithManager sets up the controller with the Manager.
func (r *SampleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&samplev1alpha1.Sample{}).
		Owns(&batchv1.Job{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
