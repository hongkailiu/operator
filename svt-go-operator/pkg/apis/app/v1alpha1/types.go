package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SVTGoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []SVTGo `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SVTGo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              SVTGoSpec   `json:"spec"`
	Status            SVTGoStatus `json:"status,omitempty"`
}

type SVTGoSpec struct {
	// Size is the size of the memcached deployment
	Size int32 `json:"size"`
}
type SVTGoStatus struct {
	// Nodes are the names of the memcached pods
	Nodes []string `json:"nodes"`
}
