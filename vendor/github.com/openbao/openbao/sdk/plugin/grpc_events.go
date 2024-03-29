// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"

	"github.com/openbao/openbao/sdk/logical"
	"github.com/openbao/openbao/sdk/plugin/pb"
	"google.golang.org/grpc"
)

func newGRPCEventsClient(conn *grpc.ClientConn) *GRPCEventsClient {
	return &GRPCEventsClient{
		client: pb.NewEventsClient(conn),
	}
}

type GRPCEventsClient struct {
	client pb.EventsClient
}

var _ logical.EventSender = (*GRPCEventsClient)(nil)

func (s *GRPCEventsClient) Send(ctx context.Context, eventType logical.EventType, event *logical.EventData) error {
	_, err := s.client.SendEvent(ctx, &pb.SendEventRequest{
		EventType: string(eventType),
		Event:     event,
	})
	return err
}

type GRPCEventsServer struct {
	pb.UnimplementedEventsServer
	impl logical.EventSender
}

func (s *GRPCEventsServer) SendEvent(ctx context.Context, req *pb.SendEventRequest) (*pb.Empty, error) {
	if s.impl == nil {
		return &pb.Empty{}, nil
	}

	err := s.impl.Send(ctx, logical.EventType(req.EventType), req.Event)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
