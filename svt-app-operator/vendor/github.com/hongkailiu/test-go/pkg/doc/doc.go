// Package doc shows how to use godoc.
//
// See https://blog.golang.org/godoc-documenting-go-code
//
// by Hongkai Liu
package doc

//NewHello does something cool.
//
//Deprecated: NewHello is defined here to show
//the usage for `Deprecated:` in goDoc.
//
// BUG(who): The NewHello did not say hello nicely.
func NewHello() error {
	return nil
}

//MYNewNiceHello says hello nicely.
func MYNewNiceHello() error {
	return nil
}

//Start starts something
func Start() {
}
