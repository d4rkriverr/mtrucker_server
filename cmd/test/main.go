package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
)

func main() {
	const (
		HOST = "localhost"
		PORT = "1885"
		TYPE = "tcp"
	)
	//
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	var msg string
	for {
		fmt.Scanln(&msg)
		str, _ := hex.DecodeString(msg)
		_, err = conn.Write(str)
		if err != nil {
			println("Write data failed:", err.Error())
			os.Exit(1)
		}
	}
	// msg := "78780D01012345678901234500018CDD0D0A"

}
