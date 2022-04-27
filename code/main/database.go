package main

import (
	// built-in
	"fmt"
	"log"
	"time"
	
	// 3rd party

	// custom
	"masterthesis/database"
)

func main() {

	fmt.Println(`
					                        *******************************************		
				 	                        *                                         *
				 	                        *                DATABASE                 *
				 	                        *                                         *
				 	                        *******************************************
	`); time.Sleep(1 * time.Second)

	database.MakeDatabase()

	log.Printf("starting"); time.Sleep(time.Second)

	// to keep the database running
	for { 
		time.Sleep(time.Second)
	}

	log.Printf("turning off"); time.Sleep(time.Second)
}
