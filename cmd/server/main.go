package main

import (
	// status "google. "
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	subpub "github.com/ermyar/subpub-service/cmd/subpub"
	utils "github.com/ermyar/subpub-service/cmd/utils"
	pb "github.com/ermyar/subpub-service/pkg/api/subpub"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedPubSubServer

	sp subpub.SubPub
}

var (
	srConf serverConfig
)

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

	_, err := s.sp.Subscribe(key, func(msg any) {
		str, ok := (msg).(string)
		if !ok {
			log.Println("msg not a string")
			return
		}
		if err := stream.Send(&pb.Event{Data: str}); err != nil {
			log.Println("MsgHandler send fault: ", err)
		}
		// we have no Unscribtion in out service
		// if we have, we should exit from this func
		// (then stream.ctx will closed and everything would fine)
	})

	if err == nil {
		log.Println("Subscribed succesfully!")
	} else {
		log.Println("Subscribed with:", err)
	}

	<-srConf.done
	// have to hold this work, if we free, we cant send data later..

	return err
}

type serverConfig struct {
	done chan struct{}
	conf utils.ConfigJSON
}

func main() {
	srConf := serverConfig{done: make(chan struct{})}

	path := flag.String("config", "configs/config.json", "path to .json which configurate our server")
	flag.Parse()

	if err := srConf.conf.Init(*path); err != nil {
		log.Fatal("reading config error ", err)
	}

	log.Println("Trying to listen", srConf.conf.HostName, srConf.conf.Port, " by ", srConf.conf.Network)
	lis, err := net.Listen(srConf.conf.Network, fmt.Sprintf("%s:%d", srConf.conf.HostName, srConf.conf.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return
	}
	log.Println("Connected")

	// Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	s := grpc.NewServer()
	pb.RegisterPubSubServer(s, NewServer())

	go func() {

		<-ctx.Done()

		s.GracefulStop()
		close(srConf.done)

		s.Stop()

		log.Println("Server finished!")
	}()

	log.Println("Server started!")

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
