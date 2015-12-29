# shutdown
offer a convenient way to intercept Ctrl+C or kill signal and do everything you want before program exit in golang

## Installation

To install, simply run in a terminal:

    go get github.com/ZGeomantic/shutdown

  
## Usage

The following example (taken from /example/example.go) shows how to add the functions or struct's entity you want to handle after Ctrl+C is pressed and before the program is shutdown. 


```go
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

```
