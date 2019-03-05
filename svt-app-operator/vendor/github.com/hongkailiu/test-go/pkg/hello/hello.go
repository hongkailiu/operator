package main

import (
	"fmt"
	"github.com/hongkailiu/test-go/pkg/stringutil"
)

func main() {
	//fmt.Printf("hello, world\n")
	fmt.Printf(stringutil.Reverse("!oG ,olleH"))
}

//NewHello does something cool.
//
//Deprecated: NewHello is defined here to show
//the usage for `Deprecated` in goDoc
//Hoevre, it does not show in goDoc because it is in main pkg.
//https://stackoverflow.com/questions/21778556/what-steps-are-needed-to-document-package-main-in-godoc
func NewHello() error {
	return nil
}

//NewNiceHello says hello nicely.
func NewNiceHello() error {
	return nil
}
