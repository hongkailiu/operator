package e2e

import (
	goctx "context"
	"fmt"
	"testing"
	"time"

	apis "github.com/hongkailiu/operators/svt-app-operator/pkg/apis"
	operator "github.com/hongkailiu/operators/svt-app-operator/pkg/apis/app/v1alpha1"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 60
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestSVT(t *testing.T) {
	svtList := &operator.SVTList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "svt",
			APIVersion: "app.test.com/v1alpha1",
		},
	}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, svtList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("svt-group", func(t *testing.T) {
		t.Run("Cluster", svtCluster)
		t.Run("Cluster2", svtCluster)
	})
}

func svtScaleTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}
	// create svt custom resource
	exampleSVT := &operator.SVT{
		TypeMeta: metav1.TypeMeta{
			Kind:       "svt",
			APIVersion: "app.test.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-svt",
			Namespace: namespace,
		},
		Spec: operator.SVTSpec{
			Size: 3,
		},
	}
	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), exampleSVT, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}
	// wait for example-svt to reach 3 replicas
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-svt", 3, retryInterval, timeout)
	if err != nil {
		return err
	}

	err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "example-svt", Namespace: namespace}, exampleSVT)
	if err != nil {
		return err
	}
	exampleSVT.Spec.Size = 2
	err = f.Client.Update(goctx.TODO(), exampleSVT)
	if err != nil {
		return err
	}

	// wait for example-svt to reach 4 replicas
	return e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-svt", 2, retryInterval, timeout)
}

func svtCluster(t *testing.T) {
	t.Parallel()
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()
	err := ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}
	t.Log("Initialized cluster resources")
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}
	// get global framework variables
	f := framework.Global
	// wait for svt-operator to be ready
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "svt-operator", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}

	if err = svtScaleTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}