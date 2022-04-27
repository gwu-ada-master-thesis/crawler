package main

import (
	"bufio"
	"fmt"
	"log"
	"masterthesis/client"
	"os"
	"strings"
	"time"
)

func help() {
	fmt.Println("Available Commands:")
	fmt.Println("\thelp\t--\tlists all possible commands")
	fmt.Println("\tcrawl\t--\tlogin to as a client with public/private key")
	fmt.Println("\tstore updates\t--\tregister user and generate rsa and symmetric keys")
	fmt.Println("\visualize\t--\tregister user and generate rsa and symmetric keys")
	fmt.Println("\texit\t--\tterminate the process")
}

func input(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')

	// remove last newline symbol
	input = strings.Replace(input, "\n", "", -1)

	return input
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Hello dear! Type help to see availabe commands")

	for {
		fmt.Print(">> ")

		inputText := input(reader)

		switch strings.ReplaceAll(strings.ToLower(inputText), " ", ""); {
		case "login":
			fmt.Print("Welcome")
			client.Client()
		case "register":
			log.Println("Register")
		case "help":
			help()

		case "exit":
			fmt.Println("Quiting...")
			time.Sleep(2 * time.Second)

			// exit
			os.Exit(0)
		default:
			fmt.Println("Such command does not exist")
		}
	}

}
