package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	pb "github.com/ermyar/subpub-service/pkg/api/subpub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	conf configJSON
	// there maybe other param
}

type configJSON struct {
	Port int `json:"port"`
}

func (c *client) readJSON(path string) (err error) {
	var jsonFile []byte
	jsonFile, err = os.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonFile, &c.conf)
	return err
}

func (c *client) init(path string) error {
	return c.readJSON(path)

}

func main() {
	cc := client{}
	path := flag.String("config", "configs/config.json", "path to .json which configurate our client")
	flag.Parse()

	if err := cc.init(*path); err != nil {
		log.Fatal(err)
	}

	log.Println("Creating Client")

	conn, err := grpc.NewClient(fmt.Sprintf(":%d", cc.conf.Port),
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
						log.Fatalln(err)
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
