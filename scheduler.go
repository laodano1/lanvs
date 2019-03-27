package main

import (
	"time"
	"fmt"
	"math/rand"
	"flag"
	"strconv"
)

var (
	// worker amount
	workerNum int

	// task amount
	taskNum int

	// repeat number
	repeat int
)

type (
	// task item
	TaskItem struct {
		TaskId int
		Body   string
	}
)

// task worker
func TaskWorker(taskQ chan TaskItem, partDone chan bool, allDone chan bool, workerId int)  {
	c := time.Tick(100 * time.Millisecond)
	for {
		select {
		case oneTask := <- taskQ :
			rand.Seed(time.Now().UnixNano())
			i := rand.Intn(4)
			// sleep
			time.Sleep( time.Duration(i) * time.Second )
			//fmt.Printf("random is %d second.\n", i)
			fmt.Printf("\nstart to task '%s' in worker '%d'\n", oneTask.Body, workerId)
			//
			//// http request
			//fmt.Printf("=======> do task '%s'.\n", oneTask.Body)
			fmt.Printf("do http request in task '%s' in worker '%d'.\n", oneTask.Body, workerId)
			//
			//// response parse
			fmt.Printf("response parse in task '%s' in worker '%d'\n.", oneTask.Body, workerId)
			//
			//// inserting to DB
			fmt.Printf("inserting to db in task '%s' in worker '%d'\n", oneTask.Body, workerId)

			fmt.Printf("------> task '%s' _done. in worker '%d' \n", oneTask.Body, workerId)

			if len(taskQ) <= workerNum * 2  {
				//time.Sleep(2 * time.Second)
				//fmt.Printf("task '%d' is last task.\n", oneTask.TaskId)
				//fmt.Printf("cap: %d/%d. need to add task.\n", len(taskQ), cap(taskQ))
				partDone <- true
			}

			if len(taskQ) <= 0 {
				//close(taskQ)
				//fmt.Println("this is last task!!!! @@@@@@@@@@@")

				time.Sleep(3 * time.Second)
				allDone <- true
				//wg.Done()
			}

		case <- c :
			time.Sleep(200 * time.Millisecond)
			//fmt.Println("this is idel.---------")
		}
	}
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
	flag.IntVar(&repeat,   "rp", 2, "repeat amount")
	flag.Parse()

	fmt.Printf("wokernum: %d\n", workerNum)
	fmt.Printf("taskNum: %d\n", taskNum)

	taskQ  := make(chan TaskItem, taskNum)
	isPartDone := make(chan bool)
	isAllDone  := make(chan bool)


	// get tasks and fill task list
	FillTaskQueue(taskQ, taskNum, "1st")

	// init go routine pool
	for i := 0; i < workerNum; i++ {
		go TaskWorker(taskQ, isPartDone, isAllDone, i)
	}

	// infinite loop to do tasks
	turn := 0
	for {
		// check tasklist empty or not
		select {
		case <- isPartDone:
			//time.Sleep(2 * time.Second)
			if turn < repeat {
				fmt.Printf("----> %d tasks done. Fill task to queue in %dth time!\n", workerNum, turn + 1)
				go FillTaskQueue(taskQ, 10, "added" + strconv.Itoa(turn))
			} else {

			}
			turn++
		case <- isAllDone :
			time.Sleep(1 * time.Second)
			fmt.Printf("channel cap: %d/%d. Bye bye!\n", len(taskQ), taskNum)

			return
		}
	}
}
