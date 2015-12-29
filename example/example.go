package main

import (
	"fmt"
	"github.com/ZGeomantic/shutdown"
)

type Foo struct {
	Source string
}

func (this *Foo) OnShutdown() {
	fmt.Println("close anything you need to close here as a struct")
}

func BarFunc() {
	fmt.Println("close anything you need to close here as a function")
}

func main() {
	var foo Foo
	shutdown.Register(&foo)
	shutdown.Register(shutdown.Func(BarFunc))
	shutdown.WaitingShutDown()
}
