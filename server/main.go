package main

import (
	"artyomliou/ptt-sale-bot/core"
	"artyomliou/ptt-sale-bot/grpc_api"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigs
		log.Println(sig)
		cancelFunc()
	}()

	crawler := core.NewCrawler(ctx)
	crawler.Run()

	parser := core.NewParser(ctx)
	parser.SetInput(crawler.Output)
	parser.Run()

	filterer := core.NewFilterer(ctx)
	filterer.SetInput(parser.Output)
	filterer.Run()

	notifier := core.NewNotifier(ctx)
	notifier.SetInput(filterer.Output)
	notifier.Run()

	crawlerApi := core.CrawlerApi{
		Instance: crawler,
	}
	filtererApi := core.FiltererApi{
		Instance: filterer,
	}

	RunApiServer(ctx, &crawlerApi, &filtererApi)
}

func RunApiServer(ctx context.Context, crawlerApi *core.CrawlerApi, filtererApi *core.FiltererApi) {
	grpcServer := grpc.NewServer()
	grpcServer.RegisterService(&grpc_api.Crawler_ServiceDesc, crawlerApi)
	grpcServer.RegisterService(&grpc_api.Filterer_ServiceDesc, filtererApi)

	var lc net.ListenConfig
	socket, err := lc.Listen(ctx, "unix", grpc_api.SocketPath)
	if err != nil {
		panic(err)
	}

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
		log.Println("[Server] stopped")
	}()

	if err := grpcServer.Serve(socket); err != nil {
		log.Printf("[Server] serve error: %s", err.Error())
	}
}
