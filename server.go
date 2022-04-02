package main

import (
	"context"
	"net"
	"sync"

	broadcast "fake.com/RPC_Chat_App/broadcast"
	errhelp "fake.com/RPC_Chat_App/errhelp"
	"google.golang.org/grpc"
)

type MsgChan struct {
	msg    chan broadcast.Message
	active bool
	index  int
}

type BroadcastServer struct {
	mu           sync.Mutex
	userChannels []MsgChan
	broadcast.UnimplementedBroadcastServer
}

func (b *BroadcastServer) CreateStream(conn *broadcast.Connect, stream broadcast.Broadcast_CreateStreamServer) error {
	currUserChannel := make(chan broadcast.Message)
	currMsgChan := b.addNewMsgChannel(currUserChannel)
	return b.sendMessagesToStream(currUserChannel, stream, currMsgChan)
}

func (b *BroadcastServer) addNewMsgChannel(currUserChannel chan broadcast.Message) MsgChan {
	b.mu.Lock()
	currUserChannelLen := len(b.userChannels)
	msgChan := MsgChan{msg: currUserChannel, index: currUserChannelLen, active: true}
	b.userChannels = append(b.userChannels, msgChan)
	b.mu.Unlock()
	return msgChan
}

func (b *BroadcastServer) sendMessagesToStream(currUserChannel chan broadcast.Message, stream broadcast.Broadcast_CreateStreamServer, currMsgChan MsgChan) error {
	for {
		msg := <-currUserChannel
		err := stream.Send(&msg)
		if errhelp.ErrorExists(err) {
			b.inactivateUserChannel(currMsgChan.index)
			return err
		}
	}
}

func (b *BroadcastServer) inactivateUserChannel(indexOfChan int) {
	b.mu.Lock()
	b.userChannels[indexOfChan].active = false
	b.mu.Unlock()
}

func (b *BroadcastServer) BroadcastMessage(ctx context.Context, message *broadcast.Message) (*broadcast.Close, error) {
	for _, userCh := range b.userChannels {
		if b.channelStillActive(userCh) {
			b.sendMessageToChannel(message, userCh.msg)
		}
	}
	return &broadcast.Close{Msg: message.Msg}, nil
}

func (b *BroadcastServer) channelStillActive(userCh MsgChan) bool {
	return userCh.active
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
