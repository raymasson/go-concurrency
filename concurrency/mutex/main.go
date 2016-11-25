package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	//mutex := new(sync.Mutex)
	runtime.GOMAXPROCS(4)

	f, _ := os.Create("./log.txt")
	f.Close()

	logCh := make(chan string, 50)

	go func() {
		for {
			msg, ok := <-logCh

			if ok {
				f, _ := os.OpenFile("./log.txt", os.O_APPEND, os.ModeAppend)

				logTime := time.Now().Format(time.RFC3339)
				f.WriteString(logTime + " - " + msg)
				f.Close()
			} else {
				break
			}
		}
	}()

	mutex := make(chan bool, 1)

	for i := 1; i < 10; i++ {
		for j := 1; j < 10; j++ {
			//mutex.Lock()
			mutex <- true
			go func() {
				//fmt.Printf("%d + %d = %d\n", i, j, i+j)
				//mutex.Unlock()
				msg := fmt.Sprintf("%d + %d = %d\n", i, j, i+j)
				logCh <- msg
				fmt.Print(msg)
				<-mutex
			}()
		}
	}

	fmt.Scanln()
}
