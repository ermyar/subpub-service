package main

import (
	// status "google. "
	"context"
	"log"
	"net"
	"sync"

	subpub "github.com/ermyar/subpub-service/cmd/subpub"
	pb "github.com/ermyar/subpub-service/pkg/api/subpub"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedPubSubServer

	sp subpub.SubPub
}

func NewServer() *server {
	return &server{sp: subpub.NewSubPub()}
}

func (s *server) Publish(_ context.Context, req *pb.PublishRequest) (*emptypb.Empty, error) {
	key, data := req.GetKey(), req.GetData()

	log.Printf("received:\nkey: %v\ndata: %v\n", key, data)

	s.sp.Publish(key, data)

	return &emptypb.Empty{}, nil
}

func (s *server) Subscribe(req *pb.SubscribeRequest, stream pb.PubSub_SubscribeServer) error {
	key := req.GetKey()
	log.Println("Subscribe request!")

	wg := sync.WaitGroup{}
	wg.Add(1)

	_, err := s.sp.Subscribe(key, func(msg any) {
		str, ok := (msg).(string)
		if !ok {
			log.Println("msg not a string")
			return
		}
		if err := stream.Send(&pb.Event{Data: str}); err != nil {
			log.Println("MsgHandler send fault: ", err)
		}
		// wg.Done()
		// we have no Unscribtion in out service
		// if we have, we should exit from this func
		// (then stream.ctx will closed and everything would fine)
	})

	log.Println("Subscribed with: ", err)

	wg.Wait()
	// have to hold this work, if we free, we cant send data..

	return err
}

func main() {
	log.Println("Trying to connect")
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return
	}
	log.Println("Connected")

	s := grpc.NewServer()
	pb.RegisterPubSubServer(s, NewServer())

	log.Println("Server started!")

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	s.GracefulStop()

	log.Println("Server finished!")

}
