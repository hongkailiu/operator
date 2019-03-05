package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/hongkailiu/test-go/pkg/lib/util"
	"github.com/openshift/api/build/v1"
	buildv1 "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	util.PanicIfError(err)

	buildV1Client, err := buildv1.NewForConfig(config)
	util.PanicIfError(err)

	namespace := "test001project"
	buildConfigs, err := buildV1Client.BuildConfigs(namespace).List(metav1.ListOptions{})
	util.PanicIfError(err)

	fmt.Printf("There are %d builds in the cluster\n", len(buildConfigs.Items))

	//Change namespace and build accordingly

	buildConfig := "cakephp-ex"
	myBuildConfig, err := buildV1Client.BuildConfigs(namespace).Get(buildConfig, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Build %s in namespace %s not found\n", buildConfig, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
			buildConfig, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found build %s in namespace %s\n", buildConfig, namespace)
		fmt.Printf("Found build %s in status %+v\n", buildConfig, myBuildConfig.Status)
	}

	buildRequest := v1.BuildRequest{}
	buildRequest.Kind = "BuildRequest"
	buildRequest.APIVersion = "build.openshift.io/v1"
	objectMeta := metav1.ObjectMeta{}
	objectMeta.Name = "cakephp-ex"
	buildRequest.ObjectMeta = objectMeta
	buildTriggerCause := v1.BuildTriggerCause{}
	buildTriggerCause.Message = "Manually triggered"
	buildRequest.TriggeredBy = []v1.BuildTriggerCause{buildTriggerCause}
	myBuild, err := buildV1Client.BuildConfigs(namespace).Instantiate(buildConfig, &buildRequest)

	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("====Trigger build %s in namespace %s\n", myBuild.Name, myBuild.Namespace)

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pod := myBuild.Name + "-build"
	// https://github.com/kubernetes/helm/blob/master/pkg/kube/wait.go
	fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
	err = wait.Poll(6*time.Second, 10*time.Minute, func() (bool, error) {

		myPod, err := clientset.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		fmt.Printf("Found pod %s in status %+v\n", myPod.Name, myPod.Status)
		if myPod.Status.Phase == "Succeeded" || myPod.Status.Phase == "Failed" {
			return true, nil
		}
		return false, nil
	})

	if err != nil {
		panic(err.Error())
	}

	// not found api on sub-resource (log is a sub-resource of pod)
	// https://stackoverflow.com/questions/32983228/kubernetes-go-client-api-for-log-of-a-particular-pod
	req := clientset.CoreV1().RESTClient().Get().
		Namespace(namespace).
		Name(pod).
		Resource("pods").
		SubResource("log").
		Param("follow", strconv.FormatBool(false)).
		Param("container", "").
		Param("previous", strconv.FormatBool(false)).
		Param("timestamps", strconv.FormatBool(true))

	readCloser, err := req.Stream()
	if err != nil {
		panic(err.Error())
	}

	defer readCloser.Close()

	if b, err := ioutil.ReadAll(readCloser); err == nil {
		fmt.Printf("logs begin ======\n")
		fmt.Printf("%s\n", string(b))
		fmt.Printf("logs end ======\n")
	} else {
		panic(err.Error())
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
