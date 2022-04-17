protoc --go_out=plugins=grpc:./passport ./proto/*.proto
protoc --grpc-gateway_out=logtostderr=true:./passport ./proto/*.proto