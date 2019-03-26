package main

import (
	"time"
	"fmt"
	"math/rand"
	"flag"
	"strconv"
	"sync"
)

var (
	// worker amount
	workerNum int

	// task amount
	taskNum int
)

type (
	// task item
	TaskItem struct {
		TaskId int
		Body   string
	}
)

// task worker
func TaskWorker(taskQ chan TaskItem, partDone chan bool, allDone chan bool, wg *sync.WaitGroup)  {
	c := time.Tick(3 * time.Second)
	for {
		select {
		case oneTask := <- taskQ :
			rand.Seed(time.Now().UnixNano())
			i := rand.Intn(3)
			// sleep
			time.Sleep( time.Duration(i) * time.Second )
			//fmt.Printf("random is %d second.\n", i)
			fmt.Printf("\nstart to task '%s'\n", oneTask.Body)
			//
			//// http request
			//fmt.Printf("=======> do task '%s'.\n", oneTask.Body)
			fmt.Println("do http request.")
			//
			//// response parse
			//fmt.Println("response parse.")
			//
			//// inserting to DB
			fmt.Println("inserting to db.")

			fmt.Printf("------> task '%s' _done.\n", oneTask.Body)

			if len(taskQ) <= workerNum * 2  {
				//time.Sleep(2 * time.Second)
				//fmt.Printf("task '%d' is last task.\n", oneTask.TaskId)
				fmt.Printf("cap: %d/%d. need to add task.\n", len(taskQ), cap(taskQ))
				partDone <- true
			}

			if len(taskQ) <= 0 {
				//close(taskQ)
				fmt.Println("this is last task!!!! @@@@@@@@@@@")
				time.Sleep(8 * time.Second)
				allDone <- true
			}
			//if len(taskQ) <= 0 {
			//	//wg.Done()
			//	fmt.Println("task terminate.")
			//	goto ter
			//}
		case <- c :
			time.Sleep(2 * time.Second)
			fmt.Println("this is idel.---------")
		}
	}
	//ter:
}

// fill task queue
func FillTaskQueue(taskQ chan TaskItem, tnum int, prefix string) {

	for i := 0; i < tnum; i++ {
		var tmpT TaskItem
		tmpT.TaskId = i + 1
		tmpT.Body = prefix + "-task-" + strconv.Itoa(i + 1)
		taskQ <- tmpT
		fmt.Printf("task: %s is filled to channel.\n", tmpT.Body)
	}
}

func main() {

	flag.IntVar(&workerNum, "wk", 5, "worker number in pool")
	flag.IntVar(&taskNum,   "tn", 10, "task amount")
	flag.Parse()

	fmt.Printf("wokernum: %d\n", workerNum)
	fmt.Printf("taskNum: %d\n", taskNum)

	taskQ  := make(chan TaskItem, taskNum)
	isPartDone := make(chan bool)
	isAllDone  := make(chan bool)
	var wg sync.WaitGroup

	wg.Add(workerNum)
	// get tasks and fill task list
	FillTaskQueue(taskQ, taskNum, "1st")

	// init go routine pool
	for i := 0; i < workerNum; i++ {
		go TaskWorker(taskQ, isPartDone, isAllDone, &wg)
	}

	// infinite loop to do tasks
	turn := 0
	for {
		// check tasklist empty or not
		select {
		case <- isPartDone:
			//time.Sleep(2 * time.Second)
			if turn < 1 {
				fmt.Printf("----> %d tasks done. Fill task to queue!\n", workerNum)
				FillTaskQueue(taskQ, 5, "pre" + strconv.Itoa(turn))
			} else {
				//fmt.Printf("channel cap: %d/%d. Bye bye!\n", len(taskQ), workerNum * 3)
				//wg.Wait()
				//time.Sleep(20 * time.Second)
				//return
				//close(taskQ)
			}
			turn++
		case <- isAllDone :
			time.Sleep(3 * time.Second)
			fmt.Printf("channel cap: %d/%d. Bye bye!\n", len(taskQ), workerNum * 3)
			return
		}
	}
}