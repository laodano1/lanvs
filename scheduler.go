package main

import (
	"time"
	"fmt"
	"math/rand"
	"flag"
	"strconv"
	"bufio"
	"os"
	"strings"
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
func TaskWorker(taskQ chan TaskItem, isDone chan bool)  {

	for {
		select {
		case oneTask := <- taskQ :
			rand.Seed(time.Now().UnixNano())
			i := rand.Intn(8)
			// sleep
			time.Sleep( time.Duration(i) * time.Second )
			//fmt.Printf("random is %d second.\n", i)
			//fmt.Printf("start to task %d\n", oneTask.TaskId)
			//
			//// http request
			//fmt.Printf("=======> do task '%s'.\n", oneTask.Body)
			//fmt.Println("do http request.")
			//
			//// response parse
			//fmt.Println("response parse.")
			//
			//// inserting to DB
			//fmt.Println("inserting to db.")

			fmt.Printf("------> task '%s' done.\n", oneTask.Body)

			if oneTask.TaskId == taskNum {
				time.Sleep(8 * time.Second)
				fmt.Printf("task '%d' is last task.\n", oneTask.TaskId)
				isDone <- true
			}
		}
	}
}

// fill task queue
func FillTaskQueue(taskQ chan TaskItem) {

	for i := 0; i < taskNum; i++ {
		var tmpT TaskItem
		tmpT.TaskId = i + 1
		tmpT.Body = "task-" + strconv.Itoa(i + 1)
		taskQ <- tmpT
	}
}

func HandleStdin(echoInfo string) string {
	rd := bufio.NewReader(os.Stdin)
	fmt.Printf(echoInfo)
	isAgain, err := rd.ReadString('\n')
	if err != nil {
		fmt.Println("reader error:", err)
	}

	isAgain = strings.TrimSpace(isAgain)
	return isAgain
}


func main() {

	flag.IntVar(&workerNum, "wk", 10, "worker number in pool")
	flag.IntVar(&taskNum, "tn", 50, "task amount")
	flag.Parse()

	fmt.Printf("wokernum: %d\n", workerNum)
	fmt.Printf("taskNum: %d\n", taskNum)

	taskQ  := make(chan TaskItem, workerNum)
	isDone := make(chan bool)

	// get tasks and fill task list
	go FillTaskQueue(taskQ)

	// init go routine pool
	for i := 0; i < workerNum; i++ {
		go TaskWorker(taskQ, isDone)
	}

	// infinite loop to do tasks
	for {
		// check tasklist empty or not
		select {
		case <- isDone:
			time.Sleep(3 * time.Second)
			fmt.Println("All tasks of this turn done. Query DB to fill task queue!")

			for {
				isAgain := HandleStdin("if do more tasks [Y|N]: ")
				isAgain = strings.ToUpper(isAgain )

				switch isAgain {
				case "Y", "YES", "":
					fmt.Printf("input is %s\n", isAgain)

				getmeta:
					str := HandleStdin("input worker and task amount [50]: ")
					if str == "" {
						str = "50"
					}
					ta, err := strconv.Atoi(str)
					if err != nil {
						fmt.Println("convert error:", err)
						goto getmeta
					}
					taskNum = ta
					FillTaskQueue(taskQ)

					goto nextSelect

				case "N", "NO" :
					fmt.Println("Bye bye!")
					return
				}
			}
		}
		nextSelect:

	}
}