package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	broadcast "fake.com/RPC_Chat_App/broadcast"
	"google.golang.org/grpc"
)

var userInputScanner = bufio.NewScanner(os.Stdin)

func main() {
	nameOfUser := getNameOfUser()

	conn, err := grpc.Dial(":5000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error while dialing port 5000: %v", err)
	}
	defer conn.Close()
	client := broadcast.NewBroadcastClient(conn)

	streamOfMsgs, err := client.CreateStream(context.Background(), &broadcast.Connect{Uid: nameOfUser})
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

	enteredChatMsg := getUserEnteredChatMsg(nameOfUser)
	client.BroadcastMessage(context.Background(), &broadcast.Message{Sender: nameOfUser, Msg: enteredChatMsg})

	exitedChatListener := make(chan os.Signal)
	signal.Notify(exitedChatListener, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-exitedChatListener
		exitedChatMsg := getUserExitedChatMsg(nameOfUser)
		client.BroadcastMessage(context.Background(), &broadcast.Message{Sender: nameOfUser, Msg: exitedChatMsg})
		os.Exit(1)
	}()

	for {
		if userInputScanner.Scan() {
			msgContent := userInputScanner.Text()
			newMessage := &broadcast.Message{Sender: nameOfUser, Msg: msgContent}
			client.BroadcastMessage(context.Background(), newMessage)
		}
	}
}

func getNameOfUser() string {
	fmt.Print("Please enter your name: ")
	return getInputFromUser()
}

func getInputFromUser() string {
	if userInputScanner.Scan() {
		return userInputScanner.Text()
	}
	return ""
}

func getUserEnteredChatMsg(name string) string {
	return fmt.Sprintf("%s has entered the chat", name)
}

func getUserExitedChatMsg(name string) string {
	return fmt.Sprintf("%s has exited the chat", name)
}
