package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	lock1 = sync.RWMutex{}
	lock2 = sync.RWMutex{}
)

func blueRobot() {
	for {
		fmt.Println("Blue: Acquiring lock1")
		lock1.Lock()
		fmt.Println("Blue: Acquiring lock2")
		lock2.Lock()
		fmt.Println("Blue: Both locks acquired")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("Blue: Locks released")
	}
}

func redRobot() {
	for {
		fmt.Println("Red: Acquiring lock1")
		lock2.Lock()
		fmt.Println("Red: Acquiring lock2")
		lock1.Lock()
		fmt.Println("Red: Both locks acquired")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("Red: Locks released")
	}
}

func main() {
	go redRobot()
	go blueRobot()
	time.Sleep(20 * time.Second)
	fmt.Println("Done")
}
