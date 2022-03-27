package main

import (
	"context"
	"fmt"
	"log"

	broadcast "fake.com/RPC_Chat_App/broadcast"
	"google.golang.org/grpc"
)

func main() {
	// for now username is hardcoded, eventually have user enter this
	username := "Rachit"
	conn, err := grpc.Dial(":5000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error while dialing port 5000: %v", err)
	}
	defer conn.Close()
	client := broadcast.NewBroadcastClient(conn)

	streamOfMsgs, err := client.CreateStream(context.Background(), &broadcast.Connect{Uid: username})
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}
	go func() {
		for {
			newMessage, err := streamOfMsgs.Recv()
			if err != nil {
				log.Fatalf("Error while recieving messages: %v", err)
			}
			fmt.Printf("%s: %s\n", newMessage.Sender, newMessage.Msg)
		}
	}()

	for {
		var msgContent string
		fmt.Scanln(&msgContent)
		newMessage := &broadcast.Message{Sender: username, Msg: msgContent}
		client.BroadcastMessage(context.Background(), newMessage)
	}
}
