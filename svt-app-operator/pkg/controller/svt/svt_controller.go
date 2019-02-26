package svt

import (
	"context"
	"fmt"
	"reflect"
	"time"

	appv1alpha1 "github.com/hongkailiu/operators/svt-app-operator/pkg/apis/app/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_svt")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new SVT Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSVT{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("svt-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource SVT
	err = c.Watch(&source.Kind{Type: &appv1alpha1.SVT{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner SVT
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appv1alpha1.SVT{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &appv1alpha1.SVT{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileSVT{}

// ReconcileSVT reconciles a SVT object
type ReconcileSVT struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a SVT object and makes changes based on the state read
// and what is in the SVT.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileSVT) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SVT==================================")

	// Fetch the SVT instance
	instance := &appv1alpha1.SVT{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	deployment := deploymentForSVT(instance)

	// Set SVT instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, deployment, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	} else {
		// Deployment already exists - don't requeue
		reqLogger.Info("Skip reconcile: Deployment already exists", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	}

	// Create the deployment if it doesn't exist
	svc := serviceForSVT(instance)

	// Set SVT instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, svc, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	foundSVC := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, foundSVC)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.client.Create(context.TODO(), svc)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	} else {
		// Service already exists - don't requeue
		reqLogger.Info("Skip reconcile: Service already exists", "Service.Namespace", foundSVC.Namespace, "Service.Name", foundSVC.Name)
	}

	// Ensure the deployment size is the same as the spec
	size := instance.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.client.Update(context.TODO(), found)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to update deployment: %v", err)
		}
		return reconcile.Result{Requeue: true}, nil
	}

	// wait 10 minutes for deployment's replicas to be satisfied
	err = wait.PollImmediate(10*time.Second, 10*time.Minute, func() (bool, error) {
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
		if err != nil {
			return false, fmt.Errorf("failed to get deployment: %v", err)
		}
		reqLogger.Info("got values ...",
			"*found.Spec.Replicas", *found.Spec.Replicas, "found.Status.AvailableReplicas", found.Status.AvailableReplicas)
		if *found.Spec.Replicas != found.Status.AvailableReplicas {
			reqLogger.Info("waiting for deployment's replicas to be satisfied ...",
				"*found.Spec.Replicas", *found.Spec.Replicas, "found.Status.AvailableReplicas", found.Status.AvailableReplicas)
			reqLogger.Info("waiting for deployment's replicas to be satisfied ...")
			return false, nil
		}
		return true, nil
	})
	if err == wait.ErrWaitTimeout {
		return reconcile.Result{}, fmt.Errorf("timed out waiting for deployment's replicas to be satisfied: %s", found.Name)
	}
	if err != nil {
		return reconcile.Result{}, err
	}

	// Update the svtgo status with the pod names
	podList := podList()
	labelSelector := labels.SelectorFromSet(labelsForSVT(instance.Name))
	listOps := &client.ListOptions{
		Namespace:     request.NamespacedName.Namespace,
		LabelSelector: labelSelector,}
	err = r.client.List(context.TODO(), listOps, podList)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to list pods: %v", err)
	}
	podNames := getPodNames(podList.Items)
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		//https://github.com/operator-framework/operator-sdk/blob/master/doc/user/client.md#updating-status-subresource
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil && errors.IsResourceExpired(err) {
			return reconcile.Result{Requeue: true}, fmt.Errorf("failed to update svt status (ResourceExpired): %v", err)
		} else {
			return reconcile.Result{}, fmt.Errorf("failed to update svt status: %v", err)
		}
	}
	return reconcile.Result{}, nil
}

// deploymentForSVT returns a svt Deployment object
func deploymentForSVT(m *appv1alpha1.SVT) *appsv1.Deployment {
	ls := labelsForSVT(m.Name)
	replicas := m.Spec.Size

	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   "quay.io/hongkailiu/test-go:http-0.0.13",
						Name:    "svt",
						Command: []string{"/http"},
						Args:    []string{"start", "-v"},
						Env:     []corev1.EnvVar{{Name: "GIN_MODE", Value: "release"}},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
							Name:          "http",
						}},
					}},
				},
			},
		},
	}
	return dep
}

// serviceForSVT returns a svt Deployment object
func serviceForSVT(m *appv1alpha1.SVT) *corev1.Service {
	ls := labelsForSVT(m.Name)

	svt := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{Port: 8080, Protocol: corev1.ProtocolTCP, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8080}},
			},
			Type:     corev1.ServiceTypeClusterIP,
			Selector: ls,
		},
	}
	return svt
}

// labelsForSVTGo returns the labels for selecting the resources
// belonging to the given svtgo CR name.
func labelsForSVT(name string) map[string]string {
	return map[string]string{"app": "svt", "svt_cr": name}
}

// podList returns a v1.PodList object
func podList() *corev1.PodList {
	return &corev1.PodList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
	}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}
