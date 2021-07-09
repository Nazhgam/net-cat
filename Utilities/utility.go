package utility

import (
	"fmt"
	"log"
	"os"
	"time"
)

// panic recover
func RecoverFunction() {
	if r := recover(); r != nil {
		fmt.Println("Recovered after Panic ", r)
		file, _ := os.OpenFile("errlog_"+time.Now().Format("2006-01-02")+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

		log.SetOutput(file)

		log.Println(r)
	}
}

// own implementation of convertation from text to number (strconv is forbidden)
func Atoi(s string) int {
	slice := []rune(s)
	nbr := 0
	ln := 0
	for _, char := range slice {
		if char == '-' || char == '+' {
			return -1
		}
		if char < '0' || char > '9' {
			return -1
		} else {
			c := 0
			for a := '0'; a < char; a++ {
				c++
			}
			if ln == 18 {
				if nbr > 922337203685477580 {
					return -1
				}
				if nbr == 922337203685477580 {
					if c > 7 {
						return -1
					}
				}
			}
			if ln == 19 {
				return -1
			}
			nbr = nbr*10 + c
			ln++
		}
	}
	return nbr
}
