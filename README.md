# Chat-Room
client server applicaiont based on RPC Mechanism and Golang


## server implementations
### the server has either two implementations
1) pooling
	- every message sent to the server has to be stored in long list
	- a client may ask for this list or a slice of it to fetch the updates


2) event-driven
	- a server is more like a coordinator
	- the server waits for clients wanting to register themselves as listeners
	- a client sends a message by calling an rpc responsible for broadcasting, the client calls a function to loop on all registered clients and send his own message to each of them separately
	- the client on the other side has a server listening for messages being pushed
