package main


import (
	// builtin
	"time"
	"fmt"
	"log"

	// custom
	"masterthesis/worker"
)


func main() {

		
	fmt.Println(`
					*****************************************		
					*                                       *
					*                WORKER                 *
					*                                       *
					*****************************************
	`); time.Sleep(1 * time.Second)

	// print worker status
	log.Printf("initializing\n"); time.Sleep(1 * time.Second)

	// create worker
	worker := worker.MakeWorker()

	// print worker status
	log.Printf("working\n"); time.Sleep(1 * time.Second)

	// start worker
	worker.Work()

	// print worker status
	log.Printf("turning off\n"); time.Sleep(1 * time.Second)
}