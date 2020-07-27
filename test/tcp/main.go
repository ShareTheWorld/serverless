package main

import (
	"fmt"
	"net"
)

func main() {
	l, _ := net.Listen("tcp", "0.0.0.0:8888")
	defer l.Close()
	for {
		fmt.Println("wait connection ... ")
		conn, _ := l.Accept()
		fmt.Println("connection success ")

		fmt.Printf("%v,%v", conn, conn.RemoteAddr().String())
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		fmt.Printf("%v,%v", conn, conn.RemoteAddr().String())
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		fmt.Print(string(buf[:n]))
	}
}
