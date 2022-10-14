package server

import (
	"github.com/MohamedShetewi/stomp-in-go/frame"
	"github.com/MohamedShetewi/stomp-in-go/utils"
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
			continue
		}
		err = tcpConnection.SetDeadline(server.config.maxTcpConnTime)
		if err != nil {
			log.Println(err)
			continue
		}
		server.AddClient(client)
		go server.tcpListener(client)
		go server.connectListener(client)
	}
}

/*
This is meant to be a goroutine that listens to the frames sent from the client.
This goroutine will stay alive as long as the connection is up.
*/
func (server *Server) tcpListener(client *Client) {
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

func (server *Server) connectListener(client *Client) {
	initDeadline := time.After(server.config.deadlineForConnect)
	select { // This is waiting for the connect-frame
	case receivedMsg := <-client.receiveChan:
		frm, err := frame.Decode(string(receivedMsg))
		if err != nil {
			sendError(client, make(map[string]string), "Decode Error: "+err.Error())
			server.RemoveClient(client)
			return
		}
		go server.receiveFromClient(client)
		go server.sendToClient(client)
		commandHandler := commandHandlerMap[frm.Command]
		commandHandler(server, client, frm)
	case <-initDeadline:
		sendError(client, make(map[string]string), "Timeout")
		return
	}
}

/*
This func is meant to be a goroutine that manages the sending + receiving
from the client. This also manages the Heartbeat between the server and the client.
*/

func (server *Server) receiveFromClient(client *Client) {
	deadline := newDeadline(client.outHB)
	for {
		select {
		case msg := <-client.receiveChan:
			frm, err := frame.Decode(string(msg))
			if err != nil {
				sendError(client, make(map[string]string), "Decode Error: "+err.Error())
				server.RemoveClient(client)
				return
			}
			commandHandler := commandHandlerMap[frm.Command]
			commandHandler(server, client, frm)
			deadline = newDeadline(client.outHB)
		case <-deadline:
			if client.outHB != -1 {
				sendError(client, make(map[string]string),
					"Timeout")
				server.RemoveClient(client)
				return
			}
		}
	}
}

func (server *Server) sendToClient(client *Client) {
	deadline := newDeadline(client.inHB)
	for {
		select {
		case msg := <-client.sendChan:
			_, err := (*client.conn).Write(msg)
			if err != nil {
				log.Println(err)
				return
			}
			deadline = newDeadline(client.inHB)
		case <-deadline:
			if client.inHB != -1 {
				// TODO send HeartBeat packet
			}
		}
	}
}

func (server *Server) HasClient(connection *Client) (bool, int) {
	for idx, conn := range server.clientList {
		if conn == connection {
			return true, idx
		}
	}
	return false, -1
}

func (server *Server) RemoveClient(client *Client) {
	if ok, idx := server.HasClient(client); ok {
		server.clientsLock.Lock()
		defer server.clientsLock.Unlock()
		err := utils.RemoveIndex(server.clientList, idx)
		if err != nil {
			return
		}
		err = (*client.conn).Close()
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (server *Server) AddClient(client *Client) {
	if ok, _ := server.HasClient(client); !ok {
		server.clientsLock.Lock()
		defer server.clientsLock.Unlock()
		server.clientList = append(server.clientList, client)
	}
}

func newDeadline(hb int64) <-chan time.Time {
	deadline := time.Duration(hb+rand.Int63n(60000)+30000) * time.Millisecond
	return time.After(deadline)
}
