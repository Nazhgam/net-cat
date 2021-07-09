package models

import (
	"net"
)

const (
	ChatlistPath = "chatlist.txt" // path to file with actual chatlist (if server will down, next run will restore chats)

	Logo = `
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    '.       | '' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     '-'       '--'
`
)

var (
	Rooms []Room
	Port  = "8989"
)

// Room - struct of chat room
type Room struct {
	Name string
	// path of own log-file (bonus)
	Path           string
	ConnectionChan chan net.Conn
	UserMessage    chan Msg
	Disconect      chan net.Conn
	AccountsOfUser map[net.Conn]string
	UsersList      map[string]bool
}

// Msg - struct of user's message
type Msg struct {
	Text string
	From string
	Room string
}
