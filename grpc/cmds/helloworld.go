package cmds

import (
	"context"
	"fmt"
	"github.com/Qitmeer/llama.go/config"
	"github.com/Qitmeer/llama.go/grpc/proto"
	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
)

type HelloWorld struct {
	callCount int

	// client
	conn   *grpc.ClientConn
	client proto.HelloWorldClient
}

func NewHelloWorld(s *grpc.Server) *HelloWorld {
	log.Trace("NewHelloWorld")
	hw := &HelloWorld{}

	proto.RegisterHelloWorldServer(s, hw)
	return hw
}

func (hw *HelloWorld) SayHelloWorld(ctx context.Context, in *proto.HelloWorldRequest) (*proto.HelloWorldResponse, error) {
	hw.callCount++
	showStr := fmt.Sprintf("CallCount=%d request: %s", hw.callCount, in.Referer)
	log.Info(showStr)
	return &proto.HelloWorldResponse{Message: showStr}, nil
}

func (hw *HelloWorld) Client() proto.HelloWorldClient {

	if hw.client == nil {
		// Set up a connection to the gRPC server.
		conn, err := grpc.Dial(config.DefaultGrpcEndpoint, grpc.WithInsecure())
		if err != nil {
			log.Error(fmt.Sprintf("did not connect: %v", err))
			return nil
		}
		hw.client = proto.NewHelloWorldClient(conn)

		log.Trace("New HelloWorldClient")
	}

	return hw.client
}

func (hw *HelloWorld) Close() {
	log.Trace("Close HelloWorld")
	if hw.conn != nil {
		hw.conn.Close()
	}
}
