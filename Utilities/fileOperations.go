package utility

import (
	m "net-cat/models"
	"os"
)

// loging into the file
func WriteToArchive(room *m.Room, s string) {

	file, _ := os.OpenFile(room.Path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)

	file.WriteString(s)
	file.Close()

}

// gets data from file
func FileReader(path string) string {
	archive, err := os.Open(path)
	if err != nil {
		return ""
	}
	stat, _ := archive.Stat()
	size := stat.Size()
	lengs := int(size)
	arr := make([]byte, lengs)
	archive.Read(arr)
	st := string(arr)
	return st
}
