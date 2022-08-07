
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"bytes"
	"strconv"
)

type user struct {
	messages []string
	channels []string
}

type channel struct {
	messages []string
	members int
}

var userMap = make(map[string]user)
var channelMap = make(map[string]channel)

var (
	DELIMITER = []byte(`\r\n`)
)

func reg(args []byte) string {
	username := bytes.TrimSpace(args)
	if username[0] != '@' {
		return err("Username must begin with @")
	}
	if len(username) == 0 {
		return err("Username cannot be blank")
	}

	var u user
	userMap[string(username)] = u
	//fmt.Println(userSet)

	return ok()
}

func join(args []byte, currentUser string) string{
	channelID := bytes.TrimSpace(args)
	if channelID[0] != '#' {
		return err("Channel ID must begin with #")
	}

	if entry, ok := channelMap[string(channelID)]; ok {
		entry.members = entry.members + 1
		channelMap[string(channelID)] = entry
	} else {
		var cm channel
		cm.members = 1
		channelMap[string(channelID)] = cm
	}

	if entry, ok := userMap[currentUser]; ok {
		entry.channels = append(entry.channels, string(channelID))
		userMap[currentUser] = entry
	}
	return ok()
}

func leave(args []byte, currentUser string) string{
	channelID := bytes.TrimSpace(args)
	if channelID[0] != '#' {
		return err("Channel ID must begin with #")
	}

	if entry, ok := channelMap[string(channelID)]; ok {
		entry.members = entry.members - 1

		if entry.members == 0 {
			delete(channelMap, string(channelID))
		} else {
			channelMap[string(channelID)] = entry
		}
	}

	userChannelIDs := userMap[currentUser].channels

	for i := 0; i < len(userChannelIDs); i++ {
		userChannelID := userChannelIDs[i]
		if userChannelID == string(channelID) {
			ret := make([]string, 0)
			ret = append(ret, userChannelIDs[:i]...)
			ret = append(ret, userChannelIDs[i+1:]...)

			if entry, ok := userMap[currentUser]; ok {

				// Then we modify the copy
				entry.channels = ret
			 
				// Then we reassign map entry
				userMap[currentUser] = entry
			}

			return ok() 
		} 
	
	}
	return err("Channel not found")
}

func chns(currentUser string) string{
	
	channelList := ""
	for channelID, _ := range channelMap {
		channelList += channelID + " "
	}

	if len(channelList) == 0 {
		return err("no channels found")
	}
	return strings.TrimSpace(channelList) + "\n"
}

func usrs(currentUser string) string{
	
	userList := ""

	for userID, _ := range userMap {
		userList += userID + " "
	}

	if len(userList) == 0 {
		return err("no users found")
	}
	return strings.TrimSpace(userList) + "\n"
}

func msg(args []byte, currentUser string) string{
	args = bytes.TrimSpace(args)
	if args[0] != '#' && args[0] != '@' {
		return err("recipient must be a channel ('#name') or user ('@user')")
	}

	recipient := bytes.Split(args, []byte(" "))[0]
	if len(recipient) == 0 {
		return err("recipient must have a name")
	}

	args = bytes.TrimSpace(bytes.TrimPrefix(args, recipient))
	l := bytes.Split(args, DELIMITER)[0]
	length, e := strconv.Atoi(string(l))
	if e != nil {
		return err("body length must be present")

	}
	if length == 0 {
		return err("body length must be at least 1")
	}

	padding := len(l) + len(DELIMITER) // Size of the body length + the delimiter
	body := args[padding : padding+length]

	if recipient[0] == '#' {
		if entry, ok := channelMap[string(recipient)]; ok {
			entry.messages = append(entry.messages, string(body))
			channelMap[string(recipient)] = entry
		}
	}

	if recipient[0] == '@' {
		if entry, ok := userMap[string(recipient)]; ok {
			entry.messages = append(entry.messages, string(body))
			userMap[string(recipient)] = entry
		}
	}

	return "\n"
}

func handle(message []byte, currentUser string) string {
	cmd := bytes.ToUpper(bytes.TrimSpace(bytes.Split(message, []byte(" "))[0]))
	args := bytes.TrimSpace(bytes.TrimPrefix(message, cmd))

	switch string(cmd) {
	case "REG":
		response := reg(args)
		return response
	
	case "JOIN":
		response := join(args, currentUser)
		return response
	
	case "LEAVE":
		response := leave(args, currentUser)
		return response

	case "CHNS":
		response := chns(currentUser)
		return response

	case "USRS":
		response := usrs(currentUser)
		return response

	case "MSG":
		response := msg(args, currentUser)
		return response
	default:
		return ""
	}
}


func err(error string) string {
	return "ERR " + error + "\n\n"
}

func ok() string {
	return "OK\n\n"
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}

	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	currentUser := ""
	for {
		//get the msg from the client
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		if strings.HasPrefix(netData, "currentUser:") {
			currentUserData := strings.Split(netData, " ")[0]
			currentUser = strings.Split(currentUserData, ":")[1]
			
			netData = string(netData[len(currentUserData) : len(netData)])
		}

		if strings.TrimSpace(string(netData)) == "STOP" {
			for userID, _ := range userMap {
				fmt.Println(userID)

				for i := 0; i < len(userMap[userID].messages); i++ {
					fmt.Println(userMap[userID].messages[i])
				}

				for i := 0; i < len(userMap[userID].channels); i++ {
					fmt.Println(userMap[userID].channels[i])

					for j := 0; j < len(channelMap[userMap[userID].channels[i]].messages); j++ {
						fmt.Println(channelMap[userMap[userID].channels[i]].messages[j])
					}
				}

				fmt.Println("\n")
			}
			
			return
		}

		netData = strings.TrimSpace(netData)
		currentUser = strings.TrimSpace(currentUser)

		//handle the command from the client
		response := handle([]byte(netData), currentUser)

		
		c.Write([]byte(response))

	}
}