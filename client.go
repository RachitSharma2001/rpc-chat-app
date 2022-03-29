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
	errhelp "fake.com/RPC_Chat_App/errhelp"
	"google.golang.org/grpc"
)

var userInputScanner = bufio.NewScanner(os.Stdin)

func main() {
	nameOfUser := getNameOfUser()
	conn := createConnection()
	defer conn.Close()
	client := createClient(conn)
	streamOfSentMsgs := createStreamFromClient(client, nameOfUser)
	listenForSentMsgs(streamOfSentMsgs)

	tellServerUserHasEntered(client, nameOfUser)
	setupUserExitMessage(client, nameOfUser)
	listenForUserInput(client, nameOfUser)
}

func getNameOfUser() string {
	fmt.Print("Please enter your name: ")
	return getInputFromUser()
}

func createConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(":5000", grpc.WithInsecure())
	if errhelp.ErrorExists(err) {
		log.Fatalf("Error while dialing port 5000: %v", err)
	}
	return conn
}

func createClient(conn *grpc.ClientConn) broadcast.BroadcastClient {
	return broadcast.NewBroadcastClient(conn)
}

func createStreamFromClient(client broadcast.BroadcastClient, nameOfUser string) broadcast.Broadcast_CreateStreamClient {
	streamOfSentMsgs, err := client.CreateStream(context.Background(), &broadcast.Connect{Uid: nameOfUser})
	if errhelp.ErrorExists(err) {
		errhelp.ThrowCreateStreamErr(err)
	}
	return streamOfSentMsgs
}

func listenForSentMsgs(streamOfSentMsgs broadcast.Broadcast_CreateStreamClient) {
	go func() {
		for {
			newMessage, err := streamOfSentMsgs.Recv()
			if errhelp.ErrorExists(err) {
				log.Fatalf("Error while recieving messages: %v", err)
			}
			fmt.Printf("%s: %s\n", newMessage.Sender, newMessage.Msg)
		}
	}()
}
func setupUserExitMessage(client broadcast.BroadcastClient, nameOfUser string) {
	exitedChatListener := make(chan os.Signal)
	signal.Notify(exitedChatListener, os.Interrupt, syscall.SIGTERM)
	go listenForUserExit(client, exitedChatListener, nameOfUser)
}

func listenForUserExit(client broadcast.BroadcastClient, exitedChatListener chan os.Signal, nameOfUser string) {
	<-exitedChatListener
	tellServerUserHasExited(client, nameOfUser)
	os.Exit(1)
}

func tellServerUserHasEntered(client broadcast.BroadcastClient, nameOfUser string) {
	enteredChatMsg := getUserEnteredChatMsg(nameOfUser)
	client.BroadcastMessage(context.Background(), &broadcast.Message{Sender: nameOfUser, Msg: enteredChatMsg})
}

func tellServerUserHasExited(client broadcast.BroadcastClient, nameOfUser string) {
	exitedChatMsg := getUserExitedChatMsg(nameOfUser)
	client.BroadcastMessage(context.Background(), &broadcast.Message{Sender: nameOfUser, Msg: exitedChatMsg})
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

func listenForUserInput(client broadcast.BroadcastClient, nameOfUser string) {
	for {
		if userInputScanner.Scan() {
			msgContent := userInputScanner.Text()
			newMessage := &broadcast.Message{Sender: nameOfUser, Msg: msgContent}
			client.BroadcastMessage(context.Background(), newMessage)
		}
	}
}
