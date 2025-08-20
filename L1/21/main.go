package main

import "fmt"

type Target interface {
	operation()
}

type Adaptable struct{}

func (adapted *Adaptable) AdaptableOperation() {
	fmt.Println("hello")
}

type Adapter struct {
	adapted *Adaptable
}

func (a *Adapter) operation() {
	a.adapted.AdaptableOperation()
}

func NewAdapter(adapted *Adaptable) Target {
	return &Adapter{adapted}
}

func main() {
	adapter := NewAdapter(&Adaptable{})
	adapter.operation()
}
