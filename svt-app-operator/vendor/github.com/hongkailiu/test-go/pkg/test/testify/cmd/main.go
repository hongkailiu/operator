package main

import (
	"fmt"

	"github.com/hongkailiu/test-go/pkg/test/testify/service"
)

func main() {
	d := service.NewDB()

	g := service.NewGreeter(d, "en")
	fmt.Println(g.Greet())             // Message is: hello
	fmt.Println(g.GreetInDefaultMsg()) // Message is: default message

	g = service.NewGreeter(d, "es")
	fmt.Println(g.Greet()) // Message is: holla

	g = service.NewGreeter(d, "random")
	fmt.Println(g.Greet()) // Message is: bzzzz
}
