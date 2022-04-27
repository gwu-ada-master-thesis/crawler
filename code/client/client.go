package client

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"strings"
	"time"

	// custom
	"masterthesis/database"
	"masterthesis/server"
	"masterthesis/custom"
)


func help() {
	fmt.Println("Available Commands:")
	fmt.Println("\thelp\t--\tlists all possible commands")
	fmt.Println("\tcrawl\t--\tlogin to as a client with public/private key")
	fmt.Println("\tstore updates\t--\tregister user and generate rsa and symmetric keys")
	fmt.Println("\visualize\t--\tregister user and generate rsa and symmetric keys")
	fmt.Println("\texit\t--\tterminate the process")
}


func Client() {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("    >> ")

		text, _ := reader.ReadString('\n')

		// remove last newline symbol
		text = strings.Replace(text, "\n", "", -1)

		switch strings.ReplaceAll(strings.ToLower(text), " ", "") {

		case "crawl":
			cmdSRead := exec.Command("go", "run", "../main/servermaster.go") // for testing same hash is given
			cmdNRead := exec.Command("go", "run", "../main/nodeworker.go")

			errSRead := cmdSRead.Start()
			errNRead := cmdNRead.Start()

			errSRead = cmdSRead.Wait()
			errNRead = cmdNRead.Wait()

			log.Printf("Master Command finished with error: %v", errSRead)
			log.Printf("Master Command finished with error: %v", errNRead)

			fmt.Println("Readed")
		case "write":
			cmdServerWrite := exec.Command("go", "run", "../main/servermaster.go", "../files/pg-being_ernest.txt", "pg-frankenstein.txt")
			cmdNodeWrite   := exec.Command("go", "run", "../main/nodeworker.go")

			errSWrite := cmdServerWrite.Start()
			errNWrite := cmdNodeWrite.Start()

			errSWrite = cmdServerWrite.Wait()
			errNWrite = cmdNodeWrite.Wait()

			log.Printf("Master Command finished with error: %v", errSWrite)
			log.Printf("Master Command finished with error: %v", errNWrite)

			fmt.Println("Write")
		case "help":
			help()
		case "exit":
			fmt.Println("Quiting")
			time.Sleep(1 * time.Second)

			// exit
			os.Exit(0)
		default:
			fmt.Println("Such command not exist")

		}
	}

}

func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := server.ServerSock()
	c, err := rpc.DialHTTP("unix", sockname)

	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	if err.Error() != "unexpected EOF" {
		fmt.Println("error: ", err)
	}
	return false
}
