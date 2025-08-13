package main

import "fmt"

type Human struct {
	name string
	age  int
}

type Action struct {
	Human
}

func (h Human) SayHello() {
	fmt.Printf("Human %s says hello.\n", h.name)
}

func main() {
	human := Human{name: "Roman", age: 19}
	action := Action{Human: human}
	action.SayHello()
}
