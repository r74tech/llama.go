package cmds

import (
	"context"
	"fmt"
	"github.com/Qitmeer/llama.go/config"
	"github.com/Qitmeer/llama.go/grpc/proto"
	"github.com/Qitmeer/llama.go/wrapper"
	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
)

type Generate struct {
	// client
	conn   *grpc.ClientConn
	client proto.GenerateClient
	cfg    *config.Config
}

func NewGenerate(s *grpc.Server, cfg *config.Config) *Generate {
	log.Trace("NewGenerate")
	hw := &Generate{cfg: cfg}

	proto.RegisterGenerateServer(s, hw)
	return hw
}

func (k *Generate) Generate(ctx context.Context, in *proto.GenerateRequest) (*proto.GenerateResponse, error) {
	content, err := wrapper.LlamaProcess(in.Prompt)
	if err != nil {
		return nil, err
	}
	return &proto.GenerateResponse{Content: content}, nil
}

func (k *Generate) Client() proto.GenerateClient {

	if k.client == nil {
		// Set up a connection to the gRPC server.
		conn, err := grpc.Dial(config.DefaultGrpcEndpoint, grpc.WithInsecure())
		if err != nil {
			log.Error(fmt.Sprintf("did not connect: %v", err))
			return nil
		}
		k.client = proto.NewGenerateClient(conn)

		log.Trace("New GenerateClient")
	}

	return k.client
}

func (k *Generate) Close() {
	log.Trace("Close Generate")
	if k.conn != nil {
		k.conn.Close()
	}
}
