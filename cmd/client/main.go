package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	pb "github.com/ermyar/subpub-service/pkg/api/subpub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("Creating Client")

	conn, err := grpc.NewClient(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Fault creating new client")
	}
	defer conn.Close()
	log.Println("Created succesfully")

	cli := pb.NewPubSubClient(conn)

	scn := bufio.NewScanner(os.Stdin)

	log.Println("Scanning stdout")

	for scn.Scan() {
		text := scn.Text()
		args := strings.SplitN(text, " ", 3)

		switch args[0] {
		case "publish":
			_, err := cli.Publish(context.Background(), &pb.PublishRequest{
				Key:  args[1],
				Data: args[2],
			})
			if err != nil {
				log.Printf("Publish error: %v\n", err)
			} else {
				log.Println("Published succesfully")
			}

		case "subscribe":
			stream, err := cli.Subscribe(context.Background(), &pb.SubscribeRequest{
				Key: args[1],
			})
			if err != nil {
				log.Fatal(err)
				break
			}

			log.Println("Subscribed succesfully")

			// "demon" that receive our data
			go func() {
				for {

					event, err := stream.Recv()

					if err != nil {
						if errors.Is(err, io.EOF) {
							continue
						}
						log.Println(err)
						return
					}

					data := event.GetData()

					// output of our service
					fmt.Printf("Received:\n%s\n", data)
				}
			}()
		}
	}

}
