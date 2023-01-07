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
	"math/rand"
	"net"
	rpc "net/rpc"
	"os"
	commons "rpc_assign/commons"
	"strconv"
	"strings"
	"sync"
	"time"
)

var messages []commons.MessageInfo

func (l *Listener) PrintMessage(Data commons.MessageInfo, reply *bool) error {
	messages = append(messages, Data)
	fmt.Printf("\r(%v) >>> %v\n", Data.UserData.NickName, Data.Message)
	fmt.Printf("(%v) >>>", nickname)
	*reply = true
	return nil
}

type Listener int

var wg = sync.WaitGroup{}
var nickname string

func main() {
	//  server

	max := 60000
	min := 45555
	// set seed
	rand.Seed(time.Now().UnixNano())
	// generate random number and print on console
	port_int := rand.Intn(max-min) + min
	port := strconv.Itoa(port_int)

	fmt.Println("Working on port:", port)

	// TODO: make it get the public ip dynamiclly
	ip := "197.63.157.57"
	wg.Add(2)

	go func(port string) {
		address := "0.0.0.0" + ":" + port
		addr, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			log.Fatalf("Faild to listen on port 55555: %v", err)
		}

		fmt.Printf("Server is working on Address: %v\n", address)

		inbound, err := net.ListenTCP("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}

		listener := new(Listener) // allocate memory
		rpc.Register(listener)    // 1. local object for procedures
		rpc.Accept(inbound)       // 2. network}

		wg.Done()
	}(port)





	// client
	go func(ip string, port string) {
		client, err := rpc.Dial("tcp", "0.0.0.0:8000")
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
		in := bufio.NewReader(os.Stdin)

		fmt.Println("==========================Welcome=========================")
		fmt.Println("=========================Chat Room========================")
		fmt.Printf("please, Enter your nickname:  ")

		nickname, err = in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		nickname = strings.TrimSpace(nickname)

		// call the coordination server to register this client

		userData := commons.User{nickname, port, ip}
		var registrationReply bool
		err = client.Call("Listener.RegisterUser", userData, &registrationReply)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("==========================================================")

		// fmt.Println("To exit the chat room, Please Enter: exit()")
		// fmt.Println("To print all messages Enter: printall()")
		fmt.Println("please, write any message to be sent to the server")

		var reply bool
		// var lastMessageIndex int = 0

		for {
			fmt.Printf("(%v) >>>", nickname)
			message, err := in.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			m := strings.TrimSpace(message)
			if m == "exit()" {
				break
			}
			// if m == "printall()" {
			// 	lastMessageIndex = 0
			// }

			data := commons.MessageInfo{m, userData}

			error1 := client.Call("Listener.ProdcastMessage", data, &reply)

			if error1 != nil {
				log.Fatalln("Failed to call the client: ", error1)
			}

			// lastMessageIndex = reply.LastMessageIndex

			// for _, val := range reply.Messages {
			// 	fmt.Printf("(%v) >>> %v \n",val.UserNickName, val.Message)
			// }
			// reply.Messages = nil
		}

		wg.Done()
	}(ip, port)

	wg.Wait()

	fmt.Println("==========================THANKS==========================")

}
