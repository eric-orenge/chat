package main

import (
	"fmt"
	"log"
	"net/http"

	"gitlab.com/Orenge/chat/trace"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize, CheckOrigin: func(r *http.Request) bool {
	return true
}}

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan *message
	//channel for clients wishing to joining
	join chan *client
	//channel for clients wishing to leave
	leave chan *client
	// clients holds all current clients in this room.
	clients map[*client]bool
	// tracer will receive trace information of activity
	// in the room.
	tracer trace.Tracer

	// avatar is how avatar information will be obtained.
	avatar Avatar
}

func newRoom(avatar Avatar) *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
		// avatar:  avatar,
	}
}

func (rm *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
	if !websocket.IsWebSocketUpgrade(req) {
		fmt.Println("Is not web socket upgrade")
		return
	} else {

		socket, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Fatal("ServeHTTP:", err)
		}
		authCookie, err := req.Cookie("auth")

		if err != nil {
			log.Fatal("Failed to get auth cookie:", err)
		}
		//new client
		client := &client{
			socket:   socket,
			send:     make(chan *message, messageBufferSize),
			room:     rm,
			userData: objx.MustFromBase64(authCookie.Value),
		}
		//allow him to join room
		rm.join <- client
		defer func() {
			rm.leave <- client // leave after executing everything else
		}()
		go client.write() //concurrently write the messages
		client.read()
	}

}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("New client left")

		case msg := <-r.forward:
			r.tracer.Trace("Message received: ", msg.Message)
			for client := range r.clients {
				select {
				case client.send <- msg:
					r.tracer.Trace("--Sent to client")

				//send the messages
				default:
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace("--not sent to client. Client cleaned up")

				}
			}
		}
	}
}
