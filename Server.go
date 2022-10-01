package main

import (
	"log"
	"net"
)

const (
	HOST = "localHost"
	PORT = "3000"
	TYPE = "tcp"
)

type Server struct {
	configuration Configuration
}

func (server *Server) init() {
	tcpServer, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		connection, err := tcpServer.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go server.handleConnection(&connection)
	}
}

func (server *Server) handleConnection(connection *net.Conn) {
}
