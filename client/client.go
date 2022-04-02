package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	broadcast "fake.com/RPC_Chat_App/broadcast"
	errhelp "fake.com/RPC_Chat_App/errhelp"
	"google.golang.org/grpc"
)

const port = ":5000"

var userInputScanner *bufio.Scanner
var client broadcast.BroadcastClient
var nameOfUser string
var wg sync.WaitGroup

func main() {
	nameOfUser = getNameOfUser()
	conn := createConnection()
	defer conn.Close()
	createServerStreamListener(conn)
	createUserListener()
	wg.Wait()
}

func getNameOfUser() string {
	userInputScanner = bufio.NewScanner(os.Stdin)
	fmt.Print("Please enter your name: ")
	return getInputFromUser()
}

func getInputFromUser() string {
	if userInputScanner.Scan() {
		return userInputScanner.Text()
	}
	return ""
}

func createConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if errhelp.ErrorExists(err) {
		errhelp.ThrowDialErr(err)
	}
	return conn
}

func createServerStreamListener(conn *grpc.ClientConn) {
	client = createClient(conn)
	streamOfSentMsgs := createStreamFromClient()
	listenForSentMsgs(streamOfSentMsgs)
}

func createClient(conn *grpc.ClientConn) broadcast.BroadcastClient {
	return broadcast.NewBroadcastClient(conn)
}

func createStreamFromClient() broadcast.Broadcast_CreateStreamClient {
	streamOfSentMsgs, err := client.CreateStream(context.Background(), &broadcast.Connect{Uid: nameOfUser})
	if errhelp.ErrorExists(err) {
		errhelp.ThrowCreateStreamErr(err)
	}
	return streamOfSentMsgs
}

func listenForSentMsgs(streamOfSentMsgs broadcast.Broadcast_CreateStreamClient) {
	wg.Add(1)
	go func() {
		for {
			newMessage, err := streamOfSentMsgs.Recv()
			if errhelp.ErrorExists(err) {
				errhelp.ThrowReceiveMsgErr(err)
				break
			}
			printMessage(newMessage.Sender, newMessage.Msg)
		}
		wg.Done()
	}()
}

func printMessage(sender, msg string) {
	fmt.Printf("%s: %s\n", sender, msg)
}

func createUserListener() {
	tellServerUserHasEntered()
	setupUserExitMessage()
	listenForUserInput()
}

func setupUserExitMessage() {
	exitedChatListener := make(chan os.Signal)
	signal.Notify(exitedChatListener, os.Interrupt, syscall.SIGTERM)
	go listenForUserExit(exitedChatListener)
}

func listenForUserExit(exitedChatListener chan os.Signal) {
	<-exitedChatListener
	tellServerUserHasExited()
	os.Exit(1)
}

func tellServerUserHasEntered() {
	enteredChatMsg := getUserEnteredChatMsg(nameOfUser)
	client.BroadcastMessage(context.Background(), &broadcast.Message{Sender: nameOfUser, Msg: enteredChatMsg})
}

func tellServerUserHasExited() {
	exitedChatMsg := getUserExitedChatMsg(nameOfUser)
	client.BroadcastMessage(context.Background(), &broadcast.Message{Sender: nameOfUser, Msg: exitedChatMsg})
}

func getUserEnteredChatMsg(name string) string {
	return fmt.Sprintf("%s has entered the chat", name)
}

func getUserExitedChatMsg(name string) string {
	return fmt.Sprintf("%s has exited the chat", name)
}

func listenForUserInput() {
	wg.Add(1)
	for {
		if userInputScanner.Scan() {
			msgContent := userInputScanner.Text()
			newMessage := &broadcast.Message{Sender: nameOfUser, Msg: msgContent}
			_, err := client.BroadcastMessage(context.Background(), newMessage)
			if errhelp.ErrorExists(err) {
				break
			}
		}
	}
	wg.Done()
}
