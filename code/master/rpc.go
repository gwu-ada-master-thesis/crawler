package master


import (
	// build-in
	"log"
	"net"
	"net/http"
	"net/rpc"
	"fmt"
	"os"

	// custom
	"masterthesis/utils"
)


type JobRequest struct {
	Job     *Job
	Status 	string
}


type JobResponse struct {
	Job     *Job
	Status 	string
}


func (master *Master) serve() {

	// register master
	rpc.Register(master)
	
	// handle HTTP requests
	rpc.HandleHTTP()

	// get socket name
	socketName := utils.GenerateSocketName(824, "master")
	
	// remove previously created socket
	os.Remove(socketName)

	// listen to the socket
	listener, err := net.Listen("unix", socketName)

	// check for error
	if err != nil {
		log.Fatal("listen error:", err)
	}

	// initialize thread for listener
	go http.Serve(listener, nil)
}


func Call(rpcName string, request interface{}, response interface{}) bool {
	
	// get socket name
	socketName := utils.GenerateSocketName(824, "master")

	// initiate connection
	connection, err := rpc.DialHTTP("unix", socketName)

	// check for error
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// close when all work is done
	defer connection.Close()

	// initiate call
	err = connection.Call(rpcName, request, response)

	if err == nil {
		return true
	}

	if err.Error() != "unexpected EOF" {
		fmt.Println("error: ", err)
	}
	
	return false
}

