package chat

import (
	"bufio"
	"net"

	m "net-cat/models"
)

// check nick
func isNickCorrect(room *m.Room, txt string) bool {
	if txt == "" {
		return false
	}
	if _, ok := room.UsersList[txt]; ok {
		return false
	}
	return true
}

// rename - changes name of user to new
// to change name user have to send "--rename" message
func rename(c net.Conn, ch *m.Room, prevName string) string {
	for {
		c.Write([]byte("[ENTER YOUR NEW NAME]:"))
		rd := bufio.NewReader(c)
		nick, err := rd.ReadString('\n')

		if err != nil {
			c.Write([]byte("Wrong name! can not create user with this name\n try again\n"))
		} else {
			nick = nick[:len(nick)-1]

			if nick == "" {
				continue
			}

			if isNickCorrect(ch, nick) {
				return nick
			} else {
				c.Write([]byte("Wrong name! can not create user with this name\n try again\n"))
			}
		}
	}
}
