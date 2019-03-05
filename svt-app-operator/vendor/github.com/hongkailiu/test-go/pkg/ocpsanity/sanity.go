package ocpsanity

import (
	"os"
	"path/filepath"

	"github.com/hongkailiu/test-go/pkg/lib/logger"
	"github.com/hongkailiu/test-go/pkg/ocutil"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// VERSION of the ocptf cmd
	VERSION = "0.0.8"
)

var (
	oc *ocutil.CLI
	// LogFilePath is the log file path
	LogFilePath = filepath.Join(os.TempDir(), "ocpsanity.log")
	log         = logger.NewLogger(LogFilePath)
)

// SanitySummary represents sanity summary
type SanitySummary struct {
	ProjectTotal       int              `yaml:"projectTotal"`
	ProjectSummaryList []ProjectSummary `yaml:"projectSummaryList"`
}

// ProjectSummary represents project summary
type ProjectSummary struct {
	ProjectName               string `yaml:"projectName"`
	DCTotal                   int    `yaml:"dcTotal"`
	DCReplicaNotSatisfied     int    `yaml:"dcReplicaNotSatisfied"`
	DeployTotal               int    `yaml:"deployTotal"`
	DeployReplicaNotSatisfied int    `yaml:"deployReplicaNotSatisfied"`
	STSTotal                  int    `yaml:"stsTotal"`
	STSReplicaNotSatisfied    int    `yaml:"stsReplicaNotSatisfied"`
	DSTotal                   int    `yaml:"dsTotal"`
	DSDesireNotSatisfied      int    `yaml:"dsDesireNotSatisfied"`

	PodTotal          int `yaml:"podTotal"`
	PodSucceeded      int `yaml:"podSucceeded"`
	PodNotRunning     int `yaml:"podNotRunning"`
	ContainerNotReady int `yaml:"containerNotReady"`
}

// StartSanityCheck starts sanity check on OCP installation
func StartSanityCheck(configPath string) error {
	oc = ocutil.NewCLI(configPath)

	sanitySummary := SanitySummary{ProjectTotal: 0, ProjectSummaryList: []ProjectSummary{}}

	projectList, err := oc.ProjectClient().Projects().List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	sanitySummary.ProjectTotal = len(projectList.Items)

	for _, project := range projectList.Items {
		projectSummary := ProjectSummary{}
		projectSummary.ProjectName = project.Name
		log.WithFields(logrus.Fields{"name": project.Name}).Info("Handle project")

		err := handleDeployConfig(project.Name, &projectSummary)
		if err != nil {
			return err
		}

		err = handleDeployment(project.Name, &projectSummary)
		if err != nil {
			return err
		}

		err = handleSTS(project.Name, &projectSummary)
		if err != nil {
			return err
		}

		err = handleDS(project.Name, &projectSummary)
		if err != nil {
			return err
		}

		err = handlePod(project.Name, &projectSummary)
		if err != nil {
			return err
		}

		sanitySummary.ProjectSummaryList = append(sanitySummary.ProjectSummaryList, projectSummary)
	}

	bytes, err := yaml.Marshal(&sanitySummary)
	if err != nil {
		return err
	}
	log.Infof("--- sanitySummary --- \n%s\n", string(bytes))

	return nil
}

func handleDeployConfig(projectName string, projectSummary *ProjectSummary) error {
	deployConfigList, err := oc.DeployConfigClient().DeploymentConfigs(projectName).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	projectSummary.DCTotal = len(deployConfigList.Items)
	for _, dc := range deployConfigList.Items {
		if dc.Spec.Replicas != dc.Status.Replicas {
			log.WithFields(logrus.Fields{"project": projectName, "dc": dc.Name, "dc.Spec.Replicas": dc.Spec.Replicas,
				"dc.Status.Replicas": dc.Status.Replicas}).
				Warn("Handle deploymentconfig.apps.openshift.io: Replicas not satisfied")
			projectSummary.DCReplicaNotSatisfied++
		} else {
			log.WithFields(logrus.Fields{"project": projectName, "dc": dc.Name, "dc.Spec.Replicas": dc.Spec.Replicas,
				"dc.Status.Replicas": dc.Status.Replicas}).
				Info("Handle deploymentconfig.apps.openshift.io")
		}
	}
	return nil
}

