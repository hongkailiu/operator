package e2e

import (
	goctx "context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/hongkailiu/operators/svt-app-operator/pkg/apis"
	operator "github.com/hongkailiu/operators/svt-app-operator/pkg/apis/app/v1alpha1"
	myhttp "github.com/hongkailiu/test-go/pkg/http"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"golang.org/x/net/context"
	"gopkg.in/resty.v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
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
		//t.Run("Cluster2", svtCluster)
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

	found := &operator.SVT{}
	err = waitForSVT(f, "example-svt", namespace, found, 3)
	if err != nil {
		return err
	}

	err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "example-svt", Namespace: namespace}, found)
	if err != nil {
		return err
	}
	found.Spec.Size = 1
	err = f.Client.Update(goctx.TODO(), found)
	if err != nil {
		return err
	}
	found.Spec.Size = 0
	err = waitForSVT(f, "example-svt", namespace, found, 1)
	if err != nil {
		return err
	}

	err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "example-svt", Namespace: namespace}, found)
	if err != nil {
		return err
	}
	found.Spec.Size = 2
	err = f.Client.Update(goctx.TODO(), found)
	if err != nil {
		return err
	}
	found.Spec.Size = 0
	err = waitForSVT(f, "example-svt", namespace, found, 2)
	if err != nil {
		return err
	}
	// TODO check the deployed svc on travis-ci only
	// might need a containerized solution if jump node is supported

	// Check if this Service already exists
	foundSVC := &corev1.Service{}
	err = f.Client.Get(context.TODO(), types.NamespacedName{Name: "example-svt", Namespace: namespace}, foundSVC)
	if err != nil {
		return fmt.Errorf("get service with err: %v", err)
	}

	if os.Getenv("CI") == "true" {
		url := fmt.Sprintf("http://%s:8080", foundSVC.Spec.ClusterIP)
		fmt.Println(fmt.Sprintf("accessing url: %s", url))
		resp, err := resty.R().Get(url)
		if err != nil {
			return fmt.Errorf("get service with err: %v", err)
		}
		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("get service with resp.StatusCode(): %d", resp.StatusCode())
		}
		fmt.Println(fmt.Sprintf("resp.Result(): %v", resp.Result()))

		switch resp.Result().(type) {
		case myhttp.Info:
			info := resp.Result().(myhttp.Info)
			fmt.Println(fmt.Sprintf("info.Version: %s", info.Version))
		default:
			return fmt.Errorf("unknown resp.Result(): %v", resp.Result())
		}
	} else {
		fmt.Println(fmt.Sprintf("${CI}!=true, skiping svc checking"))
	}

	return nil
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
	// wait for svt-app-operator to be ready
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "svt-app-operator", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}

	if err = svtScaleTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}

func waitForSVT(f *framework.Framework, svtName string, namespace string, found *operator.SVT, l int) error {
	// wait 10 minutes for len(svt.Status.Nodes) to be satisfied
	err := wait.PollImmediate(10*time.Second, 3*time.Minute, func() (bool, error) {
		err := f.Client.Get(context.TODO(), types.NamespacedName{Name: svtName, Namespace: namespace}, found)
		if err != nil {
			return false, err
		}
		fmt.Println(fmt.Sprintf("%s: len(found.Status.Nodes): %d; l: %d", time.Now().Format(time.RFC3339), len(found.Status.Nodes), l))
		if len(found.Status.Nodes) != l {
			return false, nil
		}
		return true, nil
	})
	if err == wait.ErrWaitTimeout {
		return fmt.Errorf("timed out waiting for len(svt.Status.Nodes) to be satisfied (%d): %d", l, len(found.Status.Nodes))
	}
	return err
}
