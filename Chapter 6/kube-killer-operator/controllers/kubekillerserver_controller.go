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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/p-program/kube-killer-operator/api/v1alpha1"
	bullshitprogramcomv1alpha1 "github.com/p-program/kube-killer-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var log = logf.Log.WithName("controller_server")

// KubeKillerServerReconciler reconciles a KubeKillerServer object
type KubeKillerServerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=bullshitprogram.com,resources=kubekillerservers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=bullshitprogram.com,resources=kubekillerservers/status,verbs=get;update;patch

func (r *KubeKillerServerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	// your logic here
	reqLogger := log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
	// Fetch the Server instance
	instance := &v1alpha1.KubeKillerServer{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		// Request object not found, could have been deleted after reconcile request.
		// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
		// Return and don't requeue
		reqLogger.Error(err, "CR not found")
		found := &appsv1.Deployment{}
		err = r.Get(context.TODO(), types.NamespacedName{Name: "kube-killer-operator" + req.Name, Namespace: req.Namespace}, found)
		if err == nil && found != nil {
			// delete logic
			reqLogger.Info("deleting deploy ", "name: ", instance.Name)
			err := r.Delete(context.TODO(), found)
			if err != nil {
				return reconcile.Result{}, err
			}
			oldServer := &bullshitprogramcomv1alpha1.KubeKillerServer{ObjectMeta: metav1.ObjectMeta{Name: req.Name, Namespace: req.Namespace}}
			service := newServiceForCR(oldServer)
			reqLogger.Info("deleting service ", "name: ", service.Name)
			err = r.Delete(context.TODO(), service)
			if err != nil {
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}
	found := &appsv1.Deployment{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: "kube-killer-operator" + req.Name, Namespace: req.Namespace}, found)
	if err != nil {
		reqLogger.Info("creating resource")
		deploy := newDeployForCR(instance)
		err = r.Create(context.TODO(), deploy)
		if err != nil {
			log.Error(err, "create Deploy fail")
			return reconcile.Result{}, err
		}
		service := newServiceForCR(instance)
		err = r.Create(context.TODO(), service)
		if err != nil {
			log.Error(err, "create service fail")
			return reconcile.Result{}, err
		}
		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	}

	needUpate := false
	// get deploy and then compare with CR
	if instance.Spec.Replica != *found.Spec.Replicas {
		needUpate = true
	}
	// update deploy
	if needUpate && found != nil {
		reqLogger.Info("updateing resource")
		found.Spec.Replicas = &instance.Spec.Replica
		r.Update(context.TODO(), found)
		return reconcile.Result{}, nil
	}
	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: resource already exists")
	return reconcile.Result{}, nil
}

// SetupWithManager run once
func (r *KubeKillerServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bullshitprogramcomv1alpha1.KubeKillerServer{}).
		Complete(r)
}

func newServiceForCR(cr *v1alpha1.KubeKillerServer) *corev1.Service {
	labels := map[string]string{
		"app":        cr.Name,
		"controller": "kube-killer-operator" + cr.Name,
	}
	port := corev1.ServicePort{
		Name:       "nginx",
		Protocol:   corev1.ProtocolTCP,
		Port:       80,
		TargetPort: intstr.FromInt(80)}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kube-killer-operator" + cr.Name,
			Namespace: cr.ObjectMeta.Namespace,
			// OwnerReferences: []metav1.OwnerReference{
			// 	*metav1.NewControllerRef(cr, v1alpha1.SchemeBuilder.GroupVersion.WithKind("Server")),
			// }
		},
		Spec: corev1.ServiceSpec{Ports: []corev1.ServicePort{port}, Selector: labels, Type: corev1.ServiceTypeClusterIP},
	}

}

func newDeployForCR(cr *v1alpha1.KubeKillerServer) *appsv1.Deployment {
	labels := map[string]string{
		"app":        cr.Name,
		"controller": "kube-killer-operator",
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kube-killer-operator" + cr.Name,
			Namespace: cr.ObjectMeta.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &cr.Spec.Replica,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "nginx",
							Image:           cr.Spec.Image,
							ImagePullPolicy: "IfNotPresent",
						},
					},
				},
			},
		},
	}
}
