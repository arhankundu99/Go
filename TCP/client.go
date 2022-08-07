package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	connect := arguments[1]
	c, err := net.Dial("tcp", connect)
	if err != nil {
		fmt.Println(err)
		return
	}

	currentUser := ""
	for {


		reader := bufio.NewReader(os.Stdin)

		text, _ := reader.ReadString('\n')

		text = strings.TrimSpace(string(text))

		if currentUser != "" && !strings.HasPrefix(text, "REG") && text != "STOP"{
			text = "currentUser:" + currentUser + " " + text
		}

		//send the msg to the server
		fmt.Fprintf(c, text+"\n")
		
		//get the message from the server
		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print(message + "\n")

		message = strings.TrimSpace(message)

		if message == "OK" && strings.HasPrefix(text, "REG"){
			currentUser = strings.Split(text, " ")[1]
		}

		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}

		
	}
}