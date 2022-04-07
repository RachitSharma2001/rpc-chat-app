package testing_test

import (
	"context"
	"testing"
	"time"

	"fake.com/RPC_Chat_App/broadcast"
	"fake.com/RPC_Chat_App/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

func TestTesting(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testing Suite")
}

var _ = Describe("Server and Client", func() {
	go server.SetUpServer()

	const port = ":5000"
	conn, _ := grpc.Dial(port, grpc.WithInsecure())
	client := broadcast.NewBroadcastClient(conn)
	var stream1 broadcast.Broadcast_CreateStreamClient
	var err error
	firstUserName := "Rich"
	Context("When Client creates a stream", func() {
		stream1, err = client.CreateStream(context.Background(), &broadcast.Connect{Uid: firstUserName})
		It("No error should be returned", func() {
			Expect(err).ToNot(HaveOccurred())
		})
	})

	client2 := broadcast.NewBroadcastClient(conn)
	var stream2 broadcast.Broadcast_CreateStreamClient
	secondUserName := "Mike"
	Context("When another client enters the chat", func() {
		stream2, err = client2.CreateStream(context.Background(), &broadcast.Connect{Uid: secondUserName})
		It("No error should be returned", func() {
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("When the second client sends a message", func() {
		someMsg := "Hi"
		client2.BroadcastMessage(context.Background(), &broadcast.Message{Sender: secondUserName, Msg: someMsg})
		var msgRecievedByFirstClient string
		var msgRecievedBySecondClient string
		go func() {
			returnedMsgFromStream, _ := stream1.Recv()
			msgRecievedByFirstClient = returnedMsgFromStream.Msg
		}()
		go func() {
			returnedMsgFromStream, _ := stream2.Recv()
			msgRecievedBySecondClient = returnedMsgFromStream.Msg
		}()
		time.Sleep(1 * time.Second)
		It("The first client should recieve it", func() {
			Expect(msgRecievedByFirstClient).To(Equal(someMsg))
		})
		It("The second client should recieve it", func() {
			Expect(msgRecievedBySecondClient).To(Equal(someMsg))
		})
	})

})
