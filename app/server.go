package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	//Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := l.Accept()

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	var messages []string
	scanner := bufio.NewScanner(conn)

	//scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		message := scanner.Text()
		messages = append(messages, message)
		fmt.Println(scanner.Text())
		if message == "ping" {
			fmt.Printf("cmd  %s received\n", message)
			_, err = conn.Write([]byte("+PONG\r\n"))
			if err != nil {
				return
			}
		}
		if message == "" {
			break
		}
	}
	//_, err = conn.Write([]byte("+PONG\r\n"))
	//if err != nil {
	//	return
	//}
	err = conn.Close()
	if err != nil {
		return
	}

}
