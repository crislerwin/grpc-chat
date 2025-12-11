package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/crislerwin/grpc-chat/proto/chat"
	"google.golang.org/grpc"
)

type chatServer struct {
	chat.UnimplementedChatServiceServer
	mu      sync.Mutex
	clients map[chat.ChatService_JoinServer]bool
	message chan *chat.Message
}

func newServer() *chatServer {
	s := &chatServer{
		clients: make(map[chat.ChatService_JoinServer]bool),
		message: make(chan *chat.Message),
	}
	// Start broadcast goroutine
	go s.broadcast()
	return s
}

func (s *chatServer) broadcast() {
	for msg := range s.message {
		s.mu.Lock()
		for client := range s.clients {
			if err := client.Send(msg); err != nil {
				fmt.Printf("Error sending message to client: %v\n", err)
			}
		}
		s.mu.Unlock()
	}
}

func (s *chatServer) Join(stream chat.ChatService_JoinServer) error {
	s.mu.Lock()
	s.clients[stream] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, stream)
		s.mu.Unlock()
	}()

	// Receive messages from this client and broadcast them
	for {
		msg, err := stream.Recv()
		if err != nil {
			fmt.Printf("Client disconnected: %v\n", err)
			return err
		}
		s.message <- msg
	}
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	server := grpc.NewServer()
	chat.RegisterChatServiceServer(server, newServer())
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to init: %v", err)
	}
}
