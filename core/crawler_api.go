package core

import (
	"artyomliou/ptt-sale-bot/grpc_api"
	"context"
	"log"
)

type CrawlerApi struct {
	grpc_api.UnimplementedCrawlerServer
	Instance *crawler
}

func (api *CrawlerApi) AddTarget(ctx context.Context, req *grpc_api.AddTargetRequest) (*grpc_api.AddTargetReply, error) {
	log.Printf("[Crawler API] AddTarget: %s", req.Url)
	api.Instance.AddTarget(req.Url)
	return &grpc_api.AddTargetReply{
		Ok: true,
	}, nil
}
