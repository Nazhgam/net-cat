package chat

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	u "net-cat/Utilities"
	m "net-cat/models"
)

// creating new chat room
func NewRoom(name string) m.Room {
	room := m.Room{
		Name:           name,
		Path:           name + "_" + time.Now().String()[:19] + ".txt",
		ConnectionChan: make(chan net.Conn),
		UserMessage:    make(chan m.Msg),
		Disconect:      make(chan net.Conn),
		AccountsOfUser: make(map[net.Conn]string),
		UsersList:      make(map[string]bool),
	}

	// save room to list
	u.WriteToArchive(&m.Room{Path: m.ChatlistPath}, room.Path+"\n")

	// create log-file
	file, _ := os.OpenFile(room.Path, os.O_RDWR|os.O_CREATE, 0777)
	file.Write([]byte(""))
	file.Close()
	go startHub(&room)

	return room
}

// choosing chat room on entrance
func SelectRoom(conn net.Conn) *m.Room {
	if len(m.Rooms) == 0 {
		conn.Write([]byte("No active chats. Enter the name of new chat:"))

		for {
			rd := bufio.NewReader(conn)
			chatName, _ := rd.ReadString('\n')
			if chatName == "\n" {
				conn.Write([]byte("Wrong name!\n try another\n"))
				continue
			}
			chatName = chatName[:len(chatName)-1]
			room := NewRoom(chatName)
			m.Rooms = append(m.Rooms, room)
			return &room
		}
	} else {
	INPUT:
		for i, r := range m.Rooms {
			conn.Write([]byte("To choose " + r.Name + " enter " + fmt.Sprint(i) + "\n"))
		}
		conn.Write([]byte("or enter name of new chat\n"))
		for {
			rd := bufio.NewReader(conn)
			input, _ := rd.ReadString('\n')
			if input == "\n" {
				return nil
			}

			input = input[:len(input)-1]

			n := u.Atoi(input)
			if n != -1 { // if input data - is number
				if n < 0 || n >= len(m.Rooms) {
					conn.Write([]byte("Wrong chat!\n try another\n"))
					goto INPUT
				}
				if len(m.Rooms[n].AccountsOfUser) > 9 {
					conn.Write([]byte("This chat is full!\n try another\n"))
					goto INPUT
				}
				return &m.Rooms[n]
			}
			for _, room := range m.Rooms { // check names of existing rooms
				if room.Name == input {
					conn.Write([]byte("Wrong name!\n try another\n"))
					goto INPUT
				}
			}
			room := NewRoom(input)
			m.Rooms = append(m.Rooms, room)
			return &room
		}
	}
}
