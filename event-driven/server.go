package main

/*

TODO
the server has either two implementations
1) pooling
	- every message sent to the server has to be stored in long list
	- a client may ask for this list or a slice of it to fetch the updates


2) event-driven [BONUS]
	- a server is more like a coordinator
	- the server waits for clients wanting to register themselves as listeners
	- a client sends a message by calling an rpc responsible for broadcasting, the client calls a function to loop on all registered clients and send his own message to each of them separately
	- the client on the other side has a server listening for messages being pushed
*/
import (
	"fmt"
	"log"
	"net"
	rpc "net/rpc"
	commons "rpc_assign/commons"
)

var messages []commons.MessageInfo

var userDatabase = map[string]string{}

type Listener int

func (l *Listener) ProdcastMessage(MessInfo commons.MessageInfo, reply *bool) error {

	messages = append(messages, MessInfo)
	fmt.Printf("Message#%03d: (%v) >>> %v\n",
		len(messages), MessInfo.UserData.NickName, MessInfo.Message)

	for key, _ := range userDatabase {
		TargeAddress := key
		CallerAddress := MessInfo.UserData.PublicIp + ":" + MessInfo.UserData.Port
		if TargeAddress != CallerAddress {
			// send message on this address
			data := MessInfo
			var reply bool
			go func() {
				client, err := rpc.Dial("tcp", TargeAddress)
				if err != nil {
					log.Fatal(err)
				}

				error1 := client.Call("Listener.PrintMessage", data, &reply)
				if error1 != nil {
					log.Fatalln("Failed to call the client: ", error1)
				}

				defer client.Close()
			}()
		}
	}
	*reply = true
	return nil
}

func (l *Listener) RegisterUser(user commons.User, reply *bool) error {
	CallerAddress := user.PublicIp + ":" + user.Port
	userDatabase[CallerAddress] = user.NickName
	fmt.Printf("New User is registered#%0d: %v @ %v\n", len(userDatabase), user.NickName, CallerAddress)
	// fmt.Printf("database#%v",userDatabase)
	for key, _ := range userDatabase {
		TargeAddress := key
		if TargeAddress != CallerAddress {
			// send message on this address
			message := user.NickName + " has joined the Chat Room."
			data := commons.MessageInfo{message, user}
			var reply bool
			go func() {
				client, err := rpc.Dial("tcp", TargeAddress)
				if err != nil {
					log.Fatal(err)
				}

				error1 := client.Call("Listener.PrintMessage", data, &reply)
				if error1 != nil {
					log.Fatalln("Failed to call the client: ", error1)
				}

				defer client.Close()
			}()
		}
	}
	*reply = true
	return nil
}

func main() {
	serverAddress := "0.0.0.0:8000"
	addr, err := net.ResolveTCPAddr("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Faild to listen on port 44444: %v", err)
	}

	fmt.Printf("Server is working on Address: %v<<>>%v\n", serverAddress, addr)

	inbound, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	listener := new(Listener) // allocate memory
	rpc.Register(listener)    // 1. local object for procedures
	rpc.Accept(inbound)       // 2. network
}
