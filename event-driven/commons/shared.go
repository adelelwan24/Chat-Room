package commons

/*
TODO
you may define any structs here to be used by the rpc
*/

// hint: you will need to have the server address fixed between clients and the coordinating server
func Get_server_address() string {
	return "0.0.0.0:9999"
}


type User struct{
	NickName string
	Port string
	PublicIp string
	
}

type MessageInfo struct{
	Message string
	UserData User
}

type Args struct{
	Message string
	NickName string
	Index int
}

type Res struct {
	Messages []MessageInfo
	LastMessageIndex int
}