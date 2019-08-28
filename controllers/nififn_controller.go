/*

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

// +kubebuilder:rbac:groups=nififns.nififn.nifi.b23.io,resources=nififns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nififns.nififn.nifi.b23.io,resources=nififns/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;update;patch
package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	nififnv1alpha1 "nifi.b23.io/project/api/v1alpha1"
)

// NiFiFnReconciler reconciles a NiFiFn object
type NiFiFnReconciler struct {
	client.Client
	Log    logr.Logger
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a NiFiFn object and makes changes based on the state read
// and what is in the NiFiFn.Spec
func (r *NiFiFnReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("nififn", req.NamespacedName)

	// Fetch the NiFiFn instance
	instance := &nififnv1alpha1.NiFiFn{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	args := append(
		[]string{
			"RunFromRegistry",
			"Once",
			instance.Spec.RegistryURL,
			instance.Spec.BucketID,
			instance.Spec.FlowID,
			"DestinationDirectory-/tmp/nififn/output2/",
			"", // [<Failure Output Ports>]
		},
		instance.Spec.FlowFiles...)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-job",
			Namespace: instance.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"job": instance.Name + "-job"}},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  instance.Name,
							Image: instance.Spec.Image,
							// Command: []string{""},
							Args: args,
						},
					},
					RestartPolicy: corev1.RestartPolicyOnFailure,
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(instance, job, r.scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if the Job already exists
	found := &batchv1.Job{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Job", "namespace", job.Namespace, "name", job.Name)
		err = r.Create(context.TODO(), job)
		return ctrl.Result{}, err
	} else if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *NiFiFnReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Create a new controller
	c, err := controller.New("nififn-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to NiFiFn
	err = c.Watch(&source.Kind{Type: &nififnv1alpha1.NiFiFn{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch a Job created by NiFiFn
	err = c.Watch(&source.Kind{Type: &batchv1.Job{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &nififnv1alpha1.NiFiFn{},
	})
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&nififnv1alpha1.NiFiFn{}).
		Complete(r)
}
