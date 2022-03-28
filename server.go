package main

import (
	"context"
	"net"
	"sync"

	broadcast "fake.com/RPC_Chat_App/broadcast"
	errhelp "fake.com/RPC_Chat_App/errhelp"
	"google.golang.org/grpc"
)

type BroadcastServer struct {
	mu           sync.Mutex
	userChannels []chan broadcast.Message
	broadcast.UnimplementedBroadcastServer
}

func (b *BroadcastServer) CreateStream(conn *broadcast.Connect, stream broadcast.Broadcast_CreateStreamServer) error {
	currUserChannel := make(chan broadcast.Message)
	b.updateUserChannelsSafely(currUserChannel)
	b.sendMessagesToStream(currUserChannel, stream)
	return nil
}

func (b *BroadcastServer) updateUserChannelsSafely(currUserChannel chan broadcast.Message) {
	b.mu.Lock()
	b.userChannels = append(b.userChannels, currUserChannel)
	b.mu.Unlock()
}

func (b *BroadcastServer) sendMessagesToStream(currUserChannel chan broadcast.Message, stream broadcast.Broadcast_CreateStreamServer) {
	for {
		msg := <-currUserChannel
		stream.Send(&msg)
	}
}

func (b *BroadcastServer) BroadcastMessage(ctx context.Context, message *broadcast.Message) (*broadcast.Close, error) {
	for _, userCh := range b.userChannels {
		b.sendMessageToChannel(message, userCh)
	}
	return &broadcast.Close{Msg: message.Msg}, nil
}

func (b *BroadcastServer) sendMessageToChannel(message *broadcast.Message, userCh chan broadcast.Message) {
	userCh <- *message
}

func main() {
	listener := listenAtPort(":5000")
	server := registerServer()
	connectServerToListener(listener, server)
}

func listenAtPort(port string) net.Listener {
	lis, err := net.Listen("tcp", port)
	if errhelp.ErrorExists(err) {
		errhelp.ThrowPortListenErr(err)
	}
	return lis
}

func registerServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	broadcast.RegisterBroadcastServer(grpcServer, &BroadcastServer{})
	return grpcServer
}

func connectServerToListener(listener net.Listener, server *grpc.Server) {
	err := server.Serve(listener)
	if errhelp.ErrorExists(err) {
		errhelp.ThrowServeErr(err)
	}
}
