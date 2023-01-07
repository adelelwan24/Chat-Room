package main

/*
TODO
there are two implementations
1) pooling
	- the client will dial the rpc of the coordinating server
	- the client will call the remote procedure on the server to send a message
	- the client can fetch all of the messages history from the server using remote procedure call


2) event-driven [BONUS]
	- a client starts by looking for a port to establish it's server on (like giving my phone number to my friends to call me)
	- a client can send a message through an infinite loop waiting for input text, this message will be broadcasted to other clients through an rpc call on the server
	- a client can also receive messages simultaneously using the GO keyword
	- so a client here is a server + a client at the same time
*/

import (
	"bufio"
	"fmt"
	"log"
	rpc "net/rpc"
	"os"
	"strings"
	// "strconv"
	commons "rpc_assign/commons"
)


func main() {
	client, err := rpc.Dial("tcp", "0.0.0.0:44444")
	if err != nil {
		log.Fatal(err)
	}

	in := bufio.NewReader(os.Stdin)

	fmt.Println("==========================Welcome=========================")
	fmt.Println("========================message RPC=======================")
	fmt.Println("=========================Chat Room========================")
	fmt.Println("please, Enter your nickname")

	nickname, err := in.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("please, write any message to be sent to the server")
	fmt.Println("To exit the chat room, Please Enter: exit()")


	var reply commons.Res
	var lastMessageIndex int = 0
	
	for {
		fmt.Printf(">>>")
		message, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		m := strings.TrimSpace(message)
		if m == "exit()" {
			break
		}
		if m == "printall()" {
			lastMessageIndex = 0
		}

		data := commons.Args{m , nickname, lastMessageIndex}

		error1 := client.Call("Listener.SendMessage", data, &reply)

		if error1 != nil {
			log.Fatalln("Failed to call the client: ",error1)
		}
		
		lastMessageIndex = reply.LastMessageIndex

		for _, val := range reply.Messages {
			fmt.Printf(">>> %v \n",val)
		}
		reply.Messages = nil
	}

	fmt.Println("========================THANKS=======================")

}
