package main

import (
	"context"
	"log"
	"net"
	"sync"

	broadcast "fake.com/RPC_Chat_App/broadcast"
	"google.golang.org/grpc"
)

type BroadcastServer struct {
	mu           sync.Mutex
	userChannels []chan broadcast.Message
	broadcast.UnimplementedBroadcastServer
}

func (b *BroadcastServer) CreateStream(conn *broadcast.Connect, stream broadcast.Broadcast_CreateStreamServer) error {
	currUserChannel := make(chan broadcast.Message)
	b.mu.Lock()
	b.userChannels = append(b.userChannels, currUserChannel)
	b.mu.Unlock()
	for {
		msg := <-currUserChannel
		stream.Send(&msg)
	}
	return nil
}

func (b *BroadcastServer) BroadcastMessage(ctx context.Context, message *broadcast.Message) (*broadcast.Close, error) {
	for _, userCh := range b.userChannels {
		userCh <- *message
	}
	return &broadcast.Close{Msg: message.Msg}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("Error when trying to listen at port 5000: %v", err)
	}
	server := grpc.NewServer()
	broadcast.RegisterBroadcastServer(server, &BroadcastServer{})
	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("Error while trying to serve: %v", err)
	}
}