func handleDeployment(projectName string, projectSummary *ProjectSummary) error {
	deploymentList, err := oc.K8SClientSet().AppsV1().Deployments(projectName).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	projectSummary.DeployTotal = len(deploymentList.Items)
	for _, d := range deploymentList.Items {
		if *d.Spec.Replicas != d.Status.Replicas {
			log.WithFields(logrus.Fields{"project": projectName, "dc": d.Name, "d.Spec.Replicas": *d.Spec.Replicas,
				"d.Status.Replicas": d.Status.Replicas}).
				Warn("Handle deployment: Replicas not satisfied")
			projectSummary.DeployReplicaNotSatisfied++
		} else {
			log.WithFields(logrus.Fields{"project": projectName, "dc": d.Name, "d.Spec.Replicas": *d.Spec.Replicas,
				"d.Status.Replicas": d.Status.Replicas}).
				Info("Handle deployment")
		}
	}
	return nil
}

func handleSTS(projectName string, projectSummary *ProjectSummary) error {
	statefulSetList, err := oc.K8SClientSet().AppsV1().StatefulSets(projectName).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	projectSummary.STSTotal = len(statefulSetList.Items)
	for _, ss := range statefulSetList.Items {
		if *ss.Spec.Replicas != ss.Status.Replicas {
			log.WithFields(logrus.Fields{"project": projectName, "ss": ss.Name, "ss.Spec.Replicas": *ss.Spec.Replicas,
				"ss.Status.Replicas": ss.Status.Replicas}).
				Warn("Handle sts: Replicas not satisfied")
			projectSummary.STSReplicaNotSatisfied++
		} else {
			log.WithFields(logrus.Fields{"project": projectName, "ss": ss.Name, "ss.Spec.Replicas": *ss.Spec.Replicas,
				"ss.Status.Replicas": ss.Status.Replicas}).
				Info("Handle sts")
		}
	}
	return nil
}

func handleDS(projectName string, projectSummary *ProjectSummary) error {
	daemonSetList, err := oc.K8SClientSet().AppsV1().DaemonSets(projectName).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	projectSummary.DSTotal = len(daemonSetList.Items)
	for _, ds := range daemonSetList.Items {
		if ds.Status.DesiredNumberScheduled != ds.Status.CurrentNumberScheduled {
			log.WithFields(logrus.Fields{"project": projectName, "ds": ds.Name,
				"ds.Status.DesiredNumberScheduled": ds.Status.DesiredNumberScheduled,
				"ds.Status.CurrentNumberScheduled": ds.Status.CurrentNumberScheduled}).
				Warn("Handle daemon set: Not scheduled")
			projectSummary.DSDesireNotSatisfied++
		} else {
			log.WithFields(logrus.Fields{"project": projectName, "ds": ds.Name,
				"ds.Status.DesiredNumberScheduled": ds.Status.DesiredNumberScheduled,
				"ds.Status.CurrentNumberScheduled": ds.Status.CurrentNumberScheduled}).
				Info("Handle daemon set")
		}
	}
	return nil
}

func handlePod(projectName string, projectSummary *ProjectSummary) error {
	PodList, err := oc.K8SClientSet().CoreV1().Pods(projectName).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	projectSummary.PodTotal = len(PodList.Items)
	for _, pod := range PodList.Items {
		if pod.Status.Phase == corev1.PodSucceeded {
			log.WithFields(logrus.Fields{"project": projectName, "pod": pod.Name, "pod.Status.Phase": pod.Status.Phase}).
				Info("Handle pod: Succeeded")
			projectSummary.PodSucceeded++
			break
		}
		if pod.Status.Phase != corev1.PodRunning {
			log.WithFields(logrus.Fields{"project": projectName, "pod": pod.Name, "pod.Status.Phase": pod.Status.Phase,
				"pod.Status.Reason": pod.Status.Reason}).
				Warn("Handle pod: Not Running")
			projectSummary.PodNotRunning++
		} else {
			for _, cs := range pod.Status.ContainerStatuses {
				if !cs.Ready {
					log.WithFields(logrus.Fields{"project": projectName, "pod": pod.Name,
						"pod.Status.Phase": pod.Status.Phase, "cs.Name": cs.Name,
						"cs.ContainerID": cs.ContainerID, "cs.Ready": cs.Ready}).
						Warn("Handle pod: Not Ready")
					projectSummary.ContainerNotReady++
				} else {
					log.WithFields(logrus.Fields{"project": projectName, "pod": pod.Name,
						"pod.Status.Phase": pod.Status.Phase, "cs.Name": cs.Name,
						"cs.ContainerID": cs.ContainerID, "cs.Ready": cs.Ready}).
						Info("Handle pod")
				}
			}

		}
	}
	return nil
}
