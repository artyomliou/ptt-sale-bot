PTT sale bot
============

This is a side project with 2 purposes:
1. Notify me that someone post something I want to buy on the PTT forum
2. Improve my programming skills (golang, unix socket, gRPC...)

# Build
Before build, please install [protobuf compiler](https://grpc.io/docs/protoc-installation/)
```shell
make grpc_api/%.pb.go
make
```

# Execute

## Server
```shell
./pttsalebot_server -tg-bot-key "PLEASE_GET_YOUR_OWN_TOKEN" -tg-channel-id "SAME_AS_BOT_KEY"
```

## Client
```shell
./pttsalebot_client -service Crawler -rpc AddTarget -url "https://www.ptt.cc/bbs/DC_SALE/index.html"
./pttsalebot_client -service Filterer -rpc AddInterestTopic -topic-name "Canon RF 16 f2.8" -topic-patterns "(?i)(canon|rf),(?i)16(mm)?,(?i)(f)?2.8"
```
