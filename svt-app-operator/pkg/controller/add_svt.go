package controller

import (
	"github.com/hongkailiu/operators/svt-app-operator/pkg/controller/svt"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, svt.Add)
}
