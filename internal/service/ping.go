package service

import (
	"context"

	"github.com/bufbuild/connect-go"
	pingv1 "github.com/xdorro/base-project-proto/proto-gen-go/ping/v1"
)

// Ping is the ping service.
func (s *Service) Ping(
	_ context.Context, req *connect.Request[pingv1.PingRequest],
) (*connect.Response[pingv1.PingResponse], error) {
	text := req.Msg.Text
	if text == "" {
		text = "pong"
	}

	res := connect.NewResponse(&pingv1.PingResponse{
		Text: text,
	})

	return res, nil
}
