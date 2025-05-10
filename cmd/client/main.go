package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	utils "github.com/ermyar/subpub-service/cmd/utils"
	pb "github.com/ermyar/subpub-service/pkg/api/subpub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config := utils.ConfigJSON{}

	path := flag.String("config", "configs/config.json", "path to .json which configurate our client")
	flag.Parse()

	if err := config.Init(*path); err != nil {
		log.Fatal(err)
	}

	log.Println("Creating Client on", config.HostName, config.Port)

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", config.HostName, config.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

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
			if len(args) != 3 {
				log.Println("wrong args: must be 3")
				continue
			}
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
			if len(args) != 2 {
				log.Println("wrong args: needs only 2 args")
				continue
			}
			stream, err := cli.Subscribe(context.Background(), &pb.SubscribeRequest{
				Key: args[1],
			})
			if err != nil {
				log.Fatal(err)
				break
			}

			log.Println("Subscribed succesfully")

			// "demon" that receive our data
			go func(subscription string) {
				for {
					event, err := stream.Recv()

					if err != nil {
						log.Fatalln(err)
						return
					}

					data := event.GetData()

					// output of our service
					fmt.Printf("Received from %s:\n%s\n", subscription, data)
				}
			}(args[1])
		default:
			log.Println("wrong args")
		}
	}

}
