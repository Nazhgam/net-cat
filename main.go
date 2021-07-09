package main

import (
	"fmt"
	"net"
	"os"

	u "net-cat/Utilities"
	chat "net-cat/chat"
	m "net-cat/models"
)

// the program implements TCP-chat as net-cat
// user can choose existing chat room or create new own
// user can change his name without exiting room
// every room limited by 10 users
// every chat has own log-file
// when last user exits chat, chat will be closed, but history will be saved in log-file
func main() {
	defer u.RecoverFunction()
	if len(os.Args) > 2 {
		fmt.Println("Invalid number of arguments")
		fmt.Println("[USAGE]: go run . [ $port ]")
		return
	}
	if len(os.Args) == 2 {
		n := u.Atoi(os.Args[1])
		if n < 1 || n > 65535 {
			fmt.Println("Invalid port number")
			fmt.Println("[USAGE]: go run . [ $port ]")
			return
		}
		m.Port = fmt.Sprint(n)
	}

	ln, err := net.Listen("tcp", "localhost:"+m.Port)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Listening on port :", m.Port)
	chat.ChatsRecover()
	chat.AcceptOfuser(ln)
}
