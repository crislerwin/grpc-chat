package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/crislerwin/grpc-chat/proto/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := chat.NewChatServiceClient(conn)

	// Get username
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	// Start bidirectional stream
	stream, err := client.Join(context.Background())
	if err != nil {
		log.Fatalf("Failed to join chat: %v", err)
	}

	// Goroutine to receive messages from server
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				return
			}
			fmt.Printf("\n[%s]: %s\n> ", msg.User, msg.Text)
		}
	}()

	// Main goroutine to send messages
	fmt.Println("Connected to chat! Type your messages below:")
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "" {
			continue
		}

		if text == "/quit" || text == "/exit" {
			fmt.Println("Leaving chat...")
			return
		}

		msg := &chat.Message{
			User:      username,
			Text:      text,
			Timestamp: time.Now().Unix(),
		}

		if err := stream.Send(msg); err != nil {
			log.Printf("Error sending message: %v", err)
			return
		}
	}
}
