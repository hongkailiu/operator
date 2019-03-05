package service

import "fmt"

// LoginService represents login services
type LoginService interface {
	Login(username, password string) string
}

type myLoginService struct {
	name string
}

func (s *myLoginService) Login(username, password string) string {
	fmt.Println(fmt.Sprintf("This is %s.", s.name))
	if username == "username" && password == "password" {
		return "OK"
	}
	return "NOT-OK"
}

// NewService returns a LoginService object
func NewService() LoginService {
	return &myLoginService{name: "cool login service"}
}
