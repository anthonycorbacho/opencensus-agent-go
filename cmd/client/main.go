package main

import (
	"fmt"

	pb "github.com/anthonycorbacho/opencensus-agent-go/pkg/api"
	"go.opencensus.io/plugin/ocgrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	rpc, err := grpc.Dial(
		"localhost:9999",
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),
		grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	c := pb.NewUserServiceClient(rpc)
	r, err := c.Get(context.Background(), &pb.UserRequest{Name: "aa"})
	fmt.Println(r, err)
}
