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

// Package controllers implements kubernetes controllers for nifi-stateless resources
// +kubebuilder:rbac:groups=nififns.nififn.nifi-stateless.b23.io,resources=nififns,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nififns.nififn.nifi-stateless.b23.io,resources=nififns/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;update;patch
package controllers

import (
	"context"
	"encoding/json"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	nifi "nifi-stateless.b23.io/project/api/v1alpha1"
)

const defaultImage = "dbkegley/nifi-stateless:1.10.0-SNAPSHOT"

// NiFiFnReconciler reconciles a NiFiFn object
type NiFiFnReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a NiFiFn object and makes changes based on the state read
// and what is in the NiFiFn.Spec
func (r *NiFiFnReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("nififn", req.NamespacedName)

	constructJobForNiFiFn := func(nififn *nifi.NiFiFn, jobName string) (*batchv1.Job, error) {
		runFromCmdLookup := map[string]string{
			"xml":      "RunFromFlowXml",
			"registry": "RunFromRegistry",
		}
		runFrom := runFromCmdLookup[nififn.Spec.RunFrom]

		image := defaultImage
		if nififn.Spec.Image != "" {
			image = nififn.Spec.Image
		}

		// Make sure nifi_content is defined for every flowfile, otherwise nifi-stateless throws a null ptr
		for _, item := range nififn.Spec.FlowFiles {
			if _, ok := item["nifi_content"]; !ok {
				item["nifi_content"] = ""
			}
		}

		jsonConfig, err := json.Marshal(nififn.Spec)
		if err != nil {
			return nil, err
		}

		args := []string{
			runFrom,
			"Once",
			"--json",
			string(jsonConfig),
		}

		job := &batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      jobName,
				Namespace: nififn.Namespace,
			},
			Spec: batchv1.JobSpec{
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"job": jobName}},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "nifi-stateless",
								Image: image,
								Args:  args,
							},
						},
						RestartPolicy: corev1.RestartPolicyOnFailure,
					},
				},
			},
		}

		if err := ctrl.SetControllerReference(nififn, job, r.Scheme); err != nil {
			return nil, err
		}

		return job, nil
	}

	// Fetch the NiFiFn instance
	var nififn nifi.NiFiFn
	if err := r.Get(ctx, req.NamespacedName, &nififn); err != nil {
		// Only log the error if it's _not_ a NotFound error. Otherwise this prints a stacktrace
		// anytime a NiFiFn resource is deleted
		if ignoreNotFound(err) != nil {
			log.Error(err, "Unable to fetch NiFiFn resource")
		}

		// We'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, ignoreNotFound(err)
	}

	jobName := nififn.Name + "-job"

	// Check if the job already exists
	var job batchv1.Job
	err := r.Get(ctx, types.NamespacedName{Name: jobName, Namespace: nififn.Namespace}, &job)
	if err != nil && errors.IsNotFound(err) {

		// Job wasn't found so create a new one
		log.Info("Constructing Job", "namespace", nififn.Namespace, "name", jobName)
		if job, err := constructJobForNiFiFn(&nififn, jobName); err != nil {
			log.Error(err, "Failed to construct Job resource for NiFiFn")
			return ctrl.Result{}, err
		} else if err := r.Create(ctx, job); err != nil {
			log.Error(err, "Failed to create new Job resource for NiFiFn")
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

var (
	nififnOwnerKey = ".metadata.controller"
	apiGVStr       = nifi.GroupVersion.String()
)

// SetupWithManager registers this initializes this controller and registers it with a controller manager
func (r *NiFiFnReconciler) SetupWithManager(mgr ctrl.Manager) error {

	if err := mgr.GetFieldIndexer().IndexField(&batchv1.Job{}, nififnOwnerKey, func(rawObj runtime.Object) []string {
		// Grab the job object, extract the owner
		job := rawObj.(*batchv1.Job)
		owner := metav1.GetControllerOf(job)
		if owner == nil {
			return nil
		}

		// Make sure it's a NiFiFn
		if owner.APIVersion != apiGVStr || owner.Kind != "NiFiFn" {
			return nil
		}

		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&nifi.NiFiFn{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}

func ignoreNotFound(err error) error {
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}
