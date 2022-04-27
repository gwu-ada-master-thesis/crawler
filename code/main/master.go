package main

import (
	// built-in
	"fmt"
	"os"
	"time"
	"log"

	// 3rd party imports

	// custom
	"masterthesis/utils"
	"masterthesis/master"
	

)

func main() {
	
	fmt.Println(`
					*****************************************		
					*                                       *
					*                MASTER                 *
					*                                       *
					*****************************************
	`)

	// check whether initiation has been done with all parameters
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "InputError: input file path is missing.\n")
		os.Exit(1)
	}

	// retrive filepath
	filepath := os.Args[1]

	// read tasks
	tasks, err := utils.ReadLines(filepath)
	if err != nil {
		panic(err)
	}

	// print master status
	log.Printf("initializing\n")

	// initiate server
	master := master.MakeMaster(tasks)
	time.Sleep(1 * time.Second)

	for master.Done() == false {
		// print master status
		log.Printf("serving\n");time.Sleep(5 * time.Second)
	}
	
	log.Println("letting all workers to finalize");time.Sleep(5 * time.Second)
	
	// print master status
	log.Printf("turning off\n");time.Sleep(1 * time.Second)
}
