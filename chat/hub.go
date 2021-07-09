package chat

import (
	"bufio"
	"net"
	u "net-cat/Utilities"
	m "net-cat/models"
	"os"
	"strings"
	"time"
)

func ChatsRecover() {
	list := strings.Split(u.FileReader(m.ChatlistPath), "\n")
	for _, chat := range list {
		i := strings.LastIndex(chat, "_")
		if chat == "" || i == -1 {
			continue
		}
		name := chat[:i]
		room := m.Room{
			Name:           name,
			Path:           chat,
			ConnectionChan: make(chan net.Conn),
			UserMessage:    make(chan m.Msg),
			Disconect:      make(chan net.Conn),
			AccountsOfUser: make(map[net.Conn]string),
			UsersList:      make(map[string]bool),
		}
		m.Rooms = append(m.Rooms, room)
		go startHub(&room)
	}
}

// entry point to new connection
func AcceptOfuser(ln net.Listener) {
	// waiting for new users
	for {
		conn, err := ln.Accept()
		if err != nil {
			break
		}

		go login(conn)
	}
}

// user self identification and choosing the chat
func login(conn net.Conn) {

	defer u.RecoverFunction()
	conn.Write([]byte("Welcome to TCP-Chat!\n"))
	conn.Write([]byte(m.Logo))
	ch := SelectRoom(conn) // choosing chat

	if ch == nil {
		return
	}

	conn.Write([]byte("[ENTER YOUR NAME]:"))

	for {
		rd := bufio.NewReader(conn)
		nick, err := rd.ReadString('\n')
		if err != nil {
			continue
		}

		nick = nick[:len(nick)-1]
		if isNickCorrect(ch, nick) {

			u.WriteToArchive(ch, nick+" has joined "+ch.Name+" chat at "+time.Now().String()[:19]+"\n")
			for k, v := range ch.AccountsOfUser {
				k.Write([]byte("\r" + nick + " has joined " + ch.Name + " chat at " + time.Now().String()[:19] + "\n"))
				k.Write([]byte("[" + time.Now().String()[:19] + "][" + v + "]:"))
			}
			st := u.FileReader(ch.Path) // getting old meassages
			conn.Write([]byte(st))
			// prepare to next message
			conn.Write([]byte("[" + time.Now().String()[:19] + "][" + nick + "]:"))
			ch.AccountsOfUser[conn] = nick
			ch.UsersList[nick] = true
			ch.ConnectionChan <- conn
			break
		}
		conn.Write([]byte("Wrong name! can not create user with this name\n try another\n"))
	}
}

// starting new chat room
func startHub(ch *m.Room) {
	defer u.RecoverFunction()
	for {
		select {
		case conn := <-ch.ConnectionChan: // new connection to room
			go handle(conn, ch)

		case msg := <-ch.UserMessage: // new message to room
			u.WriteToArchive(ch, msg.Text)
			for k, v := range ch.AccountsOfUser { // sends msg to other users
				if v != msg.From { // to other users
					k.Write([]byte(msg.Text))
					k.Write([]byte("[" + time.Now().String()[:19] + "][" + v + "]:"))
				} else { // to sender
					k.Write([]byte("[" + time.Now().String()[:19] + "][" + v + "]:"))
				}
			}

		case dc := <-ch.Disconect: // disconnect lost users
			name := ch.AccountsOfUser[dc]
			delete(ch.AccountsOfUser, dc)
			delete(ch.UsersList, name)
			for k, v := range ch.AccountsOfUser {
				k.Write([]byte("\r" + name + " has left " + ch.Name + " chat at " + time.Now().String()[:19] + "\n"))
				k.Write([]byte("[" + time.Now().String()[:19] + "][" + v + "]:"))
			}
			u.WriteToArchive(ch, name+" has left "+ch.Name+" chat at "+time.Now().String()[:19]+"\n")
			if len(ch.AccountsOfUser) == 0 { // if last user
				var actualList string
				c := -1
				for i, room := range m.Rooms {
					if room.Name != ch.Name {
						actualList += room.Path + "\n"
					} else {
						c = i
					}
				}
				if c >= 0 {
					m.Rooms = append(m.Rooms[:c], m.Rooms[c+1:]...)
				}
				os.Truncate(m.ChatlistPath, 0)
				u.WriteToArchive(&m.Room{Path: m.ChatlistPath}, actualList)

				return
			}
		}
	}
}

// connection to room handling
func handle(c net.Conn, ch *m.Room) {
	defer u.RecoverFunction()
	for {
		rd := bufio.NewReader(c)
		message, err := rd.ReadString('\n')
		if err != nil {
			break
		}
		if message == "\n" {
			continue
		}

		msg := m.Msg{
			Text: message,
			From: ch.AccountsOfUser[c],
			Room: ch.Name,
		}
		if msg.Text != "--rename\n" {
			msg.Text = "\r[" + time.Now().String()[:19] + "]" + "[" + ch.AccountsOfUser[c] + "]:" + message
		} else {
			prevName := ch.AccountsOfUser[c]
			delete(ch.AccountsOfUser, c)
			delete(ch.UsersList, prevName)
			newName := rename(c, ch, prevName)
			ch.AccountsOfUser[c] = newName
			ch.UsersList[newName] = true
			if newName == prevName {
				c.Write([]byte("[" + time.Now().String()[:19] + "][" + msg.From + "]:"))
				continue
			}
			msg.Text = "\r" + prevName + " changed his name to " + ch.AccountsOfUser[c] + " at " + time.Now().String()[:19] + "\n"
		}
		ch.UserMessage <- msg
	}

	ch.Disconect <- c
}
