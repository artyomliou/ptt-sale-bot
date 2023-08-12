package main

import (
	"artyomliou/ptt-sale-bot/grpc_api"
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	service            string
	rpc                string
	url                string
	topicName          string
	topicPatternString string
	topicPatterns      []string
)

func main() {
	flag.StringVar(&service, "service", "", "Service")
	flag.StringVar(&rpc, "rpc", "", "The command being sent to api server")
	flag.StringVar(&url, "url", "", "Webpage URL")
	flag.StringVar(&topicName, "topic-name", "", "Topic name")
	flag.StringVar(&topicPatternString, "topic-patterns", "", "Topic patterns, separated by comma.")
	topicPatterns = strings.Split(topicPatternString, ",")
	flag.Parse()

	// connect to grpc server over socket
	address := fmt.Sprintf("unix://%s", grpc_api.SocketPath)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Cannot connect to %s: %s", address, err.Error())
	}

	// call API
	if service == "Crawler" && rpc == "AddTarget" {
		client := grpc_api.NewCrawlerClient(conn)
		reply, err := client.AddTarget(context.TODO(), &grpc_api.AddTargetRequest{
			Url: url,
		})
		if err != nil {
			log.Fatalf("Server reply with error: %s", err.Error())
		}
		handleServerReply(reply.Ok)
	} else if service == "Filterer" && rpc == "AddInterestTopic" {
		client := grpc_api.NewFiltererClient(conn)
		reply, err := client.AddInterestTopic(context.TODO(), &grpc_api.AddInterestTopicRequest{
			Name:     topicName,
			Patterns: topicPatterns,
		})
		if err != nil {
			log.Fatalf("Server reply with error: %s", err.Error())
		}
		handleServerReply(reply.Ok)
	} else {
		log.Fatalln("Unknown service and command")
	}
}

func handleServerReply(ok bool) {
	if ok {
		log.Println("OK.")
	} else {
		log.Panicln("Failed.")
	}
}
