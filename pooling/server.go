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
	"strings"
	rpc "net/rpc"
	commons "rpc_assign/commons"
)

var messages []string

type Listener int

func (l *Listener) SendMessage(Data commons.Args, reply *commons.Res) error {
	// can't modify data in the client even if it's a pointer only can send data back throught the reply
	messages = append(messages, Data.Message)
	fmt.Printf("This is the message#%d: (%v) >>> %v\n",len(messages),strings.TrimSpace(Data.NickName), Data.Message)
	if x := len(messages) - Data.Index; x == 1{
		*reply = commons.Res{messages[Data.Index+1:], len(messages)}
		return nil
	}
	*reply = commons.Res{messages[Data.Index:], len(messages)}
	return nil
}


func main() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:44444")
	if err != nil {
		log.Fatalf("Faild to listen on port 44444: %v",err)
	}
	
	inbound, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	
	// messages = append(messages, "adel", "mamoun")
	
	listener := new(Listener) // allocate memory
	rpc.Register(listener)    // 1. local object for procedures
	rpc.Accept(inbound)       // 2. network
}