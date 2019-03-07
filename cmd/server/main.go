package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"contrib.go.opencensus.io/exporter/ocagent"
	pb "github.com/anthonycorbacho/opencensus-agent-go/pkg/api"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
)

func main() {
	oce, err := ocagent.NewExporter(
		ocagent.WithInsecure(),
		//ocagent.WithAddress("127.0.0.1:55678"),
		ocagent.WithServiceName(fmt.Sprintf("test-go-app-%d", os.Getpid())))
	if err != nil {
		panic(err)
	}
	trace.RegisterExporter(oce)
	view.RegisterExporter(oce)

	//initJaegerTracing()

	// Some configurations to get observability signals out.
	view.SetReportingPeriod(7 * time.Second)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})

	_, span := trace.StartSpan(context.Background(), "Foo")
	fmt.Println(span.String())
	span.End()

	srv := grpc.NewServer(
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_recovery.StreamServerInterceptor(),
			)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
		)),
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
	)

	pb.RegisterUserServiceServer(srv, &userService{})

	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		panic(err)
	}
	if err := srv.Serve(l); err != nil {
		panic(err)
	}
}

type userService struct{}

func (us *userService) Get(ctx context.Context, ur *pb.UserRequest) (*pb.UserResponse, error) {
	return &pb.UserResponse{
		Id:   "the ID",
		Name: ur.Name,
	}, nil
}

func initJaegerTracing() error {
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: fmt.Sprintf("http://%s%s", "127.0.0.1", ":14268/api/traces"),
		AgentEndpoint:     fmt.Sprintf("%s%s", "127.0.0.1", ":6831"),
		Process: jaeger.Process{
			ServiceName: "test", //os.Getenv("SERVICE_NAME"),
			Tags:        []jaeger.Tag{
				//jaeger.StringTag("pod-id", os.Getenv("POD_NAME")),
			},
		},
	})
	if err != nil {
		return err
	}

	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	return nil
}
