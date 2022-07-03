package grpc_server

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/bufbuild/connect-go"

	"github.com/xdorro/golang-grpc-base-project/internal/eliza"
	elizav1 "github.com/xdorro/golang-grpc-base-project/proto-gen-go/eliza/v1"
)

type elizaServer struct{}

func NewElizaServer() *elizaServer {
	return &elizaServer{}
}

func (e *elizaServer) Say(
	_ context.Context,
	req *connect.Request[elizav1.SayRequest],
) (*connect.Response[elizav1.SayResponse], error) {
	reply, _ := eliza.Reply(req.Msg.Sentence) // ignore end-of-conversation detection
	return connect.NewResponse(&elizav1.SayResponse{
		Sentence: reply,
	}), nil
}

func (e *elizaServer) Converse(
	ctx context.Context,
	stream *connect.BidiStream[elizav1.ConverseRequest, elizav1.ConverseResponse],
) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		request, err := stream.Receive()
		if err != nil && errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return fmt.Errorf("receive request: %w", err)
		}
		reply, endSession := eliza.Reply(request.Sentence)
		if err := stream.Send(&elizav1.ConverseResponse{Sentence: reply}); err != nil {
			return fmt.Errorf("send response: %w", err)
		}
		if endSession {
			return nil
		}
	}
}
