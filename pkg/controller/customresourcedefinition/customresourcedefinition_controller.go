package customresourcedefinition

import (
	"context"
	"fmt"

	"github.com/openshift/rbac-permissions-operator/config"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_customresourcedefinition")

// Add creates a new CustomResourceDefinition Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileCustomResourceDefinition{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("customresourcedefinition-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource CustomResourceDefinition
	err = c.Watch(&source.Kind{Type: &apiextensions.CustomResourceDefinition{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileCustomResourceDefinition implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCustomResourceDefinition{}

// ReconcileCustomResourceDefinition reconciles a CustomResourceDefinition object
type ReconcileCustomResourceDefinition struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a CustomResourceDefinition object and makes changes based on the state read
// and what is in the CustomResourceDefinition.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileCustomResourceDefinition) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name) //Operator namespace & object
	reqLogger.Info("Reconciling CustomResourceDefinition")

	// Fetch the CustomResourceDefinition list
	crd := &apiextensions.CustomResourceDefinition{}
	err := r.client.Get(context.TODO(), request.NamespacedName, crd)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		//return reconcile.Result{}, err
		//}

		crdName := crd.Spec.Names.Plural
		groupName := crd.Spec.Group
		clusterRole := &rbac.ClusterRole{}
		clusterRoleName := ""

		//Determine if the crd is namespace scoped or cluster scoped
		if crd.Spec.Scope == "namespaced" {
			clusterRoleName = config.CRDClusterRoleNamespaced
		} else {
			clusterRoleName = config.CRDClusterRoleGlobal
		}

		err = r.client.Get(context.TODO(), types.NamespacedName{Name: clusterRoleName}, clusterRole)
		if err != nil {
			failedToGetCRDMsg := fmt.Sprintf("Failed to get Cluster Role %s", clusterRoleName)
			reqLogger.Error(err, failedToGetCRDMsg)
			return reconcile.Result{}, err
		}

		// if found = true, break.
		// if found = false, add permission. via appending a new object, use r.client.update to update the role.
		found := isPermissionInClusterrole(crdName, groupName, clusterRole)
		if found == true {
			//Permission is already present
			return reconcile.Result{}, nil
		} else {
			// Mapping to store what will be added to the role/clusterrole
			newRule := ruleTemplate

			newRule.APIGroups = []string{
				groupName,
			}

			newRule.Resources = []string{
				crdName,
			}

			clusterRole.Rules = append(clusterRole.Rules, newRule)
			//Logic to add the newRule to the clusterrole

			return reconcile.Result{}, err
		}

		// Need a loop here, for blacklist matching via Regex to something like crd.Spec.Version for matching against apigroup
	}
	return reconcile.Result{}, nil
}

func isPermissionInClusterrole(crdName string, groupName string, clusterRole *rbac.ClusterRole) bool {
	rulesList := clusterRole.Rules
	doesAPIGroupMatch := false
	for _, permission := range rulesList {
		for _, ag := range permission.APIGroups {
			if ag == groupName {
				doesAPIGroupMatch = true
			}
		}
		for _, resource := range permission.Resources {
			if resource == crdName && doesAPIGroupMatch == true {
				return true
			}
		}
		doesAPIGroupMatch = false
	}
	return false
}
