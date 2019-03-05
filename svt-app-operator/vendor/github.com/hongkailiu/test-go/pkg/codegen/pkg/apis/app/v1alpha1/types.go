package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SVTGo describes an SVTGo app.
type SVTGo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SVTGoSpec   `json:"spec"`
	Status SVTGoStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SVTGoList is a list of Database resources
type SVTGoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []SVTGo `json:"items"`
}

// SVTGoSpec is the specification of SVTGo
type SVTGoSpec struct {
	Size int32 `json:"size"`
}

// SVTGoStatus is the status of SVTGo
type SVTGoStatus struct {
	Nodes []string `json:"nodes"`
}
