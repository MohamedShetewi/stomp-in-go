package main

import (
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	HOST = "localHost"
	PORT = "3000"
	TYPE = "tcp"
)

const maxTcpTime = 9000000 * time.Millisecond

type Server struct {
	config      Configuration
	clientList  []*Client
	clientsLock sync.Mutex
}

type Client struct {
	conn        *net.Conn
	outHB       int64
	inHB        int64
	receiveChan chan []byte
	sendChan    chan []byte
}

func (server *Server) init() {
	tcpServer, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		tcpConnection, err := tcpServer.Accept()

		client := &Client{
			conn:        &tcpConnection,
			outHB:       server.config.defaultHB,
			inHB:        server.config.defaultHB,
			receiveChan: make(chan []byte),
			sendChan:    make(chan []byte),
		}
		if err != nil {
			log.Println(err)
			return
		}
		err = tcpConnection.SetDeadline(server.config.maxTcpConnTime)
		if err != nil {
			log.Println(err)
			return
		}
		go server.clientListener(client)
		go server.clientSendReceive(client)
	}
}

/*
This is meant to be a goroutine that listens to the frames sent from the client.
This goroutine will stay alive as long as the connection is up.
*/
func (server *Server) clientListener(client *Client) {
	for {
		buffer := make([]byte, server.config.maxFrameSize)
		_, err := (*client.conn).Read(buffer)

		if err != nil {
			log.Println(err)
			return
		}
		client.receiveChan <- buffer
	}
}

/*
This func is meant to be a goroutine that manages the sending + receiving
from the client. This also manages the Heartbeat between the server and the client.
*/
func (server *Server) clientSendReceive(client *Client) {
	outBeatDeadline := time.After(maxTcpTime) // w.r.t. server
	inBeatDeadline := time.After(maxTcpTime)  // w.r.t. server
	for {
		select {
		case receivedMsg := <-client.receiveChan:
			frame, errCode := Decode(string(receivedMsg))
			if errCode != OK {
				//TODO send error msg to the client
			}
			commandHandler := CommandHandlerMap[frame.command]
			commandHandler(server, client, frame)
			newInBeat := time.Duration(client.outHB+rand.Int63n(60000)+30000) * time.Millisecond
			inBeatDeadline = time.After(newInBeat)
		case msg := <-client.sendChan:
			_, err := (*client.conn).Write(msg)

			if err != nil {
				continue
			}
			newOutBeat := time.Duration(client.inHB+rand.Int63n(60000)+30000) * time.Millisecond
			outBeatDeadline = time.After(newOutBeat)
		case <-outBeatDeadline:
			if client.inHB != -1 {
				// TODO send HeartBeat packet
			}
		case <-inBeatDeadline:
			if client.outHB != -1 {
				// TODO send error msg + close the connection
			}
		}
	}
}

func (server *Server) hasClient(connection *Client) (bool, int) {
	for idx, conn := range server.clientList {
		if conn == connection {
			return true, idx
		}
	}
	return false, -1
}

func (server *Server) removeConnection(client *Client) {
	if ok, idx := server.hasClient(client); ok {
		server.clientsLock.Lock()
		defer server.clientsLock.Unlock()
		removeIndex(server.clientList, idx)
		err := (*client.conn).Close()
		if err != nil {
			log.Println(err)
			return
		}
	}
}
