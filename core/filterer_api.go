package core

import (
	"artyomliou/ptt-sale-bot/grpc_api"
	"context"
	"log"
	"regexp"
	"strings"
)

type FiltererApi struct {
	grpc_api.UnimplementedFiltererServer
	Instance *filterer
}

func (api *FiltererApi) AddInterestTopic(ctx context.Context, req *grpc_api.AddInterestTopicRequest) (*grpc_api.AddInterestTopicReply, error) {
	log.Printf("[Filter API] AddInterestTopic: %s %s", req.Name, strings.Join(req.Patterns, ","))

	// pre-compile regexp
	compiledPatterns := []*regexp.Regexp{}
	for _, pattern := range req.Patterns {
		compiledPattern, err := regexp.Compile(pattern)
		if err != nil {
			log.Printf("[Filter API] failed to compile pattern: %s", err.Error())
			return &grpc_api.AddInterestTopicReply{
				Ok: false,
			}, err
		}
		compiledPatterns = append(compiledPatterns, compiledPattern)
	}

	api.Instance.AddInterestTopic(InterestTopic{
		Name:             req.Name,
		CompiledPatterns: compiledPatterns,
	})
	return &grpc_api.AddInterestTopicReply{
		Ok: true,
	}, nil
}
