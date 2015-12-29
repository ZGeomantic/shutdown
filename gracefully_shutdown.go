package shutdown

import (
	"container/list"
	"log"
	"os"
	"os/signal"
	"sync"
)

// 被关闭者要实现的接口
// the only interface that need to be implemented
type GracefulClose interface {
	OnShutdown()
}

// 如果无法让被关闭的对象实现接口，可以通过类型转换一个func()形式的闭包函数，来简单地实现GracefulClose接口
// if the shutdown logic is just a function but not a struct, that's ok. just cast it type to Func.
// for example:
//  func YourCloseProcedure(){
//	 your code
// }
// Call it like this:
// Register(shutdown.Func(YourCloseProcedure))
type Func func()

func (shutdown Func) OnShutdown() {
	shutdown()
}

type ShutdownQueue struct {
	registerList *list.List
	registerLock sync.Mutex
}

var defaultCloserList ShutdownQueue
var shutdownSemaphore chan string

func Register(closer GracefulClose) {

	defaultCloserList.registerLock.Lock()
	defer defaultCloserList.registerLock.Unlock()
	defaultCloserList.registerList.PushBack(closer)
}

// 对外暴露的阻塞式的等待方法
func WaitingShutDown() {
	<-shutdownSemaphore
	// os.Exit(0)
}

// 初始化一个需要优雅停止的队列list。
// 启动时就等待系统的中断、停止信号，一旦受到信号就依次停止各个注册的函数。
func init() {
	defaultCloserList.registerList = list.New()
	shutdownSemaphore = make(chan string)
	go waitingSignal()
}

// 内部的等待终止信号的协程，不能阻塞住init函数。
func waitingSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	log.Println("[CloseHub] Waiting for singal to shutdown program gracefully ...")
	s := <-c
	log.Printf("----------------------------------------\n")
	log.Printf("[CloseHub] Got signal: %s, start shutdown process: \n", s)
	log.Printf("----------------------------------------\n")
	log.Println()

	for e := defaultCloserList.registerList.Front(); e != nil; e = e.Next() {
		closer := e.Value.(GracefulClose)
		log.Printf("[CloseHub] start close %T", closer)
		closer.OnShutdown()
		log.Printf("[CloseHub] %T closed", closer)
		log.Println()

	}

	log.Println("----------------------------------------")
	log.Println("[CloseHub] shutdown process complete")
	log.Println("----------------------------------------")

	shutdownSemaphore <- "ok"

}
